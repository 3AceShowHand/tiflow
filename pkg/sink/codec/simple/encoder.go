package simple

type Message struct {
	Schema  string   `json:"database"`
	Table   string   `json:"table"`
	PKNames []string `json:"pkNames"`

	EventType string `json:"type"`
	// officially the timestamp of the event-time of the message, in milliseconds since Epoch.
	ExecutionTime int64 `json:"es"`
	// officially the timestamp of building the message, in milliseconds since Epoch.
	BuildTime int64 `json:"ts"`

	CommitTs uint64 `json:"commitTs"`

	ClaimCheckLocation string `json:"claimCheckLocation,omitempty"`

	// SQL that generated the change event, DDL or Query
	Query string `json:"sql"`

	Data map[string]interface{} `json:"data"`
	Old  map[string]interface{} `json:"old"`

	//// only works for INSERT / UPDATE / DELETE events, records each column's java representation type.
	//SQLType map[string]int32 `json:"sqlType"`
	//// only works for INSERT / UPDATE / DELETE events, records each column's mysql representation type.
	//MySQLType map[string]string `json:"mysqlType"`
	//// A Datum should be a string or nil
	//Data []map[string]interface{} `json:"data"`
	//Old  []map[string]interface{} `json:"old"`
}

type ColumnSchema struct {
	MySQLType    string      `json:"mysqlType"`
	DefaultValue interface{} `json:"defaultValue"`
	Nullable     bool        `json:"nullable"`
}

type TableSchema struct {
	Version uint64          `json:"tableVersion"`
	Columns []*ColumnSchema `json:"columns"`
}
