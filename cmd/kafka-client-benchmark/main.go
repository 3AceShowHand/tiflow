package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pingcap/log"
	"github.com/pingcap/tiflow/cdc/model"
	"github.com/pingcap/tiflow/pkg/sink/codec/common"
	"github.com/pingcap/tiflow/pkg/sink/kafka"
	v2 "github.com/pingcap/tiflow/pkg/sink/kafka/v2"
	"go.uber.org/zap"
)

func main() {
	changefeed := model.DefaultChangeFeedID("test")
	option := kafka.NewOptions()
	option.BrokerEndpoints = []string{"127.0.0.1:9092"}
	option.ClientID = "kafka-client"

	factory, err := v2.NewFactory(option, changefeed)
	//	factory, err := kafka.NewSaramaFactory(option, changefeed)
	if err != nil {
		log.Error("create kafka factory failed", zap.Error(err))
		return
	}

	asyncProducer, err := factory.AsyncProducer(make(chan struct{}), make(chan error, 1))
	if err != nil {
		log.Error("create kafka async producer failed", zap.Error(err))
		return
	}
	defer asyncProducer.Close()

	var key []byte
	value := make([]byte, 10240)
	topic := "kafka-client-benchmark"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		total    uint64
		ackTotal uint64
		wg       sync.WaitGroup
	)

	events := make(chan []*common.Message, 1024)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			rand.Read(value)
			messages := make([]*common.Message, 0, 1024)
			for i := 0; i < 1024; i++ {
				messages = append(messages, &common.Message{
					Key:   key,
					Value: value,
					Callback: func() {
						atomic.AddUint64(&ackTotal, 1)
					},
				})
			}
			select {
			case <-ctx.Done():
				return
			case events <- messages:
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case messages := <-events:
				err := asyncProducer.AsyncSendMessages(ctx, topic, 0, messages)
				if err != nil {
					log.Error("send kafka message failed",
						zap.Int("count", len(messages)),
						zap.Error(err))
					return
				} else {
					atomic.AddUint64(&total, uint64(len(messages)))
				}
			}
		}
	}()

	wg.Add(1)
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer func() {
			ticker.Stop()
			defer wg.Done()
		}()
		old := uint64(0)
		oldAck := uint64(0)
		for {
			select {
			case <-ticker.C:
				temp := atomic.LoadUint64(&total)
				qps := (float64(temp) - float64(old)) / 5.0
				old = temp
				temp = atomic.LoadUint64(&ackTotal)
				ackQps := (float64(temp) - float64(oldAck)) / 5.0
				fmt.Printf("total %d, qps is %f, ack qps is %f\n", total, qps, ackQps)
				oldAck = temp
			}
		}
	}()

	if err := asyncProducer.AsyncRunCallback(ctx); err != nil {
		log.Error("run kafka async producer failed", zap.Error(err))
	}
}
