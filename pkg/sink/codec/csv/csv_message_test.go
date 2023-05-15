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

package csv

import (
	"fmt"
	"strings"
	"testing"

	"github.com/pingcap/tidb/parser/charset"
	timodel "github.com/pingcap/tidb/parser/model"
	"github.com/pingcap/tidb/parser/mysql"
	"github.com/pingcap/tidb/types"
	"github.com/pingcap/tiflow/cdc/model"
	"github.com/pingcap/tiflow/pkg/sink/codec/common"
	"github.com/stretchr/testify/require"
)

type csvTestColumnTuple struct {
	col  model.Column
	want interface{}
}

var csvTestColumnsGroup = [][]*csvTestColumnTuple{
	{
		{
			model.Column{
				ID:        1,
				Name:      "tiny",
				Value:     int64(1),
				Type:      mysql.TypeTiny,
				FieldType: types.NewFieldType(mysql.TypeTiny),
			},
			int64(1),
		},
		{
			model.Column{
				ID:        2,
				Name:      "short",
				Value:     int64(1),
				Type:      mysql.TypeShort,
				FieldType: types.NewFieldType(mysql.TypeShort),
			},
			int64(1),
		},
		{
			model.Column{
				ID:        3,
				Name:      "int24",
				Value:     int64(1),
				Type:      mysql.TypeInt24,
				FieldType: types.NewFieldType(mysql.TypeInt24),
			},
			int64(1),
		},
		{
			model.Column{
				ID:        4,
				Name:      "long",
				Value:     int64(1),
				Type:      mysql.TypeLong,
				FieldType: types.NewFieldType(mysql.TypeLong),
			},
			int64(1),
		},
		{
			model.Column{
				ID:        5,
				Name:      "longlong",
				Value:     int64(1),
				Type:      mysql.TypeLonglong,
				FieldType: types.NewFieldType(mysql.TypeLonglong),
			},
			int64(1),
		},
		{
			model.Column{
				ID:        6,
				Name:      "tinyunsigned",
				Value:     uint64(1),
				Type:      mysql.TypeTiny,
				Flag:      model.UnsignedFlag,
				FieldType: setFlag(types.NewFieldType(mysql.TypeTiny), uint(model.UnsignedFlag)),
			},
			uint64(1),
		},
		{
			model.Column{
				ID:        7,
				Name:      "shortunsigned",
				Value:     uint64(1),
				Type:      mysql.TypeShort,
				Flag:      model.UnsignedFlag,
				FieldType: setFlag(types.NewFieldType(mysql.TypeShort), uint(model.UnsignedFlag)),
			},
			uint64(1),
		},
		{
			model.Column{
				ID:        8,
				Name:      "int24unsigned",
				Value:     uint64(1),
				Type:      mysql.TypeInt24,
				Flag:      model.UnsignedFlag,
				FieldType: setFlag(types.NewFieldType(mysql.TypeInt24), uint(model.UnsignedFlag)),
			},
			uint64(1),
		},
		{
			model.Column{
				ID:        9,
				Name:      "longunsigned",
				Value:     uint64(1),
				Type:      mysql.TypeLong,
				Flag:      model.UnsignedFlag,
				FieldType: setFlag(types.NewFieldType(mysql.TypeLong), uint(model.UnsignedFlag)),
			},
			uint64(1),
		},
		{
			model.Column{
				ID:    10,
				Name:  "longlongunsigned",
				Value: uint64(1),
				Type:  mysql.TypeLonglong,
				Flag:  model.UnsignedFlag,
				FieldType: setFlag(
					types.NewFieldType(mysql.TypeLonglong),
					uint(model.UnsignedFlag),
				),
			},
			uint64(1),
		},
	},
	{
		{
			model.Column{
				ID:        11,
				Name:      "float",
				Value:     float64(3.14),
				Type:      mysql.TypeFloat,
				FieldType: types.NewFieldType(mysql.TypeFloat),
			},
			float64(3.14),
		},
		{
			model.Column{
				ID:        12,
				Name:      "double",
				Value:     float64(3.14),
				Type:      mysql.TypeDouble,
				FieldType: types.NewFieldType(mysql.TypeDouble),
			},
			float64(3.14),
		},
	},
	{
		{
			model.Column{
				ID:        13,
				Name:      "bit",
				Value:     uint64(683),
				Type:      mysql.TypeBit,
				FieldType: types.NewFieldType(mysql.TypeBit),
			},
			uint64(683),
		},
	},
	{
		{
			model.Column{
				ID:        14,
				Name:      "decimal",
				Value:     "129012.1230000",
				Type:      mysql.TypeNewDecimal,
				FieldType: types.NewFieldType(mysql.TypeNewDecimal),
			},
			"129012.1230000",
		},
	},
	{
		{
			model.Column{
				ID:        15,
				Name:      "tinytext",
				Value:     []byte("hello world"),
				Type:      mysql.TypeTinyBlob,
				FieldType: types.NewFieldType(mysql.TypeTinyBlob),
			},
			"hello world",
		},
		{
			model.Column{
				ID:        16,
				Name:      "mediumtext",
				Value:     []byte("hello world"),
				Type:      mysql.TypeMediumBlob,
				FieldType: types.NewFieldType(mysql.TypeMediumBlob),
			},
			"hello world",
		},
		{
			model.Column{
				ID:        17,
				Name:      "text",
				Value:     []byte("hello world"),
				Type:      mysql.TypeBlob,
				FieldType: types.NewFieldType(mysql.TypeBlob),
			},
			"hello world",
		},
		{
			model.Column{
				ID:        18,
				Name:      "longtext",
				Value:     []byte("hello world"),
				Type:      mysql.TypeLongBlob,
				FieldType: types.NewFieldType(mysql.TypeLongBlob),
			},
			"hello world",
		},
		{
			model.Column{
				ID:        19,
				Name:      "varchar",
				Value:     []byte("hello world"),
				Type:      mysql.TypeVarchar,
				FieldType: types.NewFieldType(mysql.TypeVarchar),
			},
			"hello world",
		},
		{
			model.Column{
				ID:        20,
				Name:      "varstring",
				Value:     []byte("hello world"),
				Type:      mysql.TypeVarString,
				FieldType: types.NewFieldType(mysql.TypeVarString),
			},
			"hello world",
		},
		{
			model.Column{
				ID:        21,
				Name:      "string",
				Value:     []byte("hello world"),
				Type:      mysql.TypeString,
				FieldType: types.NewFieldType(mysql.TypeString),
			},
			"hello world",
		},
		{
			model.Column{
				ID:        31,
				Name:      "json",
				Value:     `{"key": "value"}`,
				Type:      mysql.TypeJSON,
				FieldType: types.NewFieldType(mysql.TypeJSON),
			},
			`{"key": "value"}`,
		},
	},
	{
		{
			model.Column{
				ID:        22,
				Name:      "tinyblob",
				Value:     []byte("hello world"),
				Type:      mysql.TypeTinyBlob,
				Flag:      model.BinaryFlag,
				FieldType: setBinChsClnFlag(types.NewFieldType(mysql.TypeTinyBlob)),
			},
			"aGVsbG8gd29ybGQ=",
		},
		{
			model.Column{
				ID:        23,
				Name:      "mediumblob",
				Value:     []byte("hello world"),
				Type:      mysql.TypeMediumBlob,
				Flag:      model.BinaryFlag,
				FieldType: setBinChsClnFlag(types.NewFieldType(mysql.TypeMediumBlob)),
			},
			"aGVsbG8gd29ybGQ=",
		},
		{
			model.Column{
				ID:        24,
				Name:      "blob",
				Value:     []byte("hello world"),
				Type:      mysql.TypeBlob,
				Flag:      model.BinaryFlag,
				FieldType: setBinChsClnFlag(types.NewFieldType(mysql.TypeBlob)),
			},
			"aGVsbG8gd29ybGQ=",
		},
		{
			model.Column{
				ID:        25,
				Name:      "longblob",
				Value:     []byte("hello world"),
				Type:      mysql.TypeLongBlob,
				Flag:      model.BinaryFlag,
				FieldType: setBinChsClnFlag(types.NewFieldType(mysql.TypeLongBlob)),
			},
			"aGVsbG8gd29ybGQ=",
		},
		{
			model.Column{
				ID:        26,
				Name:      "varbinary",
				Value:     []byte("hello world"),
				Type:      mysql.TypeVarchar,
				Flag:      model.BinaryFlag,
				FieldType: setBinChsClnFlag(types.NewFieldType(mysql.TypeVarchar)),
			},
			"aGVsbG8gd29ybGQ=",
		},
		{
			model.Column{
				ID:        27,
				Name:      "varbinary1",
				Value:     []byte("hello world"),
				Type:      mysql.TypeVarString,
				Flag:      model.BinaryFlag,
				FieldType: setBinChsClnFlag(types.NewFieldType(mysql.TypeVarString)),
			},
			"aGVsbG8gd29ybGQ=",
		},
		{
			model.Column{
				ID:        28,
				Name:      "binary",
				Value:     []byte("hello world"),
				Type:      mysql.TypeString,
				Flag:      model.BinaryFlag,
				FieldType: setBinChsClnFlag(types.NewFieldType(mysql.TypeString)),
			},
			"aGVsbG8gd29ybGQ=",
		},
	},
	{
		{
			model.Column{
				ID:        29,
				Name:      "enum",
				Value:     uint64(1),
				Type:      mysql.TypeEnum,
				FieldType: setElems(types.NewFieldType(mysql.TypeEnum), []string{"a,", "b"}),
			},
			"a,",
		},
	},
	{
		{
			model.Column{
				ID:        30,
				Name:      "set",
				Value:     uint64(9),
				Type:      mysql.TypeSet,
				FieldType: setElems(types.NewFieldType(mysql.TypeSet), []string{"a", "b", "c", "d"}),
			},
			"a,d",
		},
	},
	{
		{
			model.Column{
				ID:        32,
				Name:      "date",
				Value:     "2000-01-01",
				Type:      mysql.TypeDate,
				FieldType: types.NewFieldType(mysql.TypeDate),
			},
			"2000-01-01",
		},
		{
			model.Column{
				ID:        33,
				Name:      "datetime",
				Value:     "2015-12-20 23:58:58",
				Type:      mysql.TypeDatetime,
				FieldType: types.NewFieldType(mysql.TypeDatetime),
			},
			"2015-12-20 23:58:58",
		},
		{
			model.Column{
				ID:        34,
				Name:      "timestamp",
				Value:     "1973-12-30 15:30:00",
				Type:      mysql.TypeTimestamp,
				FieldType: types.NewFieldType(mysql.TypeTimestamp),
			},
			"1973-12-30 15:30:00",
		},
		{
			model.Column{
				ID:        35,
				Name:      "time",
				Value:     "23:59:59",
				Type:      mysql.TypeDuration,
				FieldType: types.NewFieldType(mysql.TypeDuration),
			},
			"23:59:59",
		},
	},
	{
		{
			model.Column{
				ID:        36,
				Name:      "year",
				Value:     int64(1970),
				Type:      mysql.TypeYear,
				FieldType: types.NewFieldType(mysql.TypeYear),
			},
			int64(1970),
		},
	},
}

func setBinChsClnFlag(ft *types.FieldType) *types.FieldType {
	types.SetBinChsClnFlag(ft)
	return ft
}

//nolint:unparam
func setFlag(ft *types.FieldType, flag uint) *types.FieldType {
	ft.SetFlag(flag)
	return ft
}

func setElems(ft *types.FieldType, elems []string) *types.FieldType {
	ft.SetElems(elems)
	return ft
}

func TestFormatWithQuotes(t *testing.T) {
	config := &common.Config{
		Quote: "\"",
	}

	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "string does not contain quote mark",
			input:    "a,b,c",
			expected: `"a,b,c"`,
		},
		{
			name:     "string contains quote mark",
			input:    `"a,b,c`,
			expected: `"""a,b,c"`,
		},
		{
			name:     "empty string",
			input:    "",
			expected: `""`,
		},
	}
	for _, tc := range testCases {
		csvMessage := newCSVMessage(config)
		strBuilder := new(strings.Builder)
		csvMessage.formatWithQuotes(tc.input, strBuilder)
		require.Equal(t, tc.expected, strBuilder.String(), tc.name)
	}
}

func TestFormatWithEscape(t *testing.T) {
	testCases := []struct {
		name     string
		config   *common.Config
		input    string
		expected string
	}{
		{
			name:     "string does not contain CR/LF/backslash/delimiter",
			config:   &common.Config{Delimiter: ","},
			input:    "abcdef",
			expected: "abcdef",
		},
		{
			name:     "string contains CRLF",
			config:   &common.Config{Delimiter: ","},
			input:    "abc\r\ndef",
			expected: "abc\\r\\ndef",
		},
		{
			name:     "string contains backslash",
			config:   &common.Config{Delimiter: ","},
			input:    `abc\def`,
			expected: `abc\\def`,
		},
		{
			name:     "string contains a single character delimiter",
			config:   &common.Config{Delimiter: ","},
			input:    "abc,def",
			expected: `abc\,def`,
		},
		{
			name:     "string contains multi-character delimiter",
			config:   &common.Config{Delimiter: "***"},
			input:    "abc***def",
			expected: `abc\*\*\*def`,
		},
		{
			name:     "string contains CR, LF, backslash and delimiter",
			config:   &common.Config{Delimiter: "?"},
			input:    `abc\def?ghi\r\n`,
			expected: `abc\\def\?ghi\\r\\n`,
		},
	}

	for _, tc := range testCases {
		csvMessage := newCSVMessage(tc.config)
		strBuilder := new(strings.Builder)
		csvMessage.formatWithEscapes(tc.input, strBuilder)
		require.Equal(t, tc.expected, strBuilder.String())
	}
}

func TestCSVMessageEncode(t *testing.T) {
	type fields struct {
		config     *common.Config
		opType     operation
		tableName  string
		schemaName string
		commitTs   uint64
		columns    []any
	}
	testCases := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name: "csv encode with typical configurations",
			fields: fields{
				config: &common.Config{
					Delimiter:       ",",
					Quote:           "\"",
					Terminator:      "\n",
					NullString:      "\\N",
					IncludeCommitTs: true,
				},
				opType:     operationInsert,
				tableName:  "table1",
				schemaName: "test",
				commitTs:   435661838416609281,
				columns:    []any{123, "hello,world"},
			},
			want: []byte("\"I\",\"table1\",\"test\",435661838416609281,123,\"hello,world\"\n"),
		},
		{
			name: "csv encode values containing single-character delimter string, without quote mark",
			fields: fields{
				config: &common.Config{
					Delimiter:       "!",
					Quote:           "",
					Terminator:      "\n",
					NullString:      "\\N",
					IncludeCommitTs: true,
				},
				opType:     operationUpdate,
				tableName:  "table2",
				schemaName: "test",
				commitTs:   435661838416609281,
				columns:    []any{"a!b!c", "def"},
			},
			want: []byte(`U!table2!test!435661838416609281!a\!b\!c!def` + "\n"),
		},
		{
			name: "csv encode values containing single-character delimter string, with quote mark",
			fields: fields{
				config: &common.Config{
					Delimiter:       ",",
					Quote:           "\"",
					Terminator:      "\n",
					NullString:      "\\N",
					IncludeCommitTs: true,
				},
				opType:     operationUpdate,
				tableName:  "table3",
				schemaName: "test",
				commitTs:   435661838416609281,
				columns:    []any{"a,b,c", "def", "2022-08-31 17:07:00"},
			},
			want: []byte(`"U","table3","test",435661838416609281,"a,b,c","def","2022-08-31 17:07:00"` + "\n"),
		},
		{
			name: "csv encode values containing multi-character delimiter string, without quote mark",
			fields: fields{
				config: &common.Config{
					Delimiter:       "[*]",
					Quote:           "",
					Terminator:      "\r\n",
					NullString:      "\\N",
					IncludeCommitTs: false,
				},
				opType:     operationDelete,
				tableName:  "table4",
				schemaName: "test",
				commitTs:   435661838416609281,
				columns:    []any{"a[*]b[*]c", "def"},
			},
			want: []byte(`D[*]table4[*]test[*]a\[\*\]b\[\*\]c[*]def` + "\r\n"),
		},
		{
			name: "csv encode with values containing multi-character delimiter string, with quote mark",
			fields: fields{
				config: &common.Config{
					Delimiter:       "[*]",
					Quote:           "'",
					Terminator:      "\n",
					NullString:      "\\N",
					IncludeCommitTs: false,
				},
				opType:     operationInsert,
				tableName:  "table5",
				schemaName: "test",
				commitTs:   435661838416609281,
				columns:    []any{"a[*]b[*]c", "def", nil, 12345.678},
			},
			want: []byte(`'I'[*]'table5'[*]'test'[*]'a[*]b[*]c'[*]'def'[*]\N[*]12345.678` + "\n"),
		},
		{
			name: "csv encode with values containing backslash and LF, without quote mark",
			fields: fields{
				config: &common.Config{
					Delimiter:       ",",
					Quote:           "",
					Terminator:      "\n",
					NullString:      "\\N",
					IncludeCommitTs: true,
				},
				opType:     operationUpdate,
				tableName:  "table6",
				schemaName: "test",
				commitTs:   435661838416609281,
				columns:    []any{"a\\b\\c", "def\n"},
			},
			want: []byte(`U,table6,test,435661838416609281,a\\b\\c,def\n` + "\n"),
		},
		{
			name: "csv encode with values containing backslash and CR, with quote mark",
			fields: fields{
				config: &common.Config{
					Delimiter:       ",",
					Quote:           "'",
					Terminator:      "\n",
					NullString:      "\\N",
					IncludeCommitTs: false,
				},
				opType:     operationInsert,
				tableName:  "table7",
				schemaName: "test",
				commitTs:   435661838416609281,
				columns:    []any{"\\", "\\\r", "\\\\"},
			},
			want: []byte("'I','table7','test','\\','\\\r','\\\\'" + "\n"),
		},
		{
			name: "csv encode with values containing unicode characters",
			fields: fields{
				config: &common.Config{
					Delimiter:       "\t",
					Quote:           "\"",
					Terminator:      "\n",
					NullString:      "\\N",
					IncludeCommitTs: true,
				},
				opType:     operationDelete,
				tableName:  "table8",
				schemaName: "test",
				commitTs:   435661838416609281,
				columns:    []any{"a\tb", 123.456, "你好，世界"},
			},
			want: []byte("\"D\"\t\"table8\"\t\"test\"\t435661838416609281\t\"a\tb\"\t123.456\t\"你好，世界\"\n"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := &csvMessage{
				config:     tc.fields.config,
				opType:     tc.fields.opType,
				tableName:  tc.fields.tableName,
				schemaName: tc.fields.schemaName,
				commitTs:   tc.fields.commitTs,
				columns:    tc.fields.columns,
				newRecord:  true,
			}

			require.Equal(t, tc.want, c.encode())
		})
	}
}

func TestConvertToCSVType(t *testing.T) {
	for _, group := range csvTestColumnsGroup {
		for _, c := range group {
			val, _ := fromColValToCsvVal(&c.col)
			require.Equal(t, c.want, val, c.col.Name)
		}
	}
}

func TestRowChangeEventConversion(t *testing.T) {
	for idx, group := range csvTestColumnsGroup {
		row := &model.RowChangedEvent{}
		cols := make([]*model.Column, 0)
		for _, c := range group {
			cols = append(cols, &c.col)
		}
		row.Table = &model.TableName{
			Table:  fmt.Sprintf("table%d", idx),
			Schema: "test",
		}

		if idx%3 == 0 { // delete operation
			row.PreColumns = cols
		} else if idx%3 == 1 { // insert operation
			row.Columns = cols
		} else { // update operation
			row.PreColumns = cols
			row.Columns = cols
		}
		csvMsg, err := rowChangedEvent2CSVMsg(&common.Config{
			Delimiter:       "\t",
			Quote:           "\"",
			Terminator:      "\n",
			NullString:      "\\N",
			IncludeCommitTs: true,
		}, row)
		require.NotNil(t, csvMsg)
		require.Nil(t, err)

		ticols := make([]*timodel.ColumnInfo, 0)
		for _, col := range cols {
			ticol := &timodel.ColumnInfo{
				Name:      timodel.NewCIStr(col.Name),
				FieldType: *types.NewFieldType(col.Type),
			}
			if col.Flag.IsBinary() {
				ticol.SetCharset(charset.CharsetBin)
			} else {
				ticol.SetCharset(mysql.DefaultCharset)
			}
			ticols = append(ticols, ticol)
		}

		row2, err := csvMsg2RowChangedEvent(csvMsg, ticols)
		require.Nil(t, err)
		require.NotNil(t, row2)
	}
}

func TestCSVMessageDecode(t *testing.T) {
	// datums := make([][]types.Datum, 0, 4)
	testCases := []struct {
		row              []types.Datum
		expectedCommitTs uint64
		expectedColsCnt  int
		expectedErr      string
	}{
		{
			row: []types.Datum{
				types.NewStringDatum("I"),
				types.NewStringDatum("employee"),
				types.NewStringDatum("hr"),
				types.NewStringDatum("433305438660591626"),
				types.NewStringDatum("101"),
				types.NewStringDatum("Smith"),
				types.NewStringDatum("Bob"),
				types.NewStringDatum("2014-06-04"),
				types.NewDatum(nil),
			},
			expectedCommitTs: 433305438660591626,
			expectedColsCnt:  5,
			expectedErr:      "",
		},
		{
			row: []types.Datum{
				types.NewStringDatum("U"),
				types.NewStringDatum("employee"),
				types.NewStringDatum("hr"),
				types.NewStringDatum("433305438660591627"),
				types.NewStringDatum("101"),
				types.NewStringDatum("Smith"),
				types.NewStringDatum("Bob"),
				types.NewStringDatum("2015-10-08"),
				types.NewStringDatum("Los Angeles"),
			},
			expectedCommitTs: 433305438660591627,
			expectedColsCnt:  5,
			expectedErr:      "",
		},
		{
			row: []types.Datum{
				types.NewStringDatum("D"),
				types.NewStringDatum("employee"),
				types.NewStringDatum("hr"),
			},
			expectedCommitTs: 0,
			expectedColsCnt:  0,
			expectedErr:      "the csv row should have at least four columns",
		},
		{
			row: []types.Datum{
				types.NewStringDatum("D"),
				types.NewStringDatum("employee"),
				types.NewStringDatum("hr"),
				types.NewStringDatum("hello world"),
			},
			expectedCommitTs: 0,
			expectedColsCnt:  0,
			expectedErr:      "the 4th column(hello world) of csv row should be a valid commit-ts",
		},
	}
	for _, tc := range testCases {
		csvMsg := newCSVMessage(&common.Config{
			Delimiter:       ",",
			Quote:           "\"",
			Terminator:      "\n",
			NullString:      "\\N",
			IncludeCommitTs: true,
		})
		err := csvMsg.decode(tc.row)
		if tc.expectedErr != "" {
			require.Contains(t, err.Error(), tc.expectedErr)
		} else {
			require.Nil(t, err)
			require.Equal(t, tc.expectedCommitTs, csvMsg.commitTs)
			require.Equal(t, tc.expectedColsCnt, len(csvMsg.columns))
		}
	}
}
