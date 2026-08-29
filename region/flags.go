package region

import "strings"

// Flag identifies a region flag.
type Flag string

const (
	// FlagBuild controls whether non-members can place or break blocks.
	// Values: "allow" | "deny" (default: "deny")
	FlagBuild Flag = "build"

	// FlagPvP controls player-vs-player damage inside the region.
	// Values: "allow" | "deny" (default: "allow")
	FlagPvP Flag = "pvp"

	// FlagMobSpawn controls hostile mob spawning.
	// Values: "allow" | "deny" (default: "allow")
	FlagMobSpawn Flag = "mob-spawn"

	// FlagGreeting is a message sent to players on region entry.
	FlagGreeting Flag = "greeting"

	// FlagFarewell is a message sent to players on region exit.
	FlagFarewell Flag = "farewell"
)

var knownFlags = map[Flag]string{
	FlagBuild:    "deny",
	FlagPvP:      "allow",
	FlagMobSpawn: "allow",
	FlagGreeting: "",
	FlagFarewell: "",
}

// Known reports whether f is a recognised flag name.
func Known(f Flag) bool {
	_, ok := knownFlags[f]
	return ok
}

// Default returns the default value for flag f.
func Default(f Flag) string {
	return knownFlags[f]
}

// Get returns the effective value of flag f on region r,
// falling back to the flag default.
func (r *Region) Get(f Flag) string {
	if v, ok := r.Flags[f]; ok {
		return v
	}
	return Default(f)
}

// Set sets flag f to value v on region r.
func (r *Region) Set(f Flag, v string) {
	r.Flags[f] = v
}

// BuildDenied reports whether block placement and breaking is denied
// for the given player inside the region.
func (r *Region) BuildDenied(username string) bool {
	if r.HasMember(username) {
		return false
	}
	return strings.ToLower(r.Get(FlagBuild)) == "deny"
}

func normaliseName(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
