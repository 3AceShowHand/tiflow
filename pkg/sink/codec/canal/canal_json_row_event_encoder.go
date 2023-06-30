// Copyright 2022 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package canal

import (
	"context"
	"sort"
	"time"
	"unsafe"

	"github.com/goccy/go-json"
	"github.com/mailru/easyjson/jwriter"
	"github.com/pingcap/errors"
	"github.com/pingcap/log"
	"github.com/pingcap/tiflow/cdc/model"
	"github.com/pingcap/tiflow/pkg/config"
	cerror "github.com/pingcap/tiflow/pkg/errors"
	"github.com/pingcap/tiflow/pkg/sink/codec"
	"github.com/pingcap/tiflow/pkg/sink/codec/common"
	"go.uber.org/zap"
)

func newJSONMessageForDML(
	builder *canalEntryBuilder,
	e *model.RowChangedEvent,
	config *common.Config,
) ([]byte, error) {
	isDelete := e.IsDelete()
	mysqlTypeMap := make(map[string]string, len(e.Columns))

	statistics := make(map[string]int)

	filling := func(columns []*model.Column, out *jwriter.Writer,
		onlyOutputUpdatedColumn bool,
		onlyHandleKeyColumns bool,
		newColumnMap map[string]*model.Column,
	) error {
		type item struct {
			name      string
			mysqlType string
			value     interface{}

			origin int
			after  int
			rate   float64
		}
		expansion := make([]item, 0, len(columns))

		if len(columns) == 0 {
			out.RawString("null")
			return nil
		}
		out.RawByte('[')
		out.RawByte('{')
		isFirst := true
		for _, col := range columns {
			if col != nil {
				// column equal, do not output it
				if onlyOutputUpdatedColumn && shouldIgnoreColumn(col, newColumnMap) {
					continue
				}
				if onlyHandleKeyColumns && !col.Flag.IsHandleKey() {
					continue
				}
				if isFirst {
					isFirst = false
				} else {
					out.RawByte(',')
				}
				mysqlType := getMySQLType(col)
				javaType, err := getJavaSQLType(col, mysqlType)
				if err != nil {
					return cerror.WrapError(cerror.ErrCanalEncodeFailed, err)
				}

				originSize := int(unsafe.Sizeof(col.Value))
				value, err := builder.formatValue(col.Value, javaType)
				if err != nil {
					return cerror.WrapError(cerror.ErrCanalEncodeFailed, err)
				}

				old := out.Buffer.Size()
				out.String(col.Name)
				out.RawByte(':')
				if col.Value == nil {
					out.RawString("null")
				} else {
					out.String(value)
				}
				encodedSize := out.Buffer.Size() - old

				expansion = append(expansion, item{
					name:      col.Name,
					value:     col.Value,
					mysqlType: mysqlType,
					origin:    originSize,
					after:     encodedSize,
					rate:      float64(encodedSize) / float64(originSize),
				})

			}
		}
		out.RawByte('}')
		out.RawByte(']')

		sort.Slice(expansion, func(i, j int) bool {
			return expansion[i].rate > expansion[j].rate
		})

		var originTotal int
		var afterTotal int
		for _, item := range expansion {
			originTotal += item.origin
			afterTotal += item.after
			log.Info("expansion", zap.String("name", item.name), zap.Int("origin", item.origin), zap.Int("after", item.after), zap.Float64("rate", item.rate), zap.String("mysqlType", item.mysqlType), zap.Any("value", item.value))
		}
		log.Info("expansion Total", zap.Int("originTotal", originTotal), zap.Int("afterTotal", afterTotal), zap.Float64("rate", float64(afterTotal)/float64(originTotal)))

		return nil
	}

	out := &jwriter.Writer{}
	out.RawByte('{')

	old := out.Buffer.Size()
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int64(0) // ignored by both Canal Adapter and Flink
	}
	statistics["id"] = out.Buffer.Size() - old

	old = out.Buffer.Size()
	{
		const prefix string = ",\"database\":"
		out.RawString(prefix)
		out.String(e.Table.Schema)
	}
	statistics["database"] = out.Buffer.Size() - old

	old = out.Buffer.Size()
	{
		const prefix string = ",\"table\":"
		out.RawString(prefix)
		out.String(e.Table.Table)
	}
	statistics["table"] = out.Buffer.Size() - old

	old = out.Buffer.Size()
	{
		const prefix string = ",\"pkNames\":"
		out.RawString(prefix)
		pkNames := e.PrimaryKeyColumnNames()
		if pkNames == nil {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v25, v26 := range pkNames {
				if v25 > 0 {
					out.RawByte(',')
				}
				out.String(v26)
			}
			out.RawByte(']')
		}
	}
	statistics["pkNames"] = out.Buffer.Size() - old

	old = out.Buffer.Size()
	{
		const prefix string = ",\"isDdl\":"
		out.RawString(prefix)
		out.Bool(false)
	}
	statistics["isDDL"] = out.Buffer.Size() - old

	old = out.Buffer.Size()
	{
		const prefix string = ",\"type\":"
		out.RawString(prefix)
		out.String(eventTypeString(e))
	}
	statistics["type"] = out.Buffer.Size() - old

	old = out.Buffer.Size()
	{
		const prefix string = ",\"es\":"
		out.RawString(prefix)
		out.Int64(convertToCanalTs(e.CommitTs))
	}
	statistics["es"] = out.Buffer.Size() - old

	old = out.Buffer.Size()
	{
		const prefix string = ",\"ts\":"
		out.RawString(prefix)
		out.Int64(time.Now().UnixMilli()) // ignored by both Canal Adapter and Flink
	}
	statistics["ts"] = out.Buffer.Size() - old

	old = out.Buffer.Size()
	{
		const prefix string = ",\"sql\":"
		out.RawString(prefix)
		out.String("")
	}
	statistics["sql"] = out.Buffer.Size() - old

	old = out.Buffer.Size()
	{
		columns := e.PreColumns
		if !isDelete {
			columns = e.Columns
		}
		const prefix string = ",\"sqlType\":"
		out.RawString(prefix)
		emptyColumn := true
		for _, col := range columns {
			if col != nil {
				if isDelete && config.DeleteOnlyHandleKeyColumns && !col.Flag.IsHandleKey() {
					continue
				}
				if emptyColumn {
					out.RawByte('{')
					emptyColumn = false
				} else {
					out.RawByte(',')
				}
				mysqlType := getMySQLType(col)
				javaType, err := getJavaSQLType(col, mysqlType)
				if err != nil {
					return nil, cerror.WrapError(cerror.ErrCanalEncodeFailed, err)
				}
				out.String(col.Name)
				out.RawByte(':')
				out.Int32(int32(javaType))
				mysqlTypeMap[col.Name] = mysqlType
			}
		}
		if emptyColumn {
			out.RawString(`null`)
		} else {
			out.RawByte('}')
		}
	}
	statistics["sqlType"] = out.Buffer.Size() - old

	old = out.Buffer.Size()
	{
		const prefix string = ",\"mysqlType\":"
		out.RawString(prefix)
		if mysqlTypeMap == nil {
			out.RawString(`null`)
		} else {
			out.RawByte('{')
			isFirst := true
			for typeKey, typeValue := range mysqlTypeMap {
				if isFirst {
					isFirst = false
				} else {
					out.RawByte(',')
				}
				out.String(typeKey)
				out.RawByte(':')
				out.String(typeValue)
			}
			out.RawByte('}')
		}
	}
	statistics["mysqlType"] = out.Buffer.Size() - old

	old = out.Buffer.Size()
	if e.IsDelete() {
		out.RawString(",\"old\":null")
		out.RawString(",\"data\":")
		if err := filling(e.PreColumns, out, false, config.DeleteOnlyHandleKeyColumns, nil); err != nil {
			return nil, err
		}
	} else if e.IsInsert() {
		out.RawString(",\"old\":null")
		out.RawString(",\"data\":")
		if err := filling(e.Columns, out, false, false, nil); err != nil {
			return nil, err
		}
	} else if e.IsUpdate() {
		var newColsMap map[string]*model.Column
		if config.OnlyOutputUpdatedColumns {
			newColsMap = make(map[string]*model.Column, len(e.Columns))
			for _, col := range e.Columns {
				newColsMap[col.Name] = col
			}
		}
		out.RawString(",\"old\":")
		if err := filling(e.PreColumns, out, config.OnlyOutputUpdatedColumns, false, newColsMap); err != nil {
			return nil, err
		}
		out.RawString(",\"data\":")
		if err := filling(e.Columns, out, false, false, nil); err != nil {
			return nil, err
		}
	} else {
		log.Panic("unreachable event type", zap.Any("event", e))
	}
	statistics["data"] = out.Buffer.Size() - old

	old = out.Buffer.Size()
	if config.EnableTiDBExtension {
		const prefix string = ",\"_tidb\":"
		out.RawString(prefix)
		out.RawByte('{')
		out.RawString("\"commitTs\":")
		out.Uint64(e.CommitTs)
		out.RawByte('}')
	}
	out.RawByte('}')
	statistics["extension"] = out.Buffer.Size() - old

	var total int
	for _, item := range statistics {
		total += item
	}
	log.Info("statis")
	for key, value := range statistics {
		log.Info("statistics", zap.String("key", key), zap.Int("value", value), zap.Float64("percentage", float64(value)/float64(total)))
	}

	return out.BuildBytes()
}

func eventTypeString(e *model.RowChangedEvent) string {
	if e.IsDelete() {
		return "DELETE"
	}
	if len(e.PreColumns) == 0 {
		return "INSERT"
	}
	return "UPDATE"
}

// JSONRowEventEncoder encodes row event in JSON format
type JSONRowEventEncoder struct {
	builder  *canalEntryBuilder
	messages []*common.Message

	config *common.Config
}

// newJSONRowEventEncoder creates a new JSONRowEventEncoder
func newJSONRowEventEncoder(config *common.Config) codec.RowEventEncoder {
	encoder := &JSONRowEventEncoder{
		builder:  newCanalEntryBuilder(),
		messages: make([]*common.Message, 0, 1),

		config: config,
	}
	return encoder
}

func (c *JSONRowEventEncoder) newJSONMessageForDDL(e *model.DDLEvent) canalJSONMessageInterface {
	msg := &JSONMessage{
		ID:            0, // ignored by both Canal Adapter and Flink
		Schema:        e.TableInfo.TableName.Schema,
		Table:         e.TableInfo.TableName.Table,
		IsDDL:         true,
		EventType:     convertDdlEventType(e).String(),
		ExecutionTime: convertToCanalTs(e.CommitTs),
		BuildTime:     time.Now().UnixMilli(), // timestamp
		Query:         e.Query,
	}

	if !c.config.EnableTiDBExtension {
		return msg
	}

	return &canalJSONMessageWithTiDBExtension{
		JSONMessage: msg,
		Extensions:  &tidbExtension{CommitTs: e.CommitTs},
	}
}

func (c *JSONRowEventEncoder) newJSONMessage4CheckpointEvent(
	ts uint64,
) *canalJSONMessageWithTiDBExtension {
	return &canalJSONMessageWithTiDBExtension{
		JSONMessage: &JSONMessage{
			ID:            0,
			IsDDL:         false,
			EventType:     tidbWaterMarkType,
			ExecutionTime: convertToCanalTs(ts),
			BuildTime:     time.Now().UnixNano() / int64(time.Millisecond), // converts to milliseconds
		},
		Extensions: &tidbExtension{WatermarkTs: ts},
	}
}

// EncodeCheckpointEvent implements the RowEventEncoder interface
func (c *JSONRowEventEncoder) EncodeCheckpointEvent(ts uint64) (*common.Message, error) {
	if !c.config.EnableTiDBExtension {
		return nil, nil
	}

	msg := c.newJSONMessage4CheckpointEvent(ts)
	value, err := json.Marshal(msg)
	if err != nil {
		return nil, cerror.WrapError(cerror.ErrCanalEncodeFailed, err)
	}
	return common.NewResolvedMsg(config.ProtocolCanalJSON, nil, value, ts), nil
}

// AppendRowChangedEvent implements the interface EventJSONBatchEncoder
func (c *JSONRowEventEncoder) AppendRowChangedEvent(
	_ context.Context,
	_ string,
	e *model.RowChangedEvent,
	callback func(),
) error {
	value, err := newJSONMessageForDML(c.builder, e, c.config)
	if err != nil {
		return errors.Trace(err)
	}

	length := len(value) + common.MaxRecordOverhead
	// for single message that is longer than max-message-bytes, do not send it.
	if length > c.config.MaxMessageBytes {
		log.Warn("Single message is too large for canal-json",
			zap.Int("maxMessageBytes", c.config.MaxMessageBytes),
			zap.Int("length", length),
			zap.Any("table", e.Table))
		return cerror.ErrMessageTooLarge.GenWithStackByArgs()
	}
	m := &common.Message{
		Key:      nil,
		Value:    value,
		Ts:       e.CommitTs,
		Schema:   &e.Table.Schema,
		Table:    &e.Table.Table,
		Type:     model.MessageTypeRow,
		Protocol: config.ProtocolCanalJSON,
		Callback: callback,
	}
	m.IncRowsCount()

	c.messages = append(c.messages, m)
	return nil
}

// Build implements the RowEventEncoder interface
func (c *JSONRowEventEncoder) Build() []*common.Message {
	if len(c.messages) == 0 {
		return nil
	}

	result := c.messages
	c.messages = nil
	return result
}

// EncodeDDLEvent encodes DDL events
func (c *JSONRowEventEncoder) EncodeDDLEvent(e *model.DDLEvent) (*common.Message, error) {
	message := c.newJSONMessageForDDL(e)
	value, err := json.Marshal(message)
	if err != nil {
		return nil, cerror.WrapError(cerror.ErrCanalEncodeFailed, err)
	}
	return common.NewDDLMsg(config.ProtocolCanalJSON, nil, value, e), nil
}

type jsonRowEventEncoderBuilder struct {
	config *common.Config
}

// NewJSONRowEventEncoderBuilder creates a canal-json batchEncoderBuilder.
func NewJSONRowEventEncoderBuilder(config *common.Config) codec.RowEventEncoderBuilder {
	return &jsonRowEventEncoderBuilder{config: config}
}

// Build a `jsonRowEventEncoderBuilder`
func (b *jsonRowEventEncoderBuilder) Build() codec.RowEventEncoder {
	return newJSONRowEventEncoder(b.config)
}

func shouldIgnoreColumn(col *model.Column,
	newColumnMap map[string]*model.Column,
) bool {
	newCol, ok := newColumnMap[col.Name]
	if ok && newCol != nil {
		// sql type is not equal
		if newCol.Type != col.Type {
			return false
		}
		// value equal
		if codec.IsColumnValueEqual(newCol.Value, col.Value) {
			return true
		}
	}
	return false
}
