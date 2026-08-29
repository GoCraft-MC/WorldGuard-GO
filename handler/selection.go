package handler

import (
	"sync"

	"github.com/GoCraft-MC/WorldGuard-GO/region"
)

// SelectionStore tracks the wand selection (two corners) for each player.
type SelectionStore struct {
	mu   sync.Mutex
	pos1 map[string]region.Vec3
	pos2 map[string]region.Vec3
}

// NewSelectionStore returns an empty SelectionStore.
func NewSelectionStore() *SelectionStore {
	return &SelectionStore{
		pos1: map[string]region.Vec3{},
		pos2: map[string]region.Vec3{},
	}
}

// SetPos1 records the first corner for the player.
func (s *SelectionStore) SetPos1(player string, pos region.Vec3) {
	s.mu.Lock()
	s.pos1[player] = pos
	s.mu.Unlock()
}

// SetPos2 records the second corner for the player.
func (s *SelectionStore) SetPos2(player string, pos region.Vec3) {
	s.mu.Lock()
	s.pos2[player] = pos
	s.mu.Unlock()
}

// Get returns the normalised AABB if both corners are set, false otherwise.
func (s *SelectionStore) Get(player string) (region.AABB, bool) {
	s.mu.Lock()
	p1, ok1 := s.pos1[player]
	p2, ok2 := s.pos2[player]
	s.mu.Unlock()
	if !ok1 || !ok2 {
		return region.AABB{}, false
	}
	return normalise(p1, p2), true
}

func normalise(a, b region.Vec3) region.AABB {
	return region.AABB{
		Min: region.Vec3{X: minInt(a.X, b.X), Y: minInt(a.Y, b.Y), Z: minInt(a.Z, b.Z)},
		Max: region.Vec3{X: maxInt(a.X, b.X), Y: maxInt(a.Y, b.Y), Z: maxInt(a.Z, b.Z)},
	}
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
