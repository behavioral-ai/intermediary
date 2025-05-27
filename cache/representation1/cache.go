package representation1

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	Fragment        = "v1"
	hostKey         = "host"
	cacheControlKey = "cache-control"
	sundayKey       = "sun"
	mondayKey       = "mon"
	tuesdayKey      = "tue"
	wednesdayKey    = "wed"
	thursdayKey     = "thu"
	fridayKey       = "fri"
	saturdayKey     = "sat"
	rangeSeparator  = "-"
)

type Cache struct {
	Running bool
	Host    string
	Policy  http.Header
	Days    map[string]Range
}

// Initialize - add a default policy
func Initialize() *Cache {
	c := new(Cache)
	c.Policy = make(http.Header)
	c.Policy.Add(cacheControlKey, "max-age=0")
	c.Days = make(map[string]Range)
	return c
}

func NewCache(name string) *Cache {
	c := Initialize()
	parseCache(c, nil)
	return c
}

func (c *Cache) Enabled() bool {
	return c.Host != ""
}

func (c *Cache) Now() bool {
	ts := time.Now().UTC()
	day := ts.Weekday()
	s := ""
	switch day {
	case 0:
		s = sundayKey
	case 1:
		s = mondayKey
	case 2:
		s = tuesdayKey
	case 3:
		s = wednesdayKey
	case 4:
		s = thursdayKey
	case 5:
		s = fridayKey
	case 6:
		s = saturdayKey
	}
	return c.Days[s].In(ts)
}

func parseCache(c *Cache, m map[string]string) {
	if c == nil || m == nil {
		return
	}
	c.Host = m[hostKey]
	c.Policy.Add(cacheControlKey, m[cacheControlKey])
	parseDays(c, m)

}

func parseDays(c *Cache, m map[string]string) {
	parseDay(c, sundayKey, m)
	parseDay(c, mondayKey, m)
	parseDay(c, tuesdayKey, m)
	parseDay(c, wednesdayKey, m)
	parseDay(c, thursdayKey, m)
	parseDay(c, fridayKey, m)
	parseDay(c, saturdayKey, m)
}

func parseDay(c *Cache, key string, m map[string]string) {
	r := NewRange(m[key])
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
	tokens := strings.Split(s, rangeSeparator)
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
	return r.From == 0 || r.To == 0
}

func (r Range) In(ts time.Time) bool {
	hour := ts.Hour()
	return r.From <= hour && hour <= r.To
}
