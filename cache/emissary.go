package cache

import (
	"github.com/behavioral-ai/collective/content"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/traffic/profile"
)

// emissary attention
func emissaryAttend(agent *agentT, resolver *content.Resolution) {
	paused := false
	var p *profile.Traffic

	for {
		select {
		case <-agent.ticker.C():
			if !paused {
				p = profile.NewTraffic(p, resolver)
				agent.setEnabled(p)
			}
		default:
		}
		select {
		case msg := <-agent.emissary.C:
			switch msg.Event() {
			case messaging.PauseEvent:
				paused = true
			case messaging.ResumeEvent:
				paused = false
			case messaging.ShutdownEvent:
				agent.emissaryShutdown()
				return
			default:
			}
		default:
		}
	}
}
