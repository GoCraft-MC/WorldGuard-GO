// Package handler wires GoCraft plugin events to WorldGuard region logic.
package handler

import (
	"github.com/GoCraft-MC/WorldGuard-GO/region"
	"github.com/GoCraft-MC/WorldGuard-GO/runtime"
)

// Protect handles block.break and block.place events by checking whether
// the acting player is allowed to build inside any overlapping region.
type Protect struct {
	store *region.Store
}

// NewProtect returns a Protect handler backed by store.
func NewProtect(store *region.Store) *Protect {
	return &Protect{store: store}
}

// OnBlockBreak is called for every block.break event.
func (p *Protect) OnBlockBreak(e *runtime.Event) runtime.Verdict {
	return p.checkBuild(e)
}

// OnBlockPlace is called for every block.place event.
func (p *Protect) OnBlockPlace(e *runtime.Event) runtime.Verdict {
	return p.checkBuild(e)
}

func (p *Protect) checkBuild(e *runtime.Event) runtime.Verdict {
	pos := region.Vec3{
		X: e.IntField("block_x"),
		Y: e.IntField("block_y"),
		Z: e.IntField("block_z"),
	}
	username := e.StringField("player_name")

	if e.HasPermission("worldguard.bypass") {
		return runtime.Allow()
	}

	for _, r := range p.store.At(pos) {
		if r.BuildDenied(username) {
			return runtime.Deny()
		}
	}
	return runtime.Allow()
}
