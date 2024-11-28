// Copyright 2024 PingCAP, Inc.
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

package main

import (
	"context"
	"encoding/json"
	"github.com/pingcap/log"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type CanalJSONMessageWithExtension struct {
	ID        int64    `json:"id"`
	Schema    string   `json:"database"`
	Table     string   `json:"table"`
	PKNames   []string `json:"pkNames"`
	IsDDL     bool     `json:"isDdl"`
	EventType string   `json:"type"`
	// officially the timestamp of the event-time of the message, in milliseconds since Epoch.
	ExecutionTime int64 `json:"es"`
	// officially the timestamp of building the message, in milliseconds since Epoch.
	BuildTime int64 `json:"ts"`
	// SQL that generated the change event, DDL or Query
	Query string `json:"sql"`
	// only works for INSERT / UPDATE / DELETE events, records each column's java representation type.
	SQLType map[string]int32 `json:"sqlType"`
	// only works for INSERT / UPDATE / DELETE events, records each column's mysql representation type.
	MySQLType map[string]string `json:"mysqlType"`
	// A Datum should be a string or nil
	Data []map[string]interface{} `json:"data"`
	Old  []map[string]interface{} `json:"old"`

	// Extensions is a TiCDC custom field that different from official Canal-JSON format.
	// It would be useful to store something for special usage.
	// At the moment, only store the `tso` of each event,
	// which is useful if the message consumer needs to restore the original transactions.
	Extensions *tidbExtension `json:"_tidb"`
}

type tidbExtension struct {
	CommitTs           uint64 `json:"commitTs,omitempty"`
	WatermarkTs        uint64 `json:"watermarkTs,omitempty"`
	OnlyHandleKey      bool   `json:"onlyHandleKey,omitempty"`
	ClaimCheckLocation string `json:"claimCheckLocation,omitempty"`
}

func main() {
	var (
		kafkaAddr       = "127.0.0.1:9092"
		topic           = "canal-json-example"
		consumerGroupID = "canal-json-example-group"
	)

	consumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaAddr},
		GroupID: consumerGroupID,
		Topic:   topic,
	})
	defer consumer.Close()

	ctx := context.Background()
	log.Info("start consuming ...", zap.String("kafka", kafkaAddr), zap.String("topic", topic), zap.String("groupID", consumerGroupID))
	for {
		// 1. read message out, do not commit offset
		rawMessage, err := consumer.FetchMessage(ctx)
		if err != nil {
			log.Error("failed to fetch message", zap.Error(err))
			return
		}

		// 2. unmarshal raw bytes message into structure format
		var message CanalJSONMessageWithExtension
		err = json.Unmarshal(rawMessage.Value, &message)
		if err != nil {
			log.Error("unmarshal kafka message failed", zap.Error(err))
			break
		}

		if message.IsDDL {
			log.Info("DDL message received, handle it",
				zap.String("DDL", message.Query))
			
		}

		//switch message.EventType {
		//case "TiDB_WATERMARK":
		//case "DDL":
		//default:
		//
		//}

	}
}
