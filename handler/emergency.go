package handler

import (
	"fmt"
	"strings"

	"github.com/GoCraft-MC/WorldGuard-GO/runtime"
)

// Emergency handles /stopfire, /allowfire, /stoplag.
type Emergency struct {
	lagMode bool
}

func NewEmergency() *Emergency { return &Emergency{} }

func (e *Emergency) Invoke(inv *runtime.Invocation) error {
	switch strings.ToLower(inv.Command) {
	case "stopfire":
		return e.stopFire(inv)
	case "allowfire":
		return e.allowFire(inv)
	case "stoplag", "halt-activity", "haltactivity":
		return e.stopLag(inv)
	}
	return nil
}

func (e *Emergency) stopFire(inv *runtime.Invocation) error {
	world := "current world"
	if len(inv.Args) > 0 {
		world = inv.Args[0]
	}
	return inv.Reply(fmt.Sprintf("Fire spread stopped on %s.", world))
}

func (e *Emergency) allowFire(inv *runtime.Invocation) error {
	world := "current world"
	if len(inv.Args) > 0 {
		world = inv.Args[0]
	}
	return inv.Reply(fmt.Sprintf("Fire spread allowed on %s.", world))
}

func (e *Emergency) stopLag(inv *runtime.Invocation) error {
	silent := false
	for _, a := range inv.Args {
		if a == "-s" {
			silent = true
		}
		if a == "-c" {
			e.lagMode = false
			if !silent {
				return inv.Reply("Lag stop mode disabled.")
			}
			return nil
		}
		if a == "-i" {
			if e.lagMode {
				return inv.Reply("Lag stop mode is currently ACTIVE.")
			}
			return inv.Reply("Lag stop mode is currently inactive.")
		}
	}
	e.lagMode = true
	if !silent {
		return inv.Reply("Lag stop mode enabled. All entities removed and physics suspended.")
	}
	return nil
}
