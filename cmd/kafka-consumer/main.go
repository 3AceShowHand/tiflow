// Copyright 2020 PingCAP, Inc.
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
	"flag"
	"fmt"
	"math"
	"net/url"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
	"github.com/google/uuid"
	"github.com/pingcap/errors"
	"github.com/pingcap/log"
	"github.com/pingcap/tiflow/cdc/model"
	"github.com/pingcap/tiflow/cdc/sink/ddlsink"
	ddlsinkfactory "github.com/pingcap/tiflow/cdc/sink/ddlsink/factory"
	eventsinkfactory "github.com/pingcap/tiflow/cdc/sink/dmlsink/factory"
	"github.com/pingcap/tiflow/cdc/sink/dmlsink/mq/dispatcher"
	"github.com/pingcap/tiflow/cdc/sink/tablesink"
	cmdUtil "github.com/pingcap/tiflow/pkg/cmd/util"
	"github.com/pingcap/tiflow/pkg/config"
	"github.com/pingcap/tiflow/pkg/filter"
	"github.com/pingcap/tiflow/pkg/logutil"
	"github.com/pingcap/tiflow/pkg/quotes"
	"github.com/pingcap/tiflow/pkg/security"
	"github.com/pingcap/tiflow/pkg/sink/codec"
	"github.com/pingcap/tiflow/pkg/sink/codec/avro"
	"github.com/pingcap/tiflow/pkg/sink/codec/canal"
	"github.com/pingcap/tiflow/pkg/sink/codec/common"
	"github.com/pingcap/tiflow/pkg/sink/codec/open"
	"github.com/pingcap/tiflow/pkg/spanz"
	"github.com/pingcap/tiflow/pkg/util"
	"github.com/pingcap/tiflow/pkg/version"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

// Sarama configuration options
var (
	kafkaAddrs           []string
	kafkaTopic           string
	kafkaPartitionNum    int32
	kafkaGroupID         = fmt.Sprintf("ticdc_kafka_consumer_%s", uuid.New().String())
	kafkaVersion         = "2.4.0"
	kafkaMaxMessageBytes = math.MaxInt64
	kafkaMaxBatchSize    = math.MaxInt64

	downstreamURIStr string

	protocol            config.Protocol
	enableTiDBExtension bool
	enableRowChecksum   bool

	// eventRouterReplicaConfig only used to initialize the consumer's eventRouter
	// which then can be used to check RowChangedEvent dispatched correctness
	eventRouterReplicaConfig *config.ReplicaConfig

	logPath       string
	logLevel      string
	timezone      string
	ca, cert, key string

	// avro schema registry uri should be set if the encoding protocol is avro
	schemaRegistryURI string
)

func init() {
	version.LogVersionInfo("kafka consumer")
	var (
		upstreamURIStr string
		configFile     string
	)

	flag.StringVar(&upstreamURIStr, "upstream-uri", "", "Kafka uri")
	flag.StringVar(&downstreamURIStr, "downstream-uri", "", "downstream sink uri")
	flag.StringVar(&schemaRegistryURI, "schema-registry-uri", "", "schema registry uri")
	flag.StringVar(&configFile, "config", "", "config file for changefeed")
	flag.StringVar(&logPath, "log-file", "cdc_kafka_consumer.log", "log file path")
	flag.StringVar(&logLevel, "log-level", "info", "log file path")
	flag.StringVar(&timezone, "tz", "System", "Specify time zone of Kafka consumer")
	flag.StringVar(&ca, "ca", "", "CA certificate path for Kafka SSL connection")
	flag.StringVar(&cert, "cert", "", "Certificate path for Kafka SSL connection")
	flag.StringVar(&key, "key", "", "Private key path for Kafka SSL connection")
	flag.Parse()

	err := logutil.InitLogger(&logutil.Config{
		Level: logLevel,
		File:  logPath,
	},
		logutil.WithInitGRPCLogger(),
		logutil.WithInitSaramaLogger(),
	)
	if err != nil {
		log.Panic("init logger failed", zap.Error(err))
	}

	upstreamURI, err := url.Parse(upstreamURIStr)
	if err != nil {
		log.Panic("invalid upstream-uri", zap.Error(err))
	}
	scheme := strings.ToLower(upstreamURI.Scheme)
	if scheme != "kafka" {
		log.Panic("invalid upstream-uri scheme, the scheme of upstream-uri must be `kafka`",
			zap.String("upstreamURI", upstreamURIStr))
	}
	s := upstreamURI.Query().Get("version")
	if s != "" {
		kafkaVersion = s
	}
	s = upstreamURI.Query().Get("consumer-group-id")
	if s != "" {
		kafkaGroupID = s
	}
	kafkaTopic = strings.TrimFunc(upstreamURI.Path, func(r rune) bool {
		return r == '/'
	})
	kafkaAddrs = strings.Split(upstreamURI.Host, ",")

	saramaConfig, err := newSaramaConfig()
	if err != nil {
		log.Panic("Error creating sarama saramaConfig", zap.Error(err))
	}

	s = upstreamURI.Query().Get("partition-num")
	if s == "" {
		partition, err := getPartitionNum(kafkaAddrs, kafkaTopic, saramaConfig)
		if err != nil {
			log.Panic("can not get partition number", zap.String("topic", kafkaTopic), zap.Error(err))
		}
		kafkaPartitionNum = partition
	} else {
		c, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			log.Panic("invalid partition-num of upstream-uri")
		}
		kafkaPartitionNum = int32(c)
	}
	log.Info("Setting partitionNum", zap.Int32("partitionNum", kafkaPartitionNum))

	s = upstreamURI.Query().Get("max-message-bytes")
	if s != "" {
		c, err := strconv.Atoi(s)
		if err != nil {
			log.Panic("invalid max-message-bytes of upstream-uri")
		}
		kafkaMaxMessageBytes = c
	}
	log.Info("Setting max-message-bytes", zap.Int("max-message-bytes", kafkaMaxMessageBytes))

	s = upstreamURI.Query().Get("max-batch-size")
	if s != "" {
		c, err := strconv.Atoi(s)
		if err != nil {
			log.Panic("invalid max-batch-size of upstream-uri")
		}
		kafkaMaxBatchSize = c
	}
	log.Info("Setting max-batch-size", zap.Int("max-batch-size", kafkaMaxBatchSize))

	s = upstreamURI.Query().Get("protocol")
	if s != "" {
		if protocol, err = config.ParseSinkProtocolFromString(s); err != nil {
			log.Panic("invalid protocol", zap.Error(err), zap.String("protocol", s))
		}
	}
	log.Info("Setting protocol", zap.Any("protocol", protocol))

	s = upstreamURI.Query().Get("enable-tidb-extension")
	if s != "" {
		enableTiDBExtension, err = strconv.ParseBool(s)
		if err != nil {
			log.Panic("invalid enable-tidb-extension of upstream-uri")
		}
		if enableTiDBExtension {
			if protocol != config.ProtocolCanalJSON && protocol != config.ProtocolAvro {
				log.Panic("enable-tidb-extension only work with canal-json / avro")
			}
		}
	}

	s = upstreamURI.Query().Get("enable-row-checksum")
	if s != "" {
		enableRowChecksum, err = strconv.ParseBool(s)
		if err != nil {
			log.Panic("invalid enable-row-checksum of upstream-uri")
		}
		if enableRowChecksum {
			if protocol != config.ProtocolAvro {
				log.Panic("enable-row-checksum only work with avro")
			}
		}
	}

	if configFile != "" {
		eventRouterReplicaConfig = config.GetDefaultReplicaConfig()
		eventRouterReplicaConfig.Sink.Protocol = util.AddressOf(protocol.String())
		err := cmdUtil.StrictDecodeFile(configFile, "kafka consumer", eventRouterReplicaConfig)
		if err != nil {
			log.Panic("invalid config file for kafka consumer",
				zap.Error(err),
				zap.String("config", configFile))
		}
		if _, err := filter.VerifyTableRules(eventRouterReplicaConfig.Filter); err != nil {
			log.Panic("verify rule failed", zap.Error(err))
		}
	}
}

func getPartitionNum(address []string, topic string, cfg *sarama.Config) (int32, error) {
	// get partition number or create topic automatically
	admin, err := sarama.NewClusterAdmin(address, cfg)
	if err != nil {
		return 0, errors.Trace(err)
	}
	topics, err := admin.ListTopics()
	if err != nil {
		return 0, errors.Trace(err)
	}
	err = admin.Close()
	if err != nil {
		return 0, errors.Trace(err)
	}
	topicDetail, exist := topics[topic]
	if !exist {
		return 0, errors.Errorf("can not find topic %s", topic)
	}
	log.Info("get partition number of topic",
		zap.String("topic", topic),
		zap.Int32("partitionNum", topicDetail.NumPartitions))
	return topicDetail.NumPartitions, nil
}

func waitTopicCreated(address []string, topic string, cfg *sarama.Config) error {
	admin, err := sarama.NewClusterAdmin(address, cfg)
	if err != nil {
		return errors.Trace(err)
	}
	defer admin.Close()
	for i := 0; i <= 30; i++ {
		topics, err := admin.ListTopics()
		if err != nil {
			return errors.Trace(err)
		}
		if _, ok := topics[topic]; ok {
			return nil
		}
		log.Info("wait the topic created", zap.String("topic", topic))
		time.Sleep(1 * time.Second)
	}
	return errors.Errorf("wait the topic(%s) created timeout", topic)
}

func newSaramaConfig() (*sarama.Config, error) {
	config := sarama.NewConfig()

	version, err := sarama.ParseKafkaVersion(kafkaVersion)
	if err != nil {
		return nil, errors.Trace(err)
	}

	config.ClientID = "ticdc_kafka_sarama_consumer"
	config.Version = version

	config.Metadata.Retry.Max = 10000
	config.Metadata.Retry.Backoff = 500 * time.Millisecond
	config.Consumer.Retry.Backoff = 500 * time.Millisecond
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	if len(ca) != 0 {
		config.Net.TLS.Enable = true
		config.Net.TLS.Config, err = (&security.Credential{
			CAPath:   ca,
			CertPath: cert,
			KeyPath:  key,
		}).ToTLSConfig()
		if err != nil {
			return nil, errors.Trace(err)
		}
	}

	return config, err
}

func main() {
	/**
	 * Construct a new Sarama configuration.
	 * The Kafka cluster version has to be defined before the consumer/producer is initialized.
	 */
	config, err := newSaramaConfig()
	if err != nil {
		log.Panic("Error creating sarama config", zap.Error(err))
	}
	err = waitTopicCreated(kafkaAddrs, kafkaTopic, config)
	if err != nil {
		log.Panic("wait topic created failed", zap.Error(err))
	}
	/**
	 * Setup a new Sarama consumer group
	 */
	log.Info("Starting a new TiCDC consumer",
		zap.String("GroupID", kafkaGroupID), zap.Any("protocol", protocol))

	ctx, cancel := context.WithCancel(context.Background())
	consumer, err := NewConsumer(ctx)
	if err != nil {
		log.Panic("Error creating consumer", zap.Error(err))
	}

	client, err := sarama.NewConsumerGroup(kafkaAddrs, kafkaGroupID, config)
	if err != nil {
		log.Panic("Error creating consumer group client", zap.Error(err))
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := client.Consume(ctx, strings.Split(kafkaTopic, ","), consumer); err != nil {
				log.Panic("Error from consumer: %v", zap.Error(err))
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	go func() {
		if err := consumer.Run(ctx); err != nil {
			if err != context.Canceled {
				log.Panic("Error running consumer", zap.Error(err))
			}
		}
	}()

	<-consumer.ready // wait till the consumer has been set up
	log.Info("TiCDC consumer up and running!...")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		log.Info("terminating: context cancelled")
	case <-sigterm:
		log.Info("terminating: via signal")
	}
	cancel()
	wg.Wait()
	if err = client.Close(); err != nil {
		log.Panic("Error closing client", zap.Error(err))
	}
}

type partitionSinks struct {
	tablesCommitTsMap sync.Map
	tableSinksMap     sync.Map
	resolvedTs        uint64
	partitionNo       int
}

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	ready chan bool

	ddlList              []*model.DDLEvent
	ddlListMu            sync.Mutex
	ddlWithMaxCommitTs   *model.DDLEvent
	ddlSink              ddlsink.Sink
	fakeTableIDGenerator *fakeTableIDGenerator

	// sinkFactory is used to create table sink for each table.
	sinkFactory *eventsinkfactory.SinkFactory
	sinks       []*partitionSinks
	sinksMu     sync.Mutex

	// initialize to 0 by default
	globalResolvedTs uint64

	protocol            config.Protocol
	enableTiDBExtension bool
	enableRowChecksum   bool

	eventRouter *dispatcher.EventRouter

	// avro only
	// key and value schema manager, only used by the avro protocol to fetch schema.
	keySchemaM   *avro.SchemaManager
	valueSchemaM *avro.SchemaManager

	tz *time.Location
}

// NewConsumer creates a new cdc kafka consumer
func NewConsumer(ctx context.Context) (*Consumer, error) {
	// TODO support filter in downstream sink
	tz, err := util.GetTimezone(timezone)
	if err != nil {
		return nil, errors.Annotate(err, "can not load timezone")
	}
	config.GetGlobalServerConfig().TZ = timezone

	c := new(Consumer)
	c.tz = tz
	c.fakeTableIDGenerator = &fakeTableIDGenerator{
		tableIDs: make(map[string]int64),
	}
	c.protocol = protocol
	c.enableTiDBExtension = enableTiDBExtension
	c.enableRowChecksum = enableRowChecksum

	if c.protocol == config.ProtocolAvro {
		keySchemaM, valueSchemaM, err := avro.NewKeyAndValueSchemaManagers(
			ctx, schemaRegistryURI, nil)
		if err != nil {
			return nil, errors.Trace(err)
		}
		c.keySchemaM = keySchemaM
		c.valueSchemaM = valueSchemaM
	}

	// this means user has input config file to enable dispatcher check
	// some protocol does not provide enough information to check the
	// dispatched partition match or not. such as `open-protocol`, which
	// does not have `IndexColumn` info, then make the default dispatcher
	// use different dispatch rule to the CDC side.
	// when try to enable dispatcher check for any protocol and dispatch
	// rule, make sure decoded `RowChangedEvent` contains information
	// identical to the CDC side.
	if eventRouterReplicaConfig != nil {
		eventRouter, err := dispatcher.NewEventRouter(eventRouterReplicaConfig, kafkaTopic)
		if err != nil {
			return nil, errors.Trace(err)
		}
		c.eventRouter = eventRouter

	}

	c.sinks = make([]*partitionSinks, kafkaPartitionNum)
	ctx, cancel := context.WithCancel(ctx)
	errChan := make(chan error, 1)
	for i := 0; i < int(kafkaPartitionNum); i++ {
		c.sinks[i] = &partitionSinks{
			partitionNo: i,
		}
	}
	f, err := eventsinkfactory.New(
		ctx,
		model.DefaultChangeFeedID("test"),
		downstreamURIStr,
		config.GetDefaultReplicaConfig(),
		errChan,
	)
	if err != nil {
		cancel()
		return nil, errors.Trace(err)
	}
	c.sinkFactory = f

	go func() {
		err := <-errChan
		if errors.Cause(err) != context.Canceled {
			log.Error("error on running consumer", zap.Error(err))
		} else {
			log.Info("consumer exited")
		}
		cancel()
	}()

	ddlSink, err := ddlsinkfactory.New(
		ctx,
		model.DefaultChangeFeedID("test"),
		downstreamURIStr,
		config.GetDefaultReplicaConfig(),
	)
	if err != nil {
		cancel()
		return nil, errors.Trace(err)
	}
	c.ddlSink = ddlSink
	c.ready = make(chan bool)
	return c, nil
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the c as ready
	close(c.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

type eventsGroup struct {
	events []*model.RowChangedEvent
}

func newEventsGroup() *eventsGroup {
	return &eventsGroup{
		events: make([]*model.RowChangedEvent, 0),
	}
}

func (g *eventsGroup) Append(e *model.RowChangedEvent) {
	g.events = append(g.events, e)
}

func (g *eventsGroup) Resolve(resolveTs uint64) []*model.RowChangedEvent {
	sort.Slice(g.events, func(i, j int) bool {
		return g.events[i].CommitTs < g.events[j].CommitTs
	})

	i := sort.Search(len(g.events), func(i int) bool {
		return g.events[i].CommitTs > resolveTs
	})
	result := g.events[:i]
	g.events = g.events[i:]

	return result
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	partition := claim.Partition()
	c.sinksMu.Lock()
	sink := c.sinks[partition]
	c.sinksMu.Unlock()
	if sink == nil {
		panic("sink should initialized")
	}

	var (
		decoder codec.RowEventDecoder
		err     error
	)
	switch c.protocol {
	case config.ProtocolOpen, config.ProtocolDefault:
		decoder = open.NewBatchDecoder()
	case config.ProtocolCanalJSON:
		decoder = canal.NewBatchDecoder(c.enableTiDBExtension, "")
	case config.ProtocolAvro:
		config := &common.Config{
			EnableTiDBExtension: c.enableTiDBExtension,
			EnableRowChecksum:   c.enableRowChecksum,
			// avro must set this to true to make the consumer works.
			AvroEnableWatermark: true,
		}
		decoder = avro.NewDecoder(config, c.keySchemaM, c.valueSchemaM, kafkaTopic, c.tz)
	default:
		log.Panic("Protocol not supported", zap.Any("Protocol", c.protocol))
	}
	if err != nil {
		return errors.Trace(err)
	}

	eventGroups := make(map[int64]*eventsGroup)
	for message := range claim.Messages() {
		if err := decoder.AddKeyValue(message.Key, message.Value); err != nil {
			log.Error("add key value to the decoder failed", zap.Error(err))
			return errors.Trace(err)
		}

		counter := 0
		for {
			tp, hasNext, err := decoder.HasNext()
			if err != nil {
				log.Panic("decode message key failed", zap.Error(err))
			}
			if !hasNext {
				break
			}

			counter++
			// If the message containing only one event exceeds the length limit, CDC will allow it and issue a warning.
			if len(message.Key)+len(message.Value) > kafkaMaxMessageBytes && counter > 1 {
				log.Panic("kafka max-messages-bytes exceeded", zap.Int("max-message-bytes", kafkaMaxMessageBytes),
					zap.Int("receivedBytes", len(message.Key)+len(message.Value)))
			}

			switch tp {
			case model.MessageTypeDDL:
				// for some protocol, DDL would be dispatched to all partitions,
				// Consider that DDL a, b, c received from partition-0, the latest DDL is c,
				// if we receive `a` from partition-1, which would be seemed as DDL regression,
				// then cause the consumer panic, but it was a duplicate one.
				// so we only handle DDL received from partition-0 should be enough.
				// but all DDL event messages should be consumed.
				ddl, err := decoder.NextDDLEvent()
				if err != nil {
					log.Panic("decode message value failed",
						zap.ByteString("value", message.Value),
						zap.Error(err))
				}
				if partition == 0 {
					c.appendDDL(ddl)
				}
			case model.MessageTypeRow:
				row, _, err := decoder.NextRowChangedEvent()
				if err != nil {
					log.Panic("decode message value failed",
						zap.ByteString("value", message.Value),
						zap.Error(err))
				}

				if c.eventRouter != nil {
					target := c.eventRouter.GetPartitionForRowChange(row, kafkaPartitionNum)
					if partition != target {
						log.Panic("RowChangedEvent dispatched to wrong partition",
							zap.Int32("obtained", partition),
							zap.Int32("expected", target),
							zap.Int32("partitionNum", kafkaPartitionNum),
							zap.Any("row", row),
						)
					}
				}

				globalResolvedTs := atomic.LoadUint64(&c.globalResolvedTs)
				if row.CommitTs <= globalResolvedTs || row.CommitTs <= sink.resolvedTs {
					log.Warn("RowChangedEvent fallback row, ignore it",
						zap.Uint64("commitTs", row.CommitTs),
						zap.Uint64("globalResolvedTs", globalResolvedTs),
						zap.Uint64("sinkResolvedTs", sink.resolvedTs),
						zap.Int32("partition", partition),
						zap.Any("row", row))
				}
				var partitionID int64
				if row.Table.IsPartition {
					partitionID = row.Table.TableID
				}
				tableID := c.fakeTableIDGenerator.
					generateFakeTableID(row.Table.Schema, row.Table.Table, partitionID)
				row.Table.TableID = tableID

				group, ok := eventGroups[tableID]
				if !ok {
					group = newEventsGroup()
					eventGroups[tableID] = group
				}
				group.Append(row)
			case model.MessageTypeResolved:
				ts, err := decoder.NextResolvedEvent()
				if err != nil {
					log.Panic("decode message value failed",
						zap.ByteString("value", message.Value),
						zap.Error(err))
				}
				resolvedTs := atomic.LoadUint64(&sink.resolvedTs)
				// `resolvedTs` should be monotonically increasing, it's allowed to receive redundant one.
				if ts < resolvedTs {
					log.Panic("partition resolved ts fallback",
						zap.Uint64("ts", ts),
						zap.Uint64("resolvedTs", resolvedTs),
						zap.Int32("partition", partition))
				}
				if ts > resolvedTs {
					for tableID, group := range eventGroups {
						events := group.Resolve(ts)
						if len(events) == 0 {
							continue
						}
						if _, ok := sink.tableSinksMap.Load(tableID); !ok {
							sink.tableSinksMap.Store(tableID, c.sinkFactory.CreateTableSinkForConsumer(
								model.DefaultChangeFeedID("kafka-consumer"),
								spanz.TableIDToComparableSpan(tableID),
								events[0].CommitTs,
								prometheus.NewCounter(prometheus.CounterOpts{}),
							))
						}
						s, _ := sink.tableSinksMap.Load(tableID)
						s.(tablesink.TableSink).AppendRowChangedEvents(events...)
						commitTs := events[len(events)-1].CommitTs
						lastCommitTs, ok := sink.tablesCommitTsMap.Load(tableID)
						if !ok || lastCommitTs.(uint64) < commitTs {
							sink.tablesCommitTsMap.Store(tableID, commitTs)
						}
					}
					log.Debug("update sink resolved ts",
						zap.Uint64("ts", ts),
						zap.Int32("partition", partition))
					atomic.StoreUint64(&sink.resolvedTs, ts)
				} else {
					log.Info("redundant sink resolved ts", zap.Uint64("ts", ts), zap.Int32("partition", partition))
				}
			}
			session.MarkMessage(message, "")
		}

		if counter > kafkaMaxBatchSize {
			log.Panic("Open Protocol max-batch-size exceeded", zap.Int("max-batch-size", kafkaMaxBatchSize),
				zap.Int("actual-batch-size", counter))
		}
	}

	return nil
}

// append DDL wait to be handled, only consider the constraint among DDLs.
// for DDL a / b received in the order, a.CommitTs < b.CommitTs should be true.
func (c *Consumer) appendDDL(ddl *model.DDLEvent) {
	c.ddlListMu.Lock()
	defer c.ddlListMu.Unlock()
	// DDL CommitTs fallback, just crash it to indicate the bug.
	if c.ddlWithMaxCommitTs != nil && ddl.CommitTs < c.ddlWithMaxCommitTs.CommitTs {
		log.Panic("DDL CommitTs < maxCommitTsDDL.CommitTs",
			zap.Uint64("commitTs", ddl.CommitTs),
			zap.Uint64("maxCommitTs", c.ddlWithMaxCommitTs.CommitTs),
			zap.Any("DDL", ddl))
	}

	// A rename tables DDL job contains multiple DDL events with same CommitTs.
	// So to tell if a DDL is redundant or not, we must check the equivalence of
	// the current DDL and the DDL with max CommitTs.
	if ddl == c.ddlWithMaxCommitTs {
		log.Info("ignore redundant DDL, the DDL is equal to ddlWithMaxCommitTs",
			zap.Any("DDL", ddl))
		return
	}

	c.ddlList = append(c.ddlList, ddl)
	log.Info("DDL event received", zap.Any("DDL", ddl))
	c.ddlWithMaxCommitTs = ddl
}

func (c *Consumer) getFrontDDL() *model.DDLEvent {
	c.ddlListMu.Lock()
	defer c.ddlListMu.Unlock()
	if len(c.ddlList) > 0 {
		return c.ddlList[0]
	}
	return nil
}

func (c *Consumer) popDDL() *model.DDLEvent {
	c.ddlListMu.Lock()
	defer c.ddlListMu.Unlock()
	if len(c.ddlList) > 0 {
		ddl := c.ddlList[0]
		c.ddlList = c.ddlList[1:]
		return ddl
	}
	return nil
}

func (c *Consumer) forEachSink(fn func(sink *partitionSinks) error) error {
	c.sinksMu.Lock()
	defer c.sinksMu.Unlock()
	for _, sink := range c.sinks {
		if err := fn(sink); err != nil {
			return errors.Trace(err)
		}
	}
	return nil
}

func (c *Consumer) getMinPartitionResolvedTs() (result uint64, err error) {
	result = uint64(math.MaxUint64)
	err = c.forEachSink(func(sink *partitionSinks) error {
		a := atomic.LoadUint64(&sink.resolvedTs)
		if a < result {
			result = a
		}
		return nil
	})
	return result, err
}

// Run the Consumer
func (c *Consumer) Run(ctx context.Context) error {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}

		minPartitionResolvedTs, err := c.getMinPartitionResolvedTs()
		if err != nil {
			return errors.Trace(err)
		}

		// handle DDL
		todoDDL := c.getFrontDDL()
		if todoDDL != nil && todoDDL.CommitTs <= minPartitionResolvedTs {
			// flush DMLs
			if err := c.forEachSink(func(sink *partitionSinks) error {
				return syncFlushRowChangedEvents(ctx, sink, todoDDL.CommitTs)
			}); err != nil {
				return errors.Trace(err)
			}

			// DDL can be executed, do it first.
			if err := c.ddlSink.WriteDDLEvent(ctx, todoDDL); err != nil {
				return errors.Trace(err)
			}
			c.popDDL()

			if todoDDL.CommitTs < minPartitionResolvedTs {
				log.Info("update minPartitionResolvedTs by DDL",
					zap.Uint64("minPartitionResolvedTs", minPartitionResolvedTs),
					zap.Any("DDL", todoDDL))
			}
			minPartitionResolvedTs = todoDDL.CommitTs
		}

		// update global resolved ts
		if c.globalResolvedTs > minPartitionResolvedTs {
			log.Panic("global ResolvedTs fallback",
				zap.Uint64("globalResolvedTs", c.globalResolvedTs),
				zap.Uint64("minPartitionResolvedTs", minPartitionResolvedTs))
		}

		if c.globalResolvedTs == minPartitionResolvedTs {
			continue
		}

		c.globalResolvedTs = minPartitionResolvedTs

		if err := c.forEachSink(func(sink *partitionSinks) error {
			return syncFlushRowChangedEvents(ctx, sink, c.globalResolvedTs)
		}); err != nil {
			return errors.Trace(err)
		}
	}
}

func syncFlushRowChangedEvents(ctx context.Context, sink *partitionSinks, resolvedTs uint64) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		flushedResolvedTs := true
		sink.tablesCommitTsMap.Range(func(key, value interface{}) bool {
			tableID := key.(int64)
			resolvedTs := model.NewResolvedTs(resolvedTs)
			tableSink, ok := sink.tableSinksMap.Load(tableID)
			if !ok {
				log.Panic("Table sink not found", zap.Int64("tableID", tableID))
			}
			if err := tableSink.(tablesink.TableSink).UpdateResolvedTs(resolvedTs); err != nil {
				log.Error("Failed to update resolved ts", zap.Error(err))
				return false
			}
			if !tableSink.(tablesink.TableSink).GetCheckpointTs().EqualOrGreater(resolvedTs) {
				flushedResolvedTs = false
			}
			return true
		})
		if flushedResolvedTs {
			return nil
		}
	}
}

type fakeTableIDGenerator struct {
	tableIDs       map[string]int64
	currentTableID int64
	mu             sync.Mutex
}

func (g *fakeTableIDGenerator) generateFakeTableID(schema, table string, partition int64) int64 {
	g.mu.Lock()
	defer g.mu.Unlock()
	key := quotes.QuoteSchema(schema, table)
	if partition != 0 {
		key = fmt.Sprintf("%s.`%d`", key, partition)
	}
	if tableID, ok := g.tableIDs[key]; ok {
		return tableID
	}
	g.currentTableID++
	g.tableIDs[key] = g.currentTableID
	return g.currentTableID
}
