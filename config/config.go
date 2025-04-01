package config

import (
	"github.com/behavioral-ai/core/fmtx"
	"github.com/behavioral-ai/core/messaging"
	"time"
)

const (
	AppHostKey   = "app-host"
	CacheHostKey = "cache-host"
	TimeoutKey   = "timeout"
)

func Timeout(agent messaging.Agent, m *messaging.Message) (time.Duration, bool) {
	cfg := messaging.ConfigMapContent(m)
	if cfg == nil {
		messaging.Reply(m, messaging.ConfigEmptyStatusError(agent), agent.Uri())
		return 0, false
	}
	timeout := cfg[TimeoutKey]
	if timeout == "" {
		messaging.Reply(m, messaging.ConfigContentStatusError(agent, TimeoutKey), agent.Uri())
		return 0, false
	}
	dur, err := fmtx.ParseDuration(timeout)
	if err != nil {
		messaging.Reply(m, messaging.ConfigContentStatusError(agent, TimeoutKey), agent.Uri())
		return 0, false
	}
	return dur, true
}

func AppHostName(agent messaging.Agent, m *messaging.Message) (string, bool) {
	return hostName(agent, m, AppHostKey)
}

func CacheHostName(agent messaging.Agent, m *messaging.Message) (string, bool) {
	return hostName(agent, m, CacheHostKey)
}

func hostName(agent messaging.Agent, m *messaging.Message, key string) (string, bool) {
	cfg := messaging.ConfigMapContent(m)
	if cfg == nil {
		messaging.Reply(m, messaging.ConfigEmptyStatusError(agent), agent.Uri())
		return "", false
	}
	host := cfg[key]
	if host == "" {
		messaging.Reply(m, messaging.ConfigContentStatusError(agent, key), agent.Uri())
		return "", false
	}
	return host, true
}
