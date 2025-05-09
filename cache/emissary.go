package cache

import (
	"github.com/behavioral-ai/core/messaging"
)

// emissary attention
func emissaryAttend(a *agentT) {
	paused := false

	for {
		select {
		case <-a.ticker.C():
			if !paused {
				a.enabled.Store(true) //agent.profile.Now())
			}
		default:
		}
		select {
		case msg := <-a.emissary.C:
			switch msg.Event() {
			case messaging.PauseEvent:
				paused = true
			case messaging.ResumeEvent:
				paused = false
			case messaging.ShutdownEvent:
				a.emissaryShutdown()
				return
			default:
			}
		default:
		}
	}
}
