package module

import (
	"github.com/behavioral-ai/intermediary/cache"
	"github.com/behavioral-ai/intermediary/routing"
)

var (
	CacheNamespaceName   = cache.NamespaceName
	RoutingNamespaceName = routing.NamespaceName
)

func Resolve(name string) (bool, any) {
	switch name {
	case cache.NamespaceName:
		return true, nil
	case routing.NamespaceName:
		return true, nil
	default:
		return false, nil
	}
}
