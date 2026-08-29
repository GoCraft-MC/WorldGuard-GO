# WorldGuard-GO

A region protection plugin for [GoCraft](https://github.com/GoCraft-MC/GoCraft) — the Minecraft Java + Bedrock server written in Go.

This is an **example plugin** that demonstrates how to build a GoCraft plugin in Go. It replicates the core feature set of WorldGuard: named regions, build protection, flags, and member management.

---

## Features

- Define named regions using a two-corner wand selection
- Block breaking and placing is denied inside protected regions for non-members
- Per-region flags: `build`, `pvp`, `mob-spawn`, `greeting`, `farewell`
- Member and owner management per region
- `worldguard.bypass` permission to override protection
- Regions persist to a JSON file across restarts

---

## Commands

| Command | Permission | Description |
|---|---|---|
| `/region define <name>` | `worldguard.region.define` | Define a new region from your selection |
| `/region redefine <name>` | `worldguard.region.redefine` | Update an existing region's bounds |
| `/region remove <name>` | `worldguard.region.remove` | Delete a region |
| `/region flag <name> <flag> <value>` | `worldguard.region.flag` | Set a flag on a region |
| `/region addmember <name> <player>` | `worldguard.region.addmember` | Add a member to a region |
| `/region removemember <name> <player>` | `worldguard.region.removemember` | Remove a member from a region |
| `/region info <name>` | `worldguard.region.info` | Show region details |
| `/region list` | `worldguard.region.list` | List all regions |

---

## Flags

| Flag | Default | Description |
|---|---|---|
| `build` | `deny` | Whether non-members can break or place blocks |
| `pvp` | `allow` | Whether players can damage each other |
| `mob-spawn` | `allow` | Whether hostile mobs can spawn |
| `greeting` | _(empty)_ | Message sent to players entering the region |
| `farewell` | _(empty)_ | Message sent to players leaving the region |

**Example:**
```
/region flag spawn build deny
/region flag spawn greeting Welcome to spawn!
```

---

## Permissions

| Permission | Description |
|---|---|
| `worldguard.region.define` | Define new regions |
| `worldguard.region.redefine` | Redefine existing regions |
| `worldguard.region.remove` | Remove regions |
| `worldguard.region.flag` | Set region flags |
| `worldguard.region.addmember` | Add members |
| `worldguard.region.removemember` | Remove members |
| `worldguard.region.info` | View region info |
| `worldguard.region.list` | List regions |
| `worldguard.bypass` | Bypass all region protection |

---

## Building

Requires Go 1.23+.

```bash
go build -o worldguard .
```

This produces a `worldguard` binary that the GoCraft server loads as a plugin.

---

## Installation

1. Build the binary (see above)
2. Copy `worldguard` and `plugin.toml` into your GoCraft server's `plugins/worldguard/` directory
3. Restart the server — GoCraft will load the plugin automatically

---

## Project structure

```
WorldGuard-GO/
├── plugin.toml          # GoCraft plugin manifest
├── main.go              # Entry point
├── region/
│   ├── region.go        # Region, AABB and Vec3 types
│   ├── flags.go         # Flag definitions and defaults
│   └── store.go         # JSON-backed region persistence
├── handler/
│   ├── protect.go       # block.break / block.place protection
│   ├── selection.go     # Two-corner wand selection per player
│   └── command.go       # /region subcommand handler
└── runtime/
    └── runtime.go       # GoCraft IPC client (placeholder for the Go SDK)
```

> **Note:** `runtime/runtime.go` is a thin IPC client that mirrors the GoCraft plugin wire protocol. It will be replaced by the official GoCraft Go SDK once the Go plugin runtime is released.

---

## Status

| Area | Status |
|---|---|
| Region definition and storage | Done |
| Block break / place protection | Done |
| Member and owner management | Done |
| Build, pvp, mob-spawn flags | Done |
| Greeting / farewell messages | Declared, not yet dispatched (requires `player.join` region tracking) |
| Wand item selection (left/right click) | Stub — requires item interaction events from GoCraft |
| PvP enforcement | Declared, not yet enforced (requires entity damage events) |
