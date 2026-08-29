package region

import "strings"

// Flag identifies a region flag.
type Flag string

// Group describes which players a flag applies to.
type Group string

const (
	GroupAll        Group = "all"
	GroupMembers    Group = "members"
	GroupOwners     Group = "owners"
	GroupNonMembers Group = "nonmembers"
	GroupNonOwners  Group = "nonowners"
)

// Overrides
const (
	FlagPassthrough              Flag = "passthrough"
	FlagNonPlayerProtectionDomains Flag = "nonplayer-protection-domains"
)

// Protection-related
const (
	FlagBuild             Flag = "build"
	FlagInteract          Flag = "interact"
	FlagBlockBreak        Flag = "block-break"
	FlagBlockPlace        Flag = "block-place"
	FlagUse               Flag = "use"
	FlagDamageAnimals     Flag = "damage-animals"
	FlagChestAccess       Flag = "chest-access"
	FlagRide              Flag = "ride"
	FlagPvP               Flag = "pvp"
	FlagSleep             Flag = "sleep"
	FlagRespawnAnchors    Flag = "respawn-anchors"
	FlagTNT               Flag = "tnt"
	FlagVehiclePlace      Flag = "vehicle-place"
	FlagVehicleDestroy    Flag = "vehicle-destroy"
	FlagLighter           Flag = "lighter"
	FlagBlockTrampling    Flag = "block-trampling"
	FlagFrostedIceForm    Flag = "frosted-ice-form"
	FlagItemFrameRotation Flag = "item-frame-rotation"
	FlagFireworkDamage    Flag = "firework-damage"
	FlagUseAnvil          Flag = "use-anvil"
	FlagUseDripleaf       Flag = "use-dripleaf"
)

// Mobs, fire, explosions
const (
	FlagCreeperExplosion        Flag = "creeper-explosion"
	FlagEnderdragonBlockDamage  Flag = "enderdragon-block-damage"
	FlagGhastFireball           Flag = "ghast-fireball"
	FlagOtherExplosion          Flag = "other-explosion"
	FlagFireSpread              Flag = "fire-spread"
	FlagEndermanGrief           Flag = "enderman-grief"
	FlagSnowmanTrails           Flag = "snowman-trails"
	FlagRavagerGrief            Flag = "ravager-grief"
	FlagMobDamage               Flag = "mob-damage"
	FlagMobSpawning             Flag = "mob-spawning"
	FlagDenySpawn               Flag = "deny-spawn"
	FlagEntityPaintingDestroy   Flag = "entity-painting-destroy"
	FlagEntityItemFrameDestroy  Flag = "entity-item-frame-destroy"
	FlagWitherDamage            Flag = "wither-damage"
)

// Natural events
const (
	FlagLavaFire       Flag = "lava-fire"
	FlagLightning      Flag = "lightning"
	FlagWaterFlow      Flag = "water-flow"
	FlagLavaFlow       Flag = "lava-flow"
	FlagSnowFall       Flag = "snow-fall"
	FlagSnowMelt       Flag = "snow-melt"
	FlagIceForm        Flag = "ice-form"
	FlagIceMelt        Flag = "ice-melt"
	FlagFrostedIceMelt Flag = "frosted-ice-melt"
	FlagMushroomGrowth Flag = "mushroom-growth"
	FlagLeafDecay      Flag = "leaf-decay"
	FlagGrassGrowth    Flag = "grass-growth"
	FlagMyceliumSpread Flag = "mycelium-spread"
	FlagVineGrowth     Flag = "vine-growth"
	FlagRockGrowth     Flag = "rock-growth"
	FlagSculkGrowth    Flag = "sculk-growth"
	FlagCropGrowth     Flag = "crop-growth"
	FlagSoilDry        Flag = "soil-dry"
	FlagCoralFade      Flag = "coral-fade"
	FlagCopperFade     Flag = "copper-fade"
)

// Movement
const (
	FlagEntry               Flag = "entry"
	FlagExit                Flag = "exit"
	FlagExitViaTeleport     Flag = "exit-via-teleport"
	FlagExitOverride        Flag = "exit-override"
	FlagEntryDenyMessage    Flag = "entry-deny-message"
	FlagExitDenyMessage     Flag = "exit-deny-message"
	FlagNotifyEnter         Flag = "notify-enter"
	FlagNotifyLeave         Flag = "notify-leave"
	FlagGreeting            Flag = "greeting"
	FlagGreetingTitle       Flag = "greeting-title"
	FlagFarewell            Flag = "farewell"
	FlagFarewellTitle       Flag = "farewell-title"
	FlagEnderpearl          Flag = "enderpearl"
	FlagChorusFruitTeleport Flag = "chorus-fruit-teleport"
	FlagTeleport            Flag = "teleport"
	FlagSpawn               Flag = "spawn"
	FlagTeleportMessage     Flag = "teleport-message"
)

// Map making
const (
	FlagItemPickup         Flag = "item-pickup"
	FlagItemDrop           Flag = "item-drop"
	FlagExpDrops           Flag = "exp-drops"
	FlagDenyMessage        Flag = "deny-message"
	FlagInvincible         Flag = "invincible"
	FlagFallDamage         Flag = "fall-damage"
	FlagGameMode           Flag = "game-mode"
	FlagTimeLock           Flag = "time-lock"
	FlagWeatherLock        Flag = "weather-lock"
	FlagNaturalHealthRegen Flag = "natural-health-regen"
	FlagNaturalHungerDrain Flag = "natural-hunger-drain"
	FlagHealDelay          Flag = "heal-delay"
	FlagHealAmount         Flag = "heal-amount"
	FlagHealMinHealth      Flag = "heal-min-health"
	FlagHealMaxHealth      Flag = "heal-max-health"
	FlagFeedDelay          Flag = "feed-delay"
	FlagFeedAmount         Flag = "feed-amount"
	FlagFeedMinHunger      Flag = "feed-min-hunger"
	FlagFeedMaxHunger      Flag = "feed-max-hunger"
	FlagBlockedCmds        Flag = "blocked-cmds"
	FlagAllowedCmds        Flag = "allowed-cmds"
)

// Miscellaneous
const (
	FlagPistons      Flag = "pistons"
	FlagSendChat     Flag = "send-chat"
	FlagReceiveChat  Flag = "receive-chat"
	FlagPotionSplash Flag = "potion-splash"
)

// knownFlags maps every WorldGuard flag to its default value ("" = unset/no default).
var knownFlags = map[Flag]string{
	// Overrides
	FlagPassthrough:              "",
	FlagNonPlayerProtectionDomains: "",
	// Protection
	FlagBuild:             "",
	FlagInteract:          "",
	FlagBlockBreak:        "",
	FlagBlockPlace:        "",
	FlagUse:               "",
	FlagDamageAnimals:     "",
	FlagChestAccess:       "",
	FlagRide:              "",
	FlagPvP:               "",
	FlagSleep:             "",
	FlagRespawnAnchors:    "",
	FlagTNT:               "",
	FlagVehiclePlace:      "",
	FlagVehicleDestroy:    "",
	FlagLighter:           "",
	FlagBlockTrampling:    "",
	FlagFrostedIceForm:    "",
	FlagItemFrameRotation: "",
	FlagFireworkDamage:    "",
	FlagUseAnvil:          "",
	FlagUseDripleaf:       "",
	// Mobs
	FlagCreeperExplosion:       "",
	FlagEnderdragonBlockDamage: "",
	FlagGhastFireball:          "",
	FlagOtherExplosion:         "",
	FlagFireSpread:             "",
	FlagEndermanGrief:          "",
	FlagSnowmanTrails:          "",
	FlagRavagerGrief:           "",
	FlagMobDamage:              "",
	FlagMobSpawning:            "",
	FlagDenySpawn:              "",
	FlagEntityPaintingDestroy:  "",
	FlagEntityItemFrameDestroy: "",
	FlagWitherDamage:           "",
	// Natural events
	FlagLavaFire:       "",
	FlagLightning:      "",
	FlagWaterFlow:      "",
	FlagLavaFlow:       "",
	FlagSnowFall:       "",
	FlagSnowMelt:       "",
	FlagIceForm:        "",
	FlagIceMelt:        "",
	FlagFrostedIceMelt: "",
	FlagMushroomGrowth: "",
	FlagLeafDecay:      "",
	FlagGrassGrowth:    "",
	FlagMyceliumSpread: "",
	FlagVineGrowth:     "",
	FlagRockGrowth:     "",
	FlagSculkGrowth:    "",
	FlagCropGrowth:     "",
	FlagSoilDry:        "",
	FlagCoralFade:      "",
	FlagCopperFade:     "",
	// Movement
	FlagEntry:               "",
	FlagExit:                "",
	FlagExitViaTeleport:     "",
	FlagExitOverride:        "",
	FlagEntryDenyMessage:    "",
	FlagExitDenyMessage:     "",
	FlagNotifyEnter:         "",
	FlagNotifyLeave:         "",
	FlagGreeting:            "",
	FlagGreetingTitle:       "",
	FlagFarewell:            "",
	FlagFarewellTitle:       "",
	FlagEnderpearl:          "",
	FlagChorusFruitTeleport: "",
	FlagTeleport:            "",
	FlagSpawn:               "",
	FlagTeleportMessage:     "",
	// Map making
	FlagItemPickup:         "",
	FlagItemDrop:           "",
	FlagExpDrops:           "",
	FlagDenyMessage:        "",
	FlagInvincible:         "",
	FlagFallDamage:         "",
	FlagGameMode:           "",
	FlagTimeLock:           "",
	FlagWeatherLock:        "",
	FlagNaturalHealthRegen: "",
	FlagNaturalHungerDrain: "",
	FlagHealDelay:          "",
	FlagHealAmount:         "",
	FlagHealMinHealth:      "",
	FlagHealMaxHealth:      "",
	FlagFeedDelay:          "",
	FlagFeedAmount:         "",
	FlagFeedMinHunger:      "",
	FlagFeedMaxHunger:      "",
	FlagBlockedCmds:        "",
	FlagAllowedCmds:        "",
	// Misc
	FlagPistons:      "",
	FlagSendChat:     "",
	FlagReceiveChat:  "",
	FlagPotionSplash: "",
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

// Unset removes a flag from the region so it falls back to its default.
func (r *Region) Unset(f Flag) {
	delete(r.Flags, f)
}

// BuildDenied reports whether block placement and breaking is denied
// for the given player inside the region.
func (r *Region) BuildDenied(username string) bool {
	if r.HasMember(username) {
		return false
	}
	v := strings.ToLower(r.Get(FlagBuild))
	return v == "deny" || v == ""
}

func normaliseName(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
