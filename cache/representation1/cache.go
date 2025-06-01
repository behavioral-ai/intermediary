package representation1

import (
	"github.com/behavioral-ai/core/fmtx"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

const (
	Fragment        = "v1"
	HostKey         = "host"
	CacheControlKey = "cache-control"
	TimeoutKey      = "timeout"
	IntervalKey     = "interval"
	SundayKey       = "sun"
	MondayKey       = "mon"
	TuesdayKey      = "tue"
	WednesdayKey    = "wed"
	ThursdayKey     = "thu"
	FridayKey       = "fri"
	SaturdayKey     = "sat"
	rangeSeparator  = "-"

	defaultInterval = time.Minute * 30
	defaultTimeout  = time.Millisecond * 2000
)

type Cache struct {
	Running  bool
	Enabled  *atomic.Bool
	Timeout  time.Duration
	Interval time.Duration
	Host     string           // User requirement
	Policy   http.Header      // User requirement
	Days     map[string]Range // User requirement
}

// Initialize - add a default policy
func Initialize(m map[string]string) *Cache {
	c := new(Cache)
	c.Enabled = new(atomic.Bool)
	c.Enabled.Store(false)
	c.Timeout = defaultTimeout
	c.Interval = defaultInterval
	c.Policy = make(http.Header)
	c.Days = make(map[string]Range)
	parseCache(c, m)
	return c
}

/*
func NewCache(name string) *Cache {
	//m, _ := resource.Resolve[map[string]string](name, Fragment, resource.Resolver)
	return newCache(nil)
}
*/

func newCache(m map[string]string) *Cache {
	c := Initialize(m)
	return c
}

func (c *Cache) Now() bool {
	ts := time.Now()
	day := ts.Weekday()
	s := ""
	switch day {
	case 0:
		s = SundayKey
	case 1:
		s = MondayKey
	case 2:
		s = TuesdayKey
	case 3:
		s = WednesdayKey
	case 4:
		s = ThursdayKey
	case 5:
		s = FridayKey
	case 6:
		s = SaturdayKey
	}
	return c.Days[s].In(ts)
}

func (c *Cache) Update(m map[string]string) {
	parseCache(c, m)
}

func parseCache(c *Cache, m map[string]string) {
	if c == nil || m == nil {
		return
	}
	if c.Policy == nil {
		c.Policy = make(http.Header)
	}
	if c.Days == nil {
		c.Days = make(map[string]Range)
	}
	s := m[HostKey]
	if s != "" {
		c.Host = s
	}
	s = m[CacheControlKey]
	if s != "" {
		c.Policy.Set(CacheControlKey, s)
	}
	s = m[TimeoutKey]
	if s != "" {
		dur, err := fmtx.ParseDuration(s)
		if err != nil {
			//messaging.Reply(m, messaging.ConfigContentStatusError(agent, TimeoutKey), agent.Name())
			return
		}
		c.Timeout = dur
	}
	s = m[IntervalKey]
	if s != "" {
		dur, err := fmtx.ParseDuration(s)
		if err != nil {
			//messaging.Reply(m, messaging.ConfigContentStatusError(agent, TimeoutKey), agent.Name())
			return
		}
		c.Interval = dur
	}
	parseDays(c, m)
}

func parseDays(c *Cache, m map[string]string) {
	parseDay(c, SundayKey, m)
	parseDay(c, MondayKey, m)
	parseDay(c, TuesdayKey, m)
	parseDay(c, WednesdayKey, m)
	parseDay(c, ThursdayKey, m)
	parseDay(c, FridayKey, m)
	parseDay(c, SaturdayKey, m)
}

func parseDay(c *Cache, key string, m map[string]string) {
	s := m[key]
	if s == "" {
		return
	}
	r := NewRange(s)
	if !r.Empty() {
		c.Days[key] = r
	}
}

// Range - hour range
type Range struct {
	From int
	To   int
}

func NewRange(s string) Range {
	if s == "" {
		return Range{}
	}
	tokens := strings.Split(strings.Trim(s, " "), rangeSeparator)
	if len(tokens) != 2 {
		return Range{}
	}
	r := Range{}
	if i, err := strconv.Atoi(tokens[0]); err == nil {
		r.From = i
	}
	if i, err := strconv.Atoi(tokens[1]); err == nil {
		r.To = i
	}
	return r
}

func (r Range) Empty() bool {
	if r.From < 0 || r.To <= 0 {
		return true
	}
	if r.From > 23 || r.To > 23 {
		return false
	}
	return r.From > r.To
}

func (r Range) In(ts time.Time) bool {
	hour := ts.Hour()
	return r.From <= hour && hour <= r.To
}
