package profile

import "time"

/*
const (
	TrafficOffPeak   = "off-peak"
	TrafficPeak      = "peak"
	TrafficScaleUp   = "scale-up"
	TrafficScaleDown = "scale-down"
	trafficName      = "resiliency:type/traffic/profile/traffic"
)


*/

type Cache struct {
	Week [7][24]bool
}

func (c Cache) Now() bool {
	ts := time.Now().UTC()
	day := ts.Weekday()
	hour := ts.Hour()
	return c.Week[day][hour]
}

/*
func NewTraffic(curr *Traffic, resolver *content.Resolution) *Traffic {
	p, status := content.Resolve[Traffic](trafficName, 1, resolver)
	if !status.OK() {
		agent.Message(messaging.NewStatusMessage(status, ""))
		if curr == nil {
			curr = &Traffic{}
		}
		return curr
	}
	return &p
}


*/
