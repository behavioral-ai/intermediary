package cache

import (
	"github.com/behavioral-ai/collective/content"
	"github.com/behavioral-ai/collective/eventing"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/traffic/metrics"
)

// emissary attention
func emissaryAttend(agent *agentT, resolver *content.Resolution) {
	agent.dispatch(agent.emissary, messaging.StartupEvent)
	paused := false

	for {
		select {
		case <-agent.ticker.C():
			agent.dispatch(agent.ticker, messaging.ObservationEvent)
			if !paused {
				p, status := content.Resolve[metrics.TrafficProfile](metrics.TrafficProfileName, 1, resolver)
				if !status.OK() {
					agent.handler.Message(eventing.NewNotifyMessage(status))
				} else {
					agent.setEnabled(p)
				}
			}
		default:
		}
		select {
		case msg := <-agent.emissary.C:
			agent.dispatch(agent.emissary, msg.Event())
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
