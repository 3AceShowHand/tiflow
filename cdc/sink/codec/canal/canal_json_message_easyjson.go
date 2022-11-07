// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package canal

import (
	json "encoding/json"

	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonB5056ee2DecodeGithubComPingcapTiflowCdcSinkCodecCanal(in *jlexer.Lexer, out *tidbExtension) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "commitTs":
			out.CommitTs = uint64(in.Uint64())
		case "watermarkTs":
			out.WatermarkTs = uint64(in.Uint64())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}

func easyjsonB5056ee2EncodeGithubComPingcapTiflowCdcSinkCodecCanal(out *jwriter.Writer, in tidbExtension) {
	out.RawByte('{')
	first := true
	_ = first
	if in.CommitTs != 0 {
		const prefix string = ",\"commitTs\":"
		first = false
		out.RawString(prefix[1:])
		out.Uint64(uint64(in.CommitTs))
	}
	if in.WatermarkTs != 0 {
		const prefix string = ",\"watermarkTs\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Uint64(uint64(in.WatermarkTs))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v tidbExtension) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonB5056ee2EncodeGithubComPingcapTiflowCdcSinkCodecCanal(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v tidbExtension) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonB5056ee2EncodeGithubComPingcapTiflowCdcSinkCodecCanal(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *tidbExtension) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonB5056ee2DecodeGithubComPingcapTiflowCdcSinkCodecCanal(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *tidbExtension) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonB5056ee2DecodeGithubComPingcapTiflowCdcSinkCodecCanal(l, v)
}

func easyjsonB5056ee2DecodeGithubComPingcapTiflowCdcSinkCodecCanal1(in *jlexer.Lexer, out *canalJSONMessageWithTiDBExtension) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	out.JSONMessage = new(JSONMessage)
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "_tidb":
			if in.IsNull() {
				in.Skip()
				out.Extensions = nil
			} else {
				if out.Extensions == nil {
					out.Extensions = new(tidbExtension)
				}
				(*out.Extensions).UnmarshalEasyJSON(in)
			}
		case "id":
			out.ID = int64(in.Int64())
		case "database":
			out.Schema = string(in.String())
		case "table":
			out.Table = string(in.String())
		case "pkNames":
			if in.IsNull() {
				in.Skip()
				out.PKNames = nil
			} else {
				in.Delim('[')
				if out.PKNames == nil {
					if !in.IsDelim(']') {
						out.PKNames = make([]string, 0, 4)
					} else {
						out.PKNames = []string{}
					}
				} else {
					out.PKNames = (out.PKNames)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.String())
					out.PKNames = append(out.PKNames, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "isDdl":
			out.IsDDL = bool(in.Bool())
		case "type":
			out.EventType = string(in.String())
		case "es":
			out.ExecutionTime = int64(in.Int64())
		case "ts":
			out.BuildTime = int64(in.Int64())
		case "sql":
			out.Query = string(in.String())
		case "sqlType":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				out.SQLType = make(map[string]int32)
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v2 int32
					v2 = int32(in.Int32())
					(out.SQLType)[key] = v2
					in.WantComma()
				}
				in.Delim('}')
			}
		case "mysqlType":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				out.MySQLType = make(map[string]string)
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v3 string
					v3 = string(in.String())
					(out.MySQLType)[key] = v3
					in.WantComma()
				}
				in.Delim('}')
			}
		case "data":
			if in.IsNull() {
				in.Skip()
				out.Data = nil
			} else {
				in.Delim('[')
				if out.Data == nil {
					if !in.IsDelim(']') {
						out.Data = make([]map[string]string, 0, 8)
					} else {
						out.Data = []map[string]string{}
					}
				} else {
					out.Data = (out.Data)[:0]
				}
				for !in.IsDelim(']') {
					var v4 map[string]string
					if in.IsNull() {
						in.Skip()
					} else {
						in.Delim('{')
						v4 = make(map[string]string)
						for !in.IsDelim('}') {
							key := string(in.String())
							in.WantColon()
							var v5 string
							v5 = string(in.String())
							(v4)[key] = v5
							in.WantComma()
						}
						in.Delim('}')
					}
					out.Data = append(out.Data, v4)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "old":
			if in.IsNull() {
				in.Skip()
				out.Old = nil
			} else {
				in.Delim('[')
				if out.Old == nil {
					if !in.IsDelim(']') {
						out.Old = make([]map[string]string, 0, 8)
					} else {
						out.Old = []map[string]string{}
					}
				} else {
					out.Old = (out.Old)[:0]
				}
				for !in.IsDelim(']') {
					var v6 map[string]string
					if in.IsNull() {
						in.Skip()
					} else {
						in.Delim('{')
						v6 = make(map[string]string)
						for !in.IsDelim('}') {
							key := string(in.String())
							in.WantColon()
							var v7 string
							v7 = string(in.String())
							(v6)[key] = v7
							in.WantComma()
						}
						in.Delim('}')
					}
					out.Old = append(out.Old, v6)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}

func easyjsonB5056ee2EncodeGithubComPingcapTiflowCdcSinkCodecCanal1(out *jwriter.Writer, in canalJSONMessageWithTiDBExtension) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"_tidb\":"
		out.RawString(prefix[1:])
		if in.Extensions == nil {
			out.RawString("null")
		} else {
			(*in.Extensions).MarshalEasyJSON(out)
		}
	}
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix)
		out.Int64(int64(in.ID))
	}
	{
		const prefix string = ",\"database\":"
		out.RawString(prefix)
		out.String(string(in.Schema))
	}
	{
		const prefix string = ",\"table\":"
		out.RawString(prefix)
		out.String(string(in.Table))
	}
	{
		const prefix string = ",\"pkNames\":"
		out.RawString(prefix)
		if in.PKNames == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v8, v9 := range in.PKNames {
				if v8 > 0 {
					out.RawByte(',')
				}
				out.String(string(v9))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"isDdl\":"
		out.RawString(prefix)
		out.Bool(bool(in.IsDDL))
	}
	{
		const prefix string = ",\"type\":"
		out.RawString(prefix)
		out.String(string(in.EventType))
	}
	{
		const prefix string = ",\"es\":"
		out.RawString(prefix)
		out.Int64(int64(in.ExecutionTime))
	}
	{
		const prefix string = ",\"ts\":"
		out.RawString(prefix)
		out.Int64(int64(in.BuildTime))
	}
	{
		const prefix string = ",\"sql\":"
		out.RawString(prefix)
		out.String(string(in.Query))
	}
	{
		const prefix string = ",\"sqlType\":"
		out.RawString(prefix)
		if in.SQLType == nil && (out.Flags&jwriter.NilMapAsEmpty) == 0 {
			out.RawString(`null`)
		} else {
			out.RawByte('{')
			v10First := true
			for v10Name, v10Value := range in.SQLType {
				if v10First {
					v10First = false
				} else {
					out.RawByte(',')
				}
				out.String(string(v10Name))
				out.RawByte(':')
				out.Int32(int32(v10Value))
			}
			out.RawByte('}')
		}
	}
	{
		const prefix string = ",\"mysqlType\":"
		out.RawString(prefix)
		if in.MySQLType == nil && (out.Flags&jwriter.NilMapAsEmpty) == 0 {
			out.RawString(`null`)
		} else {
			out.RawByte('{')
			v11First := true
			for v11Name, v11Value := range in.MySQLType {
				if v11First {
					v11First = false
				} else {
					out.RawByte(',')
				}
				out.String(string(v11Name))
				out.RawByte(':')
				out.String(string(v11Value))
			}
			out.RawByte('}')
		}
	}
	{
		const prefix string = ",\"data\":"
		out.RawString(prefix)
		if in.Data == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v12, v13 := range in.Data {
				if v12 > 0 {
					out.RawByte(',')
				}
				if v13 == nil && (out.Flags&jwriter.NilMapAsEmpty) == 0 {
					out.RawString(`null`)
				} else {
					out.RawByte('{')
					v14First := true
					for v14Name, v14Value := range v13 {
						if v14First {
							v14First = false
						} else {
							out.RawByte(',')
						}
						out.String(string(v14Name))
						out.RawByte(':')
						out.String(string(v14Value))
					}
					out.RawByte('}')
				}
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"old\":"
		out.RawString(prefix)
		if in.Old == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v15, v16 := range in.Old {
				if v15 > 0 {
					out.RawByte(',')
				}
				if v16 == nil && (out.Flags&jwriter.NilMapAsEmpty) == 0 {
					out.RawString(`null`)
				} else {
					out.RawByte('{')
					v17First := true
					for v17Name, v17Value := range v16 {
						if v17First {
							v17First = false
						} else {
							out.RawByte(',')
						}
						out.String(string(v17Name))
						out.RawByte(':')
						out.String(string(v17Value))
					}
					out.RawByte('}')
				}
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v canalJSONMessageWithTiDBExtension) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonB5056ee2EncodeGithubComPingcapTiflowCdcSinkCodecCanal1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v canalJSONMessageWithTiDBExtension) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonB5056ee2EncodeGithubComPingcapTiflowCdcSinkCodecCanal1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *canalJSONMessageWithTiDBExtension) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonB5056ee2DecodeGithubComPingcapTiflowCdcSinkCodecCanal1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *canalJSONMessageWithTiDBExtension) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonB5056ee2DecodeGithubComPingcapTiflowCdcSinkCodecCanal1(l, v)
}

func easyjsonB5056ee2DecodeGithubComPingcapTiflowCdcSinkCodecCanal2(in *jlexer.Lexer, out *JSONMessage) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = int64(in.Int64())
		case "database":
			out.Schema = string(in.String())
		case "table":
			out.Table = string(in.String())
		case "pkNames":
			if in.IsNull() {
				in.Skip()
				out.PKNames = nil
			} else {
				in.Delim('[')
				if out.PKNames == nil {
					if !in.IsDelim(']') {
						out.PKNames = make([]string, 0, 4)
					} else {
						out.PKNames = []string{}
					}
				} else {
					out.PKNames = (out.PKNames)[:0]
				}
				for !in.IsDelim(']') {
					var v18 string
					v18 = string(in.String())
					out.PKNames = append(out.PKNames, v18)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "isDdl":
			out.IsDDL = bool(in.Bool())
		case "type":
			out.EventType = string(in.String())
		case "es":
			out.ExecutionTime = int64(in.Int64())
		case "ts":
			out.BuildTime = int64(in.Int64())
		case "sql":
			out.Query = string(in.String())
		case "sqlType":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				out.SQLType = make(map[string]int32)
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v19 int32
					v19 = int32(in.Int32())
					(out.SQLType)[key] = v19
					in.WantComma()
				}
				in.Delim('}')
			}
		case "mysqlType":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				out.MySQLType = make(map[string]string)
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v20 string
					v20 = string(in.String())
					(out.MySQLType)[key] = v20
					in.WantComma()
				}
				in.Delim('}')
			}
		case "data":
			if in.IsNull() {
				in.Skip()
				out.Data = nil
			} else {
				in.Delim('[')
				if out.Data == nil {
					if !in.IsDelim(']') {
						out.Data = make([]map[string]string, 0, 8)
					} else {
						out.Data = []map[string]string{}
					}
				} else {
					out.Data = (out.Data)[:0]
				}
				for !in.IsDelim(']') {
					var v21 map[string]string
					if in.IsNull() {
						in.Skip()
					} else {
						in.Delim('{')
						v21 = make(map[string]string)
						for !in.IsDelim('}') {
							key := string(in.String())
							in.WantColon()
							var v22 string
							v22 = string(in.String())
							(v21)[key] = v22
							in.WantComma()
						}
						in.Delim('}')
					}
					out.Data = append(out.Data, v21)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "old":
			if in.IsNull() {
				in.Skip()
				out.Old = nil
			} else {
				in.Delim('[')
				if out.Old == nil {
					if !in.IsDelim(']') {
						out.Old = make([]map[string]string, 0, 8)
					} else {
						out.Old = []map[string]string{}
					}
				} else {
					out.Old = (out.Old)[:0]
				}
				for !in.IsDelim(']') {
					var v23 map[string]string
					if in.IsNull() {
						in.Skip()
					} else {
						in.Delim('{')
						v23 = make(map[string]string)
						for !in.IsDelim('}') {
							key := string(in.String())
							in.WantColon()
							var v24 string
							v24 = string(in.String())
							(v23)[key] = v24
							in.WantComma()
						}
						in.Delim('}')
					}
					out.Old = append(out.Old, v23)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}

func easyjsonB5056ee2EncodeGithubComPingcapTiflowCdcSinkCodecCanal2(out *jwriter.Writer, in JSONMessage) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.ID))
	}
	{
		const prefix string = ",\"database\":"
		out.RawString(prefix)
		out.String(string(in.Schema))
	}
	{
		const prefix string = ",\"table\":"
		out.RawString(prefix)
		out.String(string(in.Table))
	}
	{
		const prefix string = ",\"pkNames\":"
		out.RawString(prefix)
		if in.PKNames == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v25, v26 := range in.PKNames {
				if v25 > 0 {
					out.RawByte(',')
				}
				out.String(string(v26))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"isDdl\":"
		out.RawString(prefix)
		out.Bool(bool(in.IsDDL))
	}
	{
		const prefix string = ",\"type\":"
		out.RawString(prefix)
		out.String(string(in.EventType))
	}
	{
		const prefix string = ",\"es\":"
		out.RawString(prefix)
		out.Int64(int64(in.ExecutionTime))
	}
	{
		const prefix string = ",\"ts\":"
		out.RawString(prefix)
		out.Int64(int64(in.BuildTime))
	}
	{
		const prefix string = ",\"sql\":"
		out.RawString(prefix)
		out.String(string(in.Query))
	}
	{
		const prefix string = ",\"sqlType\":"
		out.RawString(prefix)
		if in.SQLType == nil && (out.Flags&jwriter.NilMapAsEmpty) == 0 {
			out.RawString(`null`)
		} else {
			out.RawByte('{')
			v27First := true
			for v27Name, v27Value := range in.SQLType {
				if v27First {
					v27First = false
				} else {
					out.RawByte(',')
				}
				out.String(string(v27Name))
				out.RawByte(':')
				out.Int32(int32(v27Value))
			}
			out.RawByte('}')
		}
	}
	{
		const prefix string = ",\"mysqlType\":"
		out.RawString(prefix)
		if in.MySQLType == nil && (out.Flags&jwriter.NilMapAsEmpty) == 0 {
			out.RawString(`null`)
		} else {
			out.RawByte('{')
			v28First := true
			for v28Name, v28Value := range in.MySQLType {
				if v28First {
					v28First = false
				} else {
					out.RawByte(',')
				}
				out.String(string(v28Name))
				out.RawByte(':')
				out.String(string(v28Value))
			}
			out.RawByte('}')
		}
	}
	{
		const prefix string = ",\"data\":"
		out.RawString(prefix)
		if in.Data == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v29, v30 := range in.Data {
				if v29 > 0 {
					out.RawByte(',')
				}
				if v30 == nil && (out.Flags&jwriter.NilMapAsEmpty) == 0 {
					out.RawString(`null`)
				} else {
					out.RawByte('{')
					v31First := true
					for v31Name, v31Value := range v30 {
						if v31First {
							v31First = false
						} else {
							out.RawByte(',')
						}
						out.String(string(v31Name))
						out.RawByte(':')
						out.String(string(v31Value))
					}
					out.RawByte('}')
				}
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"old\":"
		out.RawString(prefix)
		if in.Old == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v32, v33 := range in.Old {
				if v32 > 0 {
					out.RawByte(',')
				}
				if v33 == nil && (out.Flags&jwriter.NilMapAsEmpty) == 0 {
					out.RawString(`null`)
				} else {
					out.RawByte('{')
					v34First := true
					for v34Name, v34Value := range v33 {
						if v34First {
							v34First = false
						} else {
							out.RawByte(',')
						}
						out.String(string(v34Name))
						out.RawByte(':')
						out.String(string(v34Value))
					}
					out.RawByte('}')
				}
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v JSONMessage) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonB5056ee2EncodeGithubComPingcapTiflowCdcSinkCodecCanal2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v JSONMessage) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonB5056ee2EncodeGithubComPingcapTiflowCdcSinkCodecCanal2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *JSONMessage) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonB5056ee2DecodeGithubComPingcapTiflowCdcSinkCodecCanal2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *JSONMessage) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonB5056ee2DecodeGithubComPingcapTiflowCdcSinkCodecCanal2(l, v)
}
