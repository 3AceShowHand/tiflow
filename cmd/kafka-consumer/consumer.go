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
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/pingcap/log"
	"github.com/pingcap/tiflow/pkg/errors"
	"go.uber.org/zap"
	"strings"
	"time"
)

func getPartitionNum(o *option) (int32, error) {
	configMap := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(o.address, ","),
	}
	if len(o.ca) != 0 {
		_ = configMap.SetKey("security.protocol", "SSL")
		_ = configMap.SetKey("ssl.ca.location", o.ca)
		_ = configMap.SetKey("ssl.key.location", o.key)
		_ = configMap.SetKey("ssl.certificate.location", o.cert)
	}
	admin, err := kafka.NewAdminClient(configMap)
	if err != nil {
		return 0, errors.Trace(err)
	}
	defer admin.Close()

	timeout := 3000
	for i := 0; i <= o.retryTime; i++ {
		resp, err := admin.GetMetadata(&o.topic, false, timeout)
		if err != nil {
			if err.(kafka.Error).Code() == kafka.ErrTransport {
				log.Info("retry get partition number", zap.Int("retryTime", i), zap.Int("timeout", timeout))
				timeout += 100
				continue
			}
			return 0, errors.Trace(err)
		}
		if topicDetail, ok := resp.Topics[o.topic]; ok {
			numPartitions := int32(len(topicDetail.Partitions))
			log.Info("get partition number of topic",
				zap.String("topic", o.topic),
				zap.Int32("partitionNum", numPartitions))
			return numPartitions, nil
		}
		log.Info("retry get partition number", zap.String("topic", o.topic))
		time.Sleep(1 * time.Second)
	}
	return 0, errors.Errorf("get partition number(%s) timeout", o.topic)
}

type consumer struct {
	option *option
	writer *writer
}

// NewConfluentConsumer will create a consumer client.
func NewConfluentConsumer(ctx context.Context, o *option) *consumer {
	partitionNum, err := getPartitionNum(o)
	if err != nil {
		log.Panic("cannot get the partition number", zap.String("topic", o.topic), zap.Error(err))
	}
	if o.partitionNum == 0 {
		o.partitionNum = partitionNum
	}
	w, err := NewWriter(ctx, o)
	if err != nil {
		log.Panic("cannot create the writer", zap.Error(err))
	}
	return &consumer{
		writer: w,
		option: o,
	}
}

// Consume will read message from Kafka.
func (c *consumer) Consume(ctx context.Context) error {
	topics := strings.Split(c.option.topic, ",")
	if len(topics) == 0 {
		log.Panic("no topic provided for the consumer")
	}
	configMap := &kafka.ConfigMap{
		"bootstrap.servers":  strings.Join(c.option.address, ","),
		"group.id":           c.option.groupID,
		"session.timeout.ms": 6000,
		// Start reading from the first message of each assigned
		// partition if there are no previously committed offsets
		// for this group.
		"auto.offset.reset": "earliest",
		// Whether we store offsets automatically.
		"enable.auto.offset.store": false,
		"enable.auto.commit":       false,
	}
	if len(c.option.ca) != 0 {
		_ = configMap.SetKey("security.protocol", "SSL")
		_ = configMap.SetKey("ssl.ca.location", c.option.ca)
		_ = configMap.SetKey("ssl.key.location", c.option.key)
		_ = configMap.SetKey("ssl.certificate.location", c.option.cert)
	}
	client, err := kafka.NewConsumer(configMap)
	if err != nil {
		log.Panic("Error creating consumer group client", zap.Error(err))
	}
	defer func() {
		if err = client.Close(); err != nil {
			log.Panic("Error closing client", zap.Error(err))
		}
	}()

	err = client.SubscribeTopics(topics, nil)
	for {
		msg, err := client.ReadMessage(100 * time.Millisecond)
		if err != nil {
			// Timeout is not considered an error because it is raised by
			// ReadMessage in absence of messages.
			if err.(kafka.Error).IsTimeout() {
				continue
			}
			return errors.Trace(err)
		}
		// Process the message received.
		partition := msg.TopicPartition.Partition
		needCommit, err := c.writer.Decode(ctx, c.option, partition, msg.Key, msg.Value)
		if err != nil {
			log.Panic("Error decode message", zap.Error(err))
		}
		if needCommit {
			// TODO: retry commit if fail
			if _, err := client.CommitMessage(msg); err != nil {
				log.Error("Error commit message", zap.Error(err))
				return errors.Trace(err)
			}
		}
	}
}
