package region

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

// Store holds all regions and persists them to a JSON file.
type Store struct {
	mu      sync.RWMutex
	path    string
	regions map[string]*Region // keyed by normalised region name
}

// NewStore loads an existing store from path, or creates an empty one.
func NewStore(path string) (*Store, error) {
	s := &Store{path: path, regions: map[string]*Region{}}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}
	var regions []*Region
	if err := json.Unmarshal(data, &regions); err != nil {
		return nil, err
	}
	for _, r := range regions {
		s.regions[normaliseName(r.Name)] = r
	}
	return s, nil
}

// Get returns the region with the given name, or nil if not found.
func (s *Store) Get(name string) *Region {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.regions[normaliseName(name)]
}

// List returns all regions sorted by name.
func (s *Store) List() []*Region {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Region, 0, len(s.regions))
	for _, r := range s.regions {
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// At returns all regions whose bounds contain pos, sorted by priority descending.
func (s *Store) At(pos Vec3) []*Region {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []*Region
	for _, r := range s.regions {
		if r.Global || r.Bounds.Contains(pos) {
			out = append(out, r)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Priority > out[j].Priority
	})
	return out
}

// Children returns all regions whose parent is name.
func (s *Store) Children(name string) []*Region {
	s.mu.RLock()
	defer s.mu.RUnlock()
	key := normaliseName(name)
	var out []*Region
	for _, r := range s.regions {
		if normaliseName(r.Parent) == key {
			out = append(out, r)
		}
	}
	return out
}

// SetParent sets the parent of name to parent, detecting circular inheritance.
// Pass an empty parent to remove the parent.
func (s *Store) SetParent(name, parent string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	r := s.regions[normaliseName(name)]
	if r == nil {
		return fmt.Errorf("region %q does not exist", name)
	}
	if parent == "" {
		r.Parent = ""
		return s.save()
	}
	if normaliseName(parent) == normaliseName(name) {
		return fmt.Errorf("a region cannot be its own parent")
	}
	// Detect circular inheritance.
	visited := map[string]bool{normaliseName(name): true}
	cur := normaliseName(parent)
	for cur != "" {
		if visited[cur] {
			return fmt.Errorf("circular inheritance detected")
		}
		visited[cur] = true
		p := s.regions[cur]
		if p == nil {
			break
		}
		cur = normaliseName(p.Parent)
	}
	r.Parent = parent
	return s.save()
}

// Put adds or replaces a region and persists the store.
func (s *Store) Put(r *Region) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.regions[normaliseName(r.Name)] = r
	return s.save()
}

// Remove deletes the region with the given name and persists the store.
// Returns false if no such region exists.
// If unparent is true, children are unparented; if recursive is true, children are deleted.
func (s *Store) Remove(name string, unparent, recursive bool) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := normaliseName(name)
	if _, ok := s.regions[key]; !ok {
		return false, nil
	}
	for _, r := range s.regions {
		if normaliseName(r.Parent) == key {
			if recursive {
				delete(s.regions, normaliseName(r.Name))
			} else if unparent {
				r.Parent = ""
			} else {
				return false, fmt.Errorf("region %q has child regions; use -u to unparent or -f to delete them", name)
			}
		}
	}
	delete(s.regions, key)
	return true, s.save()
}

func (s *Store) save() error {
	if s.path == "" {
		return nil
	}
	regions := make([]*Region, 0, len(s.regions))
	for _, r := range s.regions {
		regions = append(regions, r)
	}
	data, err := json.MarshalIndent(regions, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}
