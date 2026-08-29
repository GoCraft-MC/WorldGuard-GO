package region

import "fmt"

// Vec3 is an integer block position.
type Vec3 struct {
	X, Y, Z int
}

// AABB is an axis-aligned bounding box in block coordinates.
type AABB struct {
	Min, Max Vec3
}

// Contains reports whether pos is inside the bounding box (inclusive).
func (a AABB) Contains(pos Vec3) bool {
	return pos.X >= a.Min.X && pos.X <= a.Max.X &&
		pos.Y >= a.Min.Y && pos.Y <= a.Max.Y &&
		pos.Z >= a.Min.Z && pos.Z <= a.Max.Z
}

// String returns a human-readable representation.
func (a AABB) String() string {
	return fmt.Sprintf("(%d,%d,%d) -> (%d,%d,%d)",
		a.Min.X, a.Min.Y, a.Min.Z,
		a.Max.X, a.Max.Y, a.Max.Z)
}

// Region is a named, protected area of the world.
type Region struct {
	Name    string          `json:"name"`
	Bounds  AABB            `json:"bounds"`
	Flags   map[Flag]string `json:"flags"`
	Members []string        `json:"members"` // usernames (lowercase)
	Owners  []string        `json:"owners"`  // usernames (lowercase)
}

// New returns a Region with default flags.
func New(name string, bounds AABB) *Region {
	return &Region{
		Name:   name,
		Bounds: bounds,
		Flags:  map[Flag]string{},
	}
}

// HasMember reports whether the username is a member or owner of the region.
func (r *Region) HasMember(username string) bool {
	username = normaliseName(username)
	for _, m := range r.Members {
		if m == username {
			return true
		}
	}
	for _, o := range r.Owners {
		if o == username {
			return true
		}
	}
	return false
}

// AddMember adds username as a member if not already present.
func (r *Region) AddMember(username string) {
	username = normaliseName(username)
	if !r.HasMember(username) {
		r.Members = append(r.Members, username)
	}
}

// RemoveMember removes username from both member and owner lists.
func (r *Region) RemoveMember(username string) {
	username = normaliseName(username)
	r.Members = removeString(r.Members, username)
	r.Owners = removeString(r.Owners, username)
}

func removeString(s []string, v string) []string {
	out := s[:0]
	for _, x := range s {
		if x != v {
			out = append(out, x)
		}
	}
	return out
}
