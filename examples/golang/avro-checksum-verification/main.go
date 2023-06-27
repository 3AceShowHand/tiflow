package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"hash/crc32"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/linkedin/goavro/v2"
	"github.com/pingcap/log"
	"github.com/pingcap/tidb/parser/mysql"
	"github.com/pingcap/tidb/types"
	"github.com/pingcap/tiflow/pkg/errors"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

const (
	// confluent avro wire format, the first byte is always 0
	// https://docs.confluent.io/platform/current/schema-registry/fundamentals/serdes-develop/index.html#wire-format
	magicByte = uint8(0)

	tidbOp = "_tidb_op"
)

func main() {
	var (
		kafkaAddr         = "127.0.0.1:9092"
		schemaRegistryURL = "http://127.0.0.1:8081"

		topic           = "avro-checksum-test"
		consumerGroupID = "avro-checksum-test"
	)

	consumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{kafkaAddr},
		GroupID:  consumerGroupID,
		Topic:    topic,
		MaxBytes: 10e6, // 10MB
	})
	defer consumer.Close()

	ctx := context.Background()
	log.Info("start consuming ...", zap.String("kafka", kafkaAddr), zap.String("topic", topic), zap.String("groupID", consumerGroupID))
	for {
		// 1. 获取 kafka 消息
		message, err := consumer.FetchMessage(ctx)
		if err != nil {
			log.Error("read kafka message failed", zap.Error(err))
		}

		value := message.Value
		if len(value) == 0 {
			log.Info("delete event does not have value, skip checksum verification", zap.String("topic", topic))
		}

		// 2. 对 value 部分进行解码，得到对应的 value 和 schema
		valueMap, valueSchema, err := getValueMapAndSchema(value, schemaRegistryURL)
		if err != nil {
			log.Panic("decode kafka value failed", zap.String("topic", topic), zap.ByteString("value", value), zap.Error(err))
		}

		// 3. 使用上一步得到的 value 和 schema，计算并且校验 checksum
		err = CalculateAndVerifyChecksum(valueMap, valueSchema)
		if err != nil {
			log.Panic("calculate checksum failed", zap.String("topic", topic), zap.ByteString("value", value), zap.Error(err))
		}

		// 4. 数据消费成功，提交 offset
		if err := consumer.CommitMessages(ctx, message); err != nil {
			log.Error("commit kafka message failed", zap.Error(err))
			break
		}
	}
}

func getValueMapAndSchema(data []byte, url string) (map[string]interface{}, map[string]interface{}, error) {
	schemaID, binary, err := extractSchemaIDAndBinaryData(data)
	if err != nil {
		return nil, nil, err
	}

	codec, err := GetSchema(url, schemaID)
	if err != nil {
		return nil, nil, err
	}

	native, _, err := codec.NativeFromBinary(binary)
	if err != nil {
		return nil, nil, err
	}

	result, ok := native.(map[string]interface{})
	if !ok {
		return nil, nil, errors.New("raw avro message is not a map")
	}

	schema := make(map[string]interface{})
	if err := json.Unmarshal([]byte(codec.Schema()), &schema); err != nil {
		return nil, nil, errors.Trace(err)
	}

	return result, schema, nil
}

func extractSchemaIDAndBinaryData(data []byte) (int, []byte, error) {
	if len(data) < 5 {
		return 0, nil, errors.ErrAvroInvalidMessage.FastGenByArgs()
	}
	if data[0] != magicByte {
		return 0, nil, errors.ErrAvroInvalidMessage.FastGenByArgs()
	}
	return int(binary.BigEndian.Uint32(data[1:5])), data[5:], nil
}

func CalculateAndVerifyChecksum(valueMap, valueSchema map[string]interface{}) error {
	// fields 存放有数据变更事件的每一个列的类型信息，按照每一列的 ID 排序，该顺序和 Checksum 计算顺序相同。
	fields, ok := valueSchema["fields"].([]interface{})
	if !ok {
		return errors.New("schema fields should be a map")
	}

	// 1. 从收到的 valueMap 里面，获取到期望的 checksum 值，它被编码成 string 类型
	// 如果找不到期望的 checksum 值，说明 TiCDC 发送该条数据时，还没有开启 checksum 功能，直接返回。
	o, ok := valueMap["_tidb_row_level_checksum"]
	if !ok {
		return nil
	}
	expected := o.(string)
	if expected == "" {
		return nil
	}

	// expectedChecksum 即是从 TiCDC 传递而来的期望的 checksum 值
	expectedChecksum, err := strconv.ParseUint(expected, 10, 64)
	if err != nil {
		return errors.Trace(err)
	}

	// 2. 遍历每一个 field，计算 Checksum 值
	var actualChecksum uint32
	// buf 用来存储每次更新 checksum 时使用的 bytes。
	buf := make([]byte, 0)
	for _, item := range fields {
		field, ok := item.(map[string]interface{})
		if !ok {
			return errors.New("schema field should be a map")
		}

		// `tidbOp` 及之后的列不参与到 checksum 计算中，因为它们是一些用于辅助数据消费的列，并非真实的 TiDB 列数据
		colName := field["name"].(string)
		if colName == tidbOp {
			break
		}

		// holder 存放有列类型信息
		var holder map[string]interface{}
		switch ty := field["type"].(type) {
		case []interface{}:
			for _, item := range ty {
				if m, ok := item.(map[string]interface{}); ok {
					holder = m["connect.parameters"].(map[string]interface{})
					break
				}
			}
		case map[string]interface{}:
			holder = ty["connect.parameters"].(map[string]interface{})
		default:
			log.Panic("type info is anything else", zap.Any("typeInfo", field["type"]))
		}
		tidbType := holder["tidb_type"].(string)

		mysqlType := mysqlTypeFromTiDBType(tidbType)

		// 根据每一列的名字，从解码之后的 value map 里拿到该列的值。
		value, ok := valueMap[colName]
		if !ok {
			return errors.New("value not found")
		}
		value, err := getColumnValue(value, holder, mysqlType)
		if err != nil {
			return errors.Trace(err)
		}

		if len(buf) > 0 {
			buf = buf[:0]
		}

		// 根据每一列的 value 和 mysqlType，生成用于更新 checksum 的字节切片
		buf, err = buildChecksumBytes(buf, value, mysqlType)
		if err != nil {
			return errors.Trace(err)
		}
		actualChecksum = crc32.Update(actualChecksum, crc32.IEEETable, buf)
	}

	if uint64(actualChecksum) != expectedChecksum {
		log.Error("checksum mismatch",
			zap.Uint64("expected", expectedChecksum),
			zap.Uint64("actual", uint64(actualChecksum)))
		return errors.New("checksum mismatch")
	}

	log.Info("checksum verified", zap.Uint64("checksum", uint64(actualChecksum)))

	return nil
}

func mysqlTypeFromTiDBType(tidbType string) byte {
	var result byte
	switch tidbType {
	case "INT", "INT UNSIGNED":
		result = mysql.TypeLong
	case "BIGINT", "BIGINT UNSIGNED":
		result = mysql.TypeLonglong
	case "FLOAT":
		result = mysql.TypeFloat
	case "DOUBLE":
		result = mysql.TypeDouble
	case "BIT":
		result = mysql.TypeBit
	case "DECIMAL":
		result = mysql.TypeNewDecimal
	case "TEXT":
		result = mysql.TypeVarchar
	case "BLOB":
		result = mysql.TypeLongBlob
	case "ENUM":
		result = mysql.TypeEnum
	case "SET":
		result = mysql.TypeSet
	case "JSON":
		result = mysql.TypeJSON
	case "DATE":
		result = mysql.TypeDate
	case "DATETIME":
		result = mysql.TypeDatetime
	case "TIMESTAMP":
		result = mysql.TypeTimestamp
	case "TIME":
		result = mysql.TypeDuration
	case "YEAR":
		result = mysql.TypeYear
	default:
		log.Panic("this should not happen, unknown TiDB type", zap.String("type", tidbType))
	}
	return result
}

// value is an interface, need to convert it to the real value with the help of type info.
// holder has the value's column info.
func getColumnValue(value interface{}, holder map[string]interface{}, mysqlType byte) (interface{}, error) {
	switch t := value.(type) {
	// for nullable columns, the value is encoded as a map with one pair.
	// key is the encoded type, value is the encoded value, only care about the value here.
	case map[string]interface{}:
		for _, v := range t {
			value = v
		}
	}

	switch mysqlType {
	case mysql.TypeEnum:
		// enum type is encoded as string,
		// we need to convert it to int by the order of the enum values definition.
		allowed := strings.Split(holder["allowed"].(string), ",")
		switch t := value.(type) {
		case string:
			enum, err := types.ParseEnum(allowed, t, "")
			if err != nil {
				return nil, errors.Trace(err)
			}
			value = enum.Value
		case nil:
			value = nil
		}
	case mysql.TypeSet:
		// set type is encoded as string,
		// we need to convert it to int by the order of the set values definition.
		elems := strings.Split(holder["allowed"].(string), ",")
		switch t := value.(type) {
		case string:
			s, err := types.ParseSet(elems, t, "")
			if err != nil {
				return nil, errors.Trace(err)
			}
			value = s.Value
		case nil:
			value = nil
		}
	}
	return value, nil
}

// buildChecksumBytes append value the buf, mysqlType is used to is used to convert value interface to concrete type.
// by follow: https://github.com/pingcap/tidb/blob/e3417913f58cdd5a136259b902bf177eaf3aa637/util/rowcodec/common.go#L308
func buildChecksumBytes(buf []byte, value interface{}, mysqlType byte) ([]byte, error) {
	if value == nil {
		return buf, nil
	}

	switch mysqlType {
	// TypeTiny, TypeShort, TypeInt32 is encoded as int32
	// TypeLong is encoded as int32 if signed, else int64.
	// TypeLongLong is encoded as int64 if signed, else uint64,
	// if bigintUnsignedHandlingMode set as string, encode as string.
	case mysql.TypeTiny, mysql.TypeShort, mysql.TypeLong, mysql.TypeLonglong, mysql.TypeInt24, mysql.TypeYear:
		switch a := value.(type) {
		case int32:
			buf = binary.LittleEndian.AppendUint64(buf, uint64(a))
		case uint32:
			buf = binary.LittleEndian.AppendUint64(buf, uint64(a))
		case int64:
			buf = binary.LittleEndian.AppendUint64(buf, uint64(a))
		case uint64:
			buf = binary.LittleEndian.AppendUint64(buf, a)
		case string:
			v, err := strconv.ParseUint(a, 10, 64)
			if err != nil {
				return nil, errors.Trace(err)
			}
			buf = binary.LittleEndian.AppendUint64(buf, v)
		default:
			log.Panic("unknown golang type for the integral value",
				zap.Any("value", value), zap.Any("mysqlType", mysqlType))
		}
	// TypeFloat encoded as float32, TypeDouble encoded as float64
	case mysql.TypeFloat, mysql.TypeDouble:
		var v float64
		switch a := value.(type) {
		case float32:
			v = float64(a)
		case float64:
			v = a
		}
		if math.IsInf(v, 0) || math.IsNaN(v) {
			v = 0
		}
		buf = binary.LittleEndian.AppendUint64(buf, math.Float64bits(v))
	// TypeEnum, TypeSet encoded as string
	// but convert to int by the getColumnValue function
	case mysql.TypeEnum, mysql.TypeSet:
		buf = binary.LittleEndian.AppendUint64(buf, value.(uint64))
	// TypeBit encoded as bytes
	case mysql.TypeBit:
		// bit is store as bytes, convert to uint64.
		v, err := binaryLiteralToInt(value.([]byte))
		if err != nil {
			return nil, errors.Trace(err)
		}
		buf = binary.LittleEndian.AppendUint64(buf, v)
	// encoded as bytes if binary flag set to true, else string
	case mysql.TypeVarchar, mysql.TypeVarString, mysql.TypeString, mysql.TypeTinyBlob, mysql.TypeMediumBlob, mysql.TypeLongBlob, mysql.TypeBlob:
		switch a := value.(type) {
		case string:
			buf = appendLengthValue(buf, []byte(a))
		case []byte:
			buf = appendLengthValue(buf, a)
		default:
			log.Panic("unknown golang type for the string value",
				zap.Any("value", value), zap.Any("mysqlType", mysqlType))
		}
	// all encoded as string
	case mysql.TypeTimestamp, mysql.TypeDatetime, mysql.TypeDate, mysql.TypeDuration, mysql.TypeNewDate:
		v := value.(string)
		buf = appendLengthValue(buf, []byte(v))
	// encoded as string if decimalHandlingMode set to string, it's required to enable checksum.
	case mysql.TypeNewDecimal:
		buf = appendLengthValue(buf, []byte(value.(string)))
	// encoded as string
	case mysql.TypeJSON:
		buf = appendLengthValue(buf, []byte(value.(string)))
	// this should not happen, does not take into the checksum calculation.
	case mysql.TypeNull, mysql.TypeGeometry:
		// do nothing
	default:
		return buf, errors.New("invalid type for the checksum calculation")
	}
	return buf, nil
}

func appendLengthValue(buf []byte, val []byte) []byte {
	buf = binary.LittleEndian.AppendUint32(buf, uint32(len(val)))
	buf = append(buf, val...)
	return buf
}

// binaryLiteralToInt convert bytes into uint64,
// by follow https://github.com/pingcap/tidb/blob/e3417913f58cdd5a136259b902bf177eaf3aa637/types/binary_literal.go#L105
func binaryLiteralToInt(bytes []byte) (uint64, error) {
	bytes = trimLeadingZeroBytes(bytes)
	length := len(bytes)

	if length > 8 {
		log.Error("invalid bit value found", zap.ByteString("value", bytes))
		return math.MaxUint64, errors.New("invalid bit value")
	}

	if length == 0 {
		return 0, nil
	}

	// Note: the byte-order is BigEndian.
	val := uint64(bytes[0])
	for i := 1; i < length; i++ {
		val = (val << 8) | uint64(bytes[i])
	}
	return val, nil
}

func trimLeadingZeroBytes(bytes []byte) []byte {
	if len(bytes) == 0 {
		return bytes
	}
	pos, posMax := 0, len(bytes)-1
	for ; pos < posMax; pos++ {
		if bytes[pos] != 0 {
			break
		}
	}
	return bytes[pos:]
}

// GetSchema query the schema registry to fetch the schema by the schema id.
// return the goavro.Codec which can be used to encode and decode the data.
func GetSchema(url string, schemaID int) (*goavro.Codec, error) {
	requestURI := url + "/schemas/ids/" + strconv.Itoa(schemaID)

	req, err := http.NewRequest("GET", requestURI, nil)
	if err != nil {
		log.Error("Cannot create the request to look up the schema", zap.Error(err))
		return nil, errors.WrapError(errors.ErrAvroSchemaAPIError, err)
	}
	req.Header.Add(
		"Accept",
		"application/vnd.schemaregistry.v1+json, application/vnd.schemaregistry+json, "+
			"application/json",
	)

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Cannot parse the lookup schema response", zap.Error(err))
		return nil, errors.WrapError(errors.ErrAvroSchemaAPIError, err)
	}

	if resp.StatusCode == 404 {
		log.Warn("Specified schema not found in Registry", zap.String("requestURI", requestURI), zap.Int("schemaID", schemaID))
		return nil, errors.ErrAvroSchemaAPIError.GenWithStackByArgs("Schema not found in Registry")
	}

	if resp.StatusCode != 200 {
		log.Error("Failed to query schema from the Registry, HTTP error",
			zap.Int("status", resp.StatusCode), zap.String("uri", requestURI), zap.ByteString("responseBody", body))
		return nil, errors.ErrAvroSchemaAPIError.GenWithStack("Failed to query schema from the Registry, HTTP error")
	}

	var jsonResp lookupResponse
	err = json.Unmarshal(body, &jsonResp)
	if err != nil {
		log.Error("Failed to parse result from Registry", zap.Error(err))
		return nil, errors.WrapError(errors.ErrAvroSchemaAPIError, err)
	}

	codec, err := goavro.NewCodec(jsonResp.Schema)
	if err != nil {
		return nil, errors.WrapError(errors.ErrAvroSchemaAPIError, err)
	}
	return codec, nil
}

type lookupResponse struct {
	Name     string `json:"name"`
	SchemaID int    `json:"id"`
	Schema   string `json:"schema"`
}
