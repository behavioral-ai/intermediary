package cache

import "github.com/behavioral-ai/core/messaging"

var (
	Agent messaging.Agent
)

func Initialize(ops messaging.Agent) {
	Agent = New(ops)
}
