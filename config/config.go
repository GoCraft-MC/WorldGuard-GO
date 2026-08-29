// Package config loads and exposes the WorldGuard-GO configuration.
package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config mirrors the WorldGuard config.yml structure exactly.
type Config struct {
	Regions struct {
		UUIDMigration struct {
			PerformOnNextStart    bool `yaml:"perform-on-next-start"`
			KeepNamesThatLackUUID bool `yaml:"keep-names-that-lack-uuids"`
		} `yaml:"uuid-migration"`
		UseCreatureSpawnEvent       bool   `yaml:"use-creature-spawn-event"`
		DisableBypassByDefault      bool   `yaml:"disable-bypass-by-default"`
		AnnounceBypassStatus        bool   `yaml:"announce-bypass-status"`
		UsePaperEntityOrigin        bool   `yaml:"use-paper-entity-origin"`
		Enable                      bool   `yaml:"enable"`
		InvincibilityRemovesMobs    bool   `yaml:"invincibility-removes-mobs"`
		CancelChatWithoutRecipients bool   `yaml:"cancel-chat-without-recipients"`
		NetherPortalProtection      bool   `yaml:"nether-portal-protection"`
		FakePlayerBuildOverride     bool   `yaml:"fake-player-build-override"`
		ExplosionFlagsBlockEntity   bool   `yaml:"explosion-flags-block-entity-damage"`
		HighFrequencyFlags          bool   `yaml:"high-frequency-flags"`
		ProtectAgainstLiquidFlow    bool   `yaml:"protect-against-liquid-flow"`
		Wand                        string `yaml:"wand"`
		MaxClaimVolume              int    `yaml:"max-claim-volume"`
		ClaimOnlyInsideExisting     bool   `yaml:"claim-only-inside-existing-regions"`
		SetParentOnClaim            string `yaml:"set-parent-on-claim"`
		LocationFlagsOnlyInside     bool   `yaml:"location-flags-only-inside-regions"`
		MaxRegionCountPerPlayer     struct {
			Default int `yaml:"default"`
		} `yaml:"max-region-count-per-player"`
	} `yaml:"regions"`

	AutoInvincible      bool `yaml:"auto-invincible"`
	AutoInvincibleGroup bool `yaml:"auto-invincible-group"`
	AutoNoDrowningGroup bool `yaml:"auto-no-drowning-group"`
	UsePlayerMoveEvent  bool `yaml:"use-player-move-event"`
	UsePlayerTeleports  bool `yaml:"use-player-teleports"`
	UseParticleEffects  bool `yaml:"use-particle-effects"`
	DisablePermCache    bool `yaml:"disable-permission-cache"`

	Security struct {
		DeopEveryoneOnJoin        bool `yaml:"deop-everyone-on-join"`
		BlockInGameOpCommand      bool `yaml:"block-in-game-op-command"`
		HostKeysAllowForgeClients bool `yaml:"host-keys-allow-forge-clients"`
	} `yaml:"security"`

	HostKeys       map[string]string `yaml:"host-keys"`
	SummaryOnStart bool              `yaml:"summary-on-start"`
	OpPermissions  bool              `yaml:"op-permissions"`

	BuildPermissionNodes struct {
		Enable      bool   `yaml:"enable"`
		DenyMessage string `yaml:"deny-message"`
	} `yaml:"build-permission-nodes"`

	EventHandling struct {
		BlockEntitySpawnsWithUntraceableCause bool     `yaml:"block-entity-spawns-with-untraceable-cause"`
		InteractionWhitelist                  []string `yaml:"interaction-whitelist"`
		EmitBlockUseAtFeet                    []string `yaml:"emit-block-use-at-feet"`
		IgnoreHopperItemMoveEvents            bool     `yaml:"ignore-hopper-item-move-events"`
		BreakHoppersOnDeniedMove              bool     `yaml:"break-hoppers-on-denied-move"`
	} `yaml:"event-handling"`

	Protection struct {
		ItemDurability       bool `yaml:"item-durability"`
		RemoveInfiniteStacks bool `yaml:"remove-infinite-stacks"`
		DisableXPOrbDrops    bool `yaml:"disable-xp-orb-drops"`
		UseMaxPriorityAssoc  bool `yaml:"use-max-priority-association"`
	} `yaml:"protection"`

	Gameplay struct {
		BlockPotions          []string `yaml:"block-potions"`
		BlockPotionsOverly    bool     `yaml:"block-potions-overly-reliably"`
		DisableConduitEffects bool     `yaml:"disable-conduit-effects"`
	} `yaml:"gameplay"`

	Default struct {
		PumpkinScuba        bool `yaml:"pumpkin-scuba"`
		DisableHealthRegain bool `yaml:"disable-health-regain"`
	} `yaml:"default"`

	Physics struct {
		NoPhysicsGravel          bool     `yaml:"no-physics-gravel"`
		NoPhysicsSand            bool     `yaml:"no-physics-sand"`
		VineLikeRopeLadders      bool     `yaml:"vine-like-rope-ladders"`
		AllowPortalAnywhere      bool     `yaml:"allow-portal-anywhere"`
		DisableWaterDamageBlocks []string `yaml:"disable-water-damage-blocks"`
	} `yaml:"physics"`

	Ignition struct {
		BlockTNT            bool `yaml:"block-tnt"`
		BlockTNTBlockDamage bool `yaml:"block-tnt-block-damage"`
		BlockLighter        bool `yaml:"block-lighter"`
	} `yaml:"ignition"`

	Fire struct {
		DisableLavaFireSpread   bool     `yaml:"disable-lava-fire-spread"`
		DisableAllFireSpread    bool     `yaml:"disable-all-fire-spread"`
		DisableFireSpreadBlocks []string `yaml:"disable-fire-spread-blocks"`
		LavaSpreadBlocks        []string `yaml:"lava-spread-blocks"`
	} `yaml:"fire"`

	Mobs struct {
		BlockCreeperExplosions         bool     `yaml:"block-creeper-explosions"`
		BlockCreeperBlockDamage        bool     `yaml:"block-creeper-block-damage"`
		BlockWitherExplosions          bool     `yaml:"block-wither-explosions"`
		BlockWitherBlockDamage         bool     `yaml:"block-wither-block-damage"`
		BlockWitherSkullExplosions     bool     `yaml:"block-wither-skull-explosions"`
		BlockWitherSkullBlockDamage    bool     `yaml:"block-wither-skull-block-damage"`
		BlockEnderdragonBlockDamage    bool     `yaml:"block-enderdragon-block-damage"`
		BlockEnderdragonPortalCreation bool     `yaml:"block-enderdragon-portal-creation"`
		BlockFireballExplosions        bool     `yaml:"block-fireball-explosions"`
		BlockFireballBlockDamage       bool     `yaml:"block-fireball-block-damage"`
		BlockWindchargeExplosions      bool     `yaml:"block-windcharge-explosions"`
		AntiWolfDumbness               bool     `yaml:"anti-wolf-dumbness"`
		AllowTamedSpawns               bool     `yaml:"allow-tamed-spawns"`
		DisableEndermanGriefing        bool     `yaml:"disable-enderman-griefing"`
		DisableSnowmanTrails           bool     `yaml:"disable-snowman-trails"`
		BlockPaintingDestroy           bool     `yaml:"block-painting-destroy"`
		BlockItemFrameDestroy          bool     `yaml:"block-item-frame-destroy"`
		BlockArmorStandDestroy         bool     `yaml:"block-armor-stand-destroy"`
		BlockPluginSpawning            bool     `yaml:"block-plugin-spawning"`
		BlockAboveGroundSlimes         bool     `yaml:"block-above-ground-slimes"`
		BlockOtherExplosions           bool     `yaml:"block-other-explosions"`
		BlockZombieDoorDestruction     bool     `yaml:"block-zombie-door-destruction"`
		BlockVehicleEntry              bool     `yaml:"block-vehicle-entry"`
		BlockCreatureSpawn             []string `yaml:"block-creature-spawn"`
	} `yaml:"mobs"`

	PlayerDamage struct {
		DisableFallDamage        bool `yaml:"disable-fall-damage"`
		DisableLavaDamage        bool `yaml:"disable-lava-damage"`
		DisableFireDamage        bool `yaml:"disable-fire-damage"`
		DisableLightningDamage   bool `yaml:"disable-lightning-damage"`
		DisableDrowningDamage    bool `yaml:"disable-drowning-damage"`
		DisableSuffocationDamage bool `yaml:"disable-suffocation-damage"`
		DisableContactDamage     bool `yaml:"disable-contact-damage"`
		TeleportOnSuffocation    bool `yaml:"teleport-on-suffocation"`
		DisableVoidDamage        bool `yaml:"disable-void-damage"`
		TeleportOnVoidFalling    bool `yaml:"teleport-on-void-falling"`
		ResetFallOnVoidTeleport  bool `yaml:"reset-fall-on-void-teleport"`
		DisableExplosionDamage   bool `yaml:"disable-explosion-damage"`
		DisableMobDamage         bool `yaml:"disable-mob-damage"`
		DisableDeathMessages     bool `yaml:"disable-death-messages"`
	} `yaml:"player-damage"`

	Crops struct {
		DisableCreatureTrampling bool `yaml:"disable-creature-trampling"`
		DisablePlayerTrampling   bool `yaml:"disable-player-trampling"`
	} `yaml:"crops"`

	TurtleEgg struct {
		DisableCreatureTrampling bool `yaml:"disable-creature-trampling"`
		DisablePlayerTrampling   bool `yaml:"disable-player-trampling"`
	} `yaml:"turtle-egg"`

	SnifferEgg struct {
		DisableCreatureTrampling bool `yaml:"disable-creature-trampling"`
		DisablePlayerTrampling   bool `yaml:"disable-player-trampling"`
	} `yaml:"sniffer-egg"`

	Weather struct {
		PreventLightningStrikeBlocks  []string `yaml:"prevent-lightning-strike-blocks"`
		DisableLightningStrikeFire    bool     `yaml:"disable-lightning-strike-fire"`
		DisableThunderstorm           bool     `yaml:"disable-thunderstorm"`
		DisableWeather                bool     `yaml:"disable-weather"`
		DisablePigZombification       bool     `yaml:"disable-pig-zombification"`
		DisableVillagerWitchification bool     `yaml:"disable-villager-witchification"`
		DisablePoweredCreepers        bool     `yaml:"disable-powered-creepers"`
		AlwaysRaining                 bool     `yaml:"always-raining"`
		AlwaysThundering              bool     `yaml:"always-thundering"`
	} `yaml:"weather"`

	Dynamics struct {
		DisableMushroomSpread  bool     `yaml:"disable-mushroom-spread"`
		DisableIceMelting      bool     `yaml:"disable-ice-melting"`
		DisableSnowMelting     bool     `yaml:"disable-snow-melting"`
		DisableSnowFormation   bool     `yaml:"disable-snow-formation"`
		DisableIceFormation    bool     `yaml:"disable-ice-formation"`
		DisableLeafDecay       bool     `yaml:"disable-leaf-decay"`
		DisableGrassGrowth     bool     `yaml:"disable-grass-growth"`
		DisableMyceliumSpread  bool     `yaml:"disable-mycelium-spread"`
		DisableVineGrowth      bool     `yaml:"disable-vine-growth"`
		DisableRockGrowth      bool     `yaml:"disable-rock-growth"`
		DisableSculkGrowth     bool     `yaml:"disable-sculk-growth"`
		DisableCropGrowth      bool     `yaml:"disable-crop-growth"`
		DisableSoilDehydration bool     `yaml:"disable-soil-dehydration"`
		DisableSoilMoisture    bool     `yaml:"disable-soil-moisture-change"`
		DisableCoralBlockFade  bool     `yaml:"disable-coral-block-fade"`
		DisableCopperBlockFade bool     `yaml:"disable-copper-block-fade"`
		SnowFallBlocks         []string `yaml:"snow-fall-blocks"`
	} `yaml:"dynamics"`

	Blacklist struct {
		UseAsWhitelist bool `yaml:"use-as-whitelist"`
		Logging        struct {
			Console struct {
				Enable bool `yaml:"enable"`
			} `yaml:"console"`
			Database struct {
				Enable bool   `yaml:"enable"`
				DSN    string `yaml:"dsn"`
				User   string `yaml:"user"`
				Pass   string `yaml:"pass"`
				Table  string `yaml:"table"`
			} `yaml:"database"`
			File struct {
				Enable    bool   `yaml:"enable"`
				Path      string `yaml:"path"`
				OpenFiles int    `yaml:"open-files"`
			} `yaml:"file"`
		} `yaml:"logging"`
	} `yaml:"blacklist"`

	CustomMetricsCharts bool `yaml:"custom-metrics-charts"`
}

// Load reads config.yml from path, or returns defaults if the file does not exist.
func Load(path string) (*Config, error) {
	c := defaults()
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return c, nil
	}
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(data, c); err != nil {
		return nil, err
	}
	return c, nil
}

// defaults returns a Config pre-populated with the same defaults as WorldGuard.
func defaults() *Config {
	c := &Config{}
	c.Regions.UUIDMigration.KeepNamesThatLackUUID = true
	c.Regions.UseCreatureSpawnEvent = true
	c.Regions.Enable = true
	c.Regions.CancelChatWithoutRecipients = true
	c.Regions.NetherPortalProtection = true
	c.Regions.FakePlayerBuildOverride = true
	c.Regions.ExplosionFlagsBlockEntity = true
	c.Regions.Wand = "minecraft:leather"
	c.Regions.MaxClaimVolume = 30000
	c.Regions.MaxRegionCountPerPlayer.Default = 7
	c.UsePlayerMoveEvent = true
	c.UsePlayerTeleports = true
	c.UseParticleEffects = true
	c.SummaryOnStart = true
	c.OpPermissions = true
	c.BuildPermissionNodes.DenyMessage = "&eSorry, but you are not permitted to do that here."
	c.EventHandling.BreakHoppersOnDeniedMove = true
	c.Protection.ItemDurability = true
	c.Mobs.AllowTamedSpawns = true
	c.Mobs.BlockPluginSpawning = true
	c.Blacklist.Logging.Console.Enable = true
	c.Blacklist.Logging.Database.DSN = "jdbc:mysql://localhost:3306/minecraft"
	c.Blacklist.Logging.Database.User = "root"
	c.Blacklist.Logging.Database.Table = "blacklist_events"
	c.Blacklist.Logging.File.Path = "worldguard/logs/%Y-%m-%d.log"
	c.Blacklist.Logging.File.OpenFiles = 10
	c.CustomMetricsCharts = true
	return c
}
