package region

import (
	"encoding/json"
	"os"
	"path/filepath"
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

// List returns all regions.
func (s *Store) List() []*Region {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Region, 0, len(s.regions))
	for _, r := range s.regions {
		out = append(out, r)
	}
	return out
}

// At returns all regions whose bounds contain pos.
func (s *Store) At(pos Vec3) []*Region {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []*Region
	for _, r := range s.regions {
		if r.Bounds.Contains(pos) {
			out = append(out, r)
		}
	}
	return out
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
func (s *Store) Remove(name string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := normaliseName(name)
	if _, ok := s.regions[key]; !ok {
		return false, nil
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
