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

package agent

import "github.com/prometheus/client_golang/prometheus"

var (
	agentTickDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "ticdc",
			Subsystem: "scheduler",
			Name:      "agent_tick_duration",
			Help:      "Bucketed histogram of agent tick processor time (s).",
			Buckets:   prometheus.ExponentialBuckets(0.01 /* 10 ms */, 2, 18),
		}, []string{"namespace", "changefeed"})
	agentReceiveMessageDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "ticdc",
			Subsystem: "scheduler",
			Name:      "agent_receive_message_duration",
			Help:      "Bucketed histogram of agent tick processor time (s).",
			Buckets:   prometheus.ExponentialBuckets(0.01 /* 10 ms */, 2, 18),
		}, []string{"namespace", "changefeed"})
	agentHandleMessageDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "ticdc",
			Subsystem: "scheduler",
			Name:      "agent_handle_message_duration",
			Help:      "Bucketed histogram of agent tick processor time (s).",
			Buckets:   prometheus.ExponentialBuckets(0.01 /* 10 ms */, 2, 18),
		}, []string{"namespace", "changefeed"})
	agentPollDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "ticdc",
			Subsystem: "scheduler",
			Name:      "agent_poll_duration",
			Help:      "Bucketed histogram of agent tick processor time (s).",
			Buckets:   prometheus.ExponentialBuckets(0.01 /* 10 ms */, 2, 18),
		}, []string{"namespace", "changefeed"})
	agentSendMessageDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "ticdc",
			Subsystem: "scheduler",
			Name:      "agent_send_message_duration",
			Help:      "Bucketed histogram of agent tick processor time (s).",
			Buckets:   prometheus.ExponentialBuckets(0.01 /* 10 ms */, 2, 18),
		}, []string{"namespace", "changefeed"})
)

func InitMetrics(registry *prometheus.Registry) {
	registry.MustRegister(agentTickDuration)
	registry.MustRegister(agentReceiveMessageDuration)
	registry.MustRegister(agentHandleMessageDuration)
	registry.MustRegister(agentPollDuration)
	registry.MustRegister(agentSendMessageDuration)
}
