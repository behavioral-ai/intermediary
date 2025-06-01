package representation1

import (
	"github.com/behavioral-ai/core/fmtx"
	"time"
)

const (
	Fragment     = "v1"
	logRouteName = "app"

	AppHostKey  = "app-host"
	LogKey      = "log"
	LogRouteKey = "route-name"
	TimeoutKey  = "timeout"

	defaultTimeout = time.Millisecond * 2500
)

type Routing struct {
	Log          bool
	AppHost      string // User requirement
	LogRouteName string
	Timeout      time.Duration
}

func Initialize(m map[string]string) *Routing {
	r := new(Routing)
	r.Log = true
	r.LogRouteName = logRouteName
	r.Timeout = defaultTimeout
	parseRouting(r, m)
	return r
}

/*
func NewRouting(name string) *Routing {
	//m, _ := resource.Resolve[map[string]string](name, Fragment, resource.Resolver)
	return newRouting(nil)
}


*/

func newRouting(m map[string]string) *Routing {
	c := Initialize(m)
	return c
}

func (r *Routing) Update(m map[string]string) {
	parseRouting(r, m)
}

func parseRouting(r *Routing, m map[string]string) {
	if r == nil || m == nil {
		return
	}
	s := m[LogKey]
	if s != "" {
		if s == "true" {
			r.Log = true
		} else {
			r.Log = false
		}
	}
	s = m[LogRouteKey]
	if s != "" {
		r.LogRouteName = s
	}
	s = m[AppHostKey]
	if s != "" {
		r.AppHost = s
	}
	s = m[TimeoutKey]
	if s != "" {
		dur, err := fmtx.ParseDuration(s)
		if err != nil {
			//messaging.Reply(m, messaging.ConfigContentStatusError(agent, TimeoutKey), agent.Name())
			return
		}
		r.Timeout = dur
	}
}
