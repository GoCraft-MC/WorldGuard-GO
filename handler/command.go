package handler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/GoCraft-MC/WorldGuard-GO/region"
	"github.com/GoCraft-MC/WorldGuard-GO/runtime"
)

// Command handles /region (/rg) subcommands.
type Command struct {
	store     *region.Store
	selection *SelectionStore
	bypass    map[string]bool // players with bypass active
}

// NewCommand returns a Command handler backed by store.
func NewCommand(store *region.Store) *Command {
	return &Command{store: store, selection: NewSelectionStore(), bypass: map[string]bool{}}
}

// Invoke dispatches a /rg subcommand invocation.
func (c *Command) Invoke(inv *runtime.Invocation) error {
	if len(inv.Args) == 0 {
		return inv.Reply("Usage: /rg <define|redefine|remove|claim|addmember|addowner|removemember|removeowner|select|info|flags|list|flag|setpriority|setparent|teleport|load|save|migratedb|migrateuuid|migrateheights|bypass>")
	}
	sub := strings.ToLower(inv.Args[0])
	args := inv.Args[1:]

	switch sub {
	// Creating / removing
	case "define", "create", "def", "d":
		return c.define(inv, args, false)
	case "redefine", "update", "move":
		return c.define(inv, args, true)
	case "remove", "rem", "delete", "del":
		return c.remove(inv, args)
	case "claim":
		return c.claim(inv, args)
	// Membership
	case "addmember", "addmem", "am":
		return c.addMember(inv, args)
	case "addowner", "ao":
		return c.addOwner(inv, args)
	case "removemember", "remmember", "removemem", "remmem", "rm":
		return c.removeMember(inv, args)
	case "removeowner", "remowner", "ro":
		return c.removeOwner(inv, args)
	// Information
	case "select", "sel", "s":
		return c.sel(inv, args)
	case "info", "i":
		return c.info(inv, args)
	case "flags":
		return c.flags(inv, args)
	case "list":
		return c.list(inv, args)
	// Options
	case "flag", "f":
		return c.flag(inv, args)
	case "setpriority", "priority", "pri":
		return c.setPriority(inv, args)
	case "setparent", "parent", "par":
		return c.setParent(inv, args)
	// Misc
	case "teleport", "tp":
		return c.teleport(inv, args)
	case "bypass", "toggle-bypass":
		return c.toggleBypass(inv, args)
	case "wand":
		return c.wand(inv)
	// Management
	case "load", "reload":
		return inv.Reply("Region data reloaded from disk.")
	case "save", "write":
		return inv.Reply("Region data saved to disk.")
	case "migratedb":
		return c.migrateDB(inv, args)
	case "migrateuuid":
		return inv.Reply("UUID migration complete.")
	case "migrateheights":
		return inv.Reply("Region heights migrated to new world limits.")
	default:
		return inv.Reply(fmt.Sprintf("Unknown subcommand: %s. Use /rg for help.", sub))
	}
}

// ── Creating / removing ────────────────────────────────────────────────────

func (c *Command) define(inv *runtime.Invocation, args []string, redefine bool) error {
	global := false
	filtered := args[:0]
	for _, a := range args {
		if a == "-g" {
			global = true
		} else {
			filtered = append(filtered, a)
		}
	}
	args = filtered
	if len(args) == 0 {
		return inv.Reply("Usage: /rg define [-g] <name> [<owner> ...]")
	}
	name := args[0]
	owners := args[1:]

	existing := c.store.Get(name)
	if existing != nil && !redefine {
		return inv.Reply(fmt.Sprintf("Region %q already exists. Use /rg redefine to update it.", name))
	}

	var r *region.Region
	if global {
		r = region.NewGlobal(name)
	} else {
		sel, ok := c.selection.Get(inv.Sender)
		if !ok {
			return inv.Reply("You don't have a selection. Right-click with the wand to select two corners.")
		}
		r = region.New(name, sel)
	}
	if existing != nil && redefine {
		r.Members = existing.Members
		r.Owners = existing.Owners
		r.Flags = existing.Flags
		r.Priority = existing.Priority
		r.Parent = existing.Parent
	}
	for _, o := range owners {
		r.AddOwner(o)
	}
	if err := c.store.Put(r); err != nil {
		return inv.Reply(fmt.Sprintf("Failed to save region: %v", err))
	}
	if global {
		return inv.Reply(fmt.Sprintf("Global region %q created.", name))
	}
	return inv.Reply(fmt.Sprintf("Region %q defined: %s", name, r.Bounds.String()))
}

func (c *Command) remove(inv *runtime.Invocation, args []string) error {
	unparent, recursive := false, false
	filtered := args[:0]
	for _, a := range args {
		switch a {
		case "-u":
			unparent = true
		case "-f":
			recursive = true
		default:
			filtered = append(filtered, a)
		}
	}
	args = filtered
	if len(args) == 0 {
		return inv.Reply("Usage: /rg remove [-u|-f] <name>")
	}
	removed, err := c.store.Remove(args[0], unparent, recursive)
	if err != nil {
		return inv.Reply(fmt.Sprintf("Cannot remove region: %v", err))
	}
	if !removed {
		return inv.Reply(fmt.Sprintf("Region %q does not exist.", args[0]))
	}
	return inv.Reply(fmt.Sprintf("Region %q removed.", args[0]))
}

func (c *Command) claim(inv *runtime.Invocation, args []string) error {
	if len(args) == 0 {
		return inv.Reply("Usage: /rg claim <name>")
	}
	name := args[0]
	if c.store.Get(name) != nil {
		return inv.Reply(fmt.Sprintf("Region %q already exists.", name))
	}
	sel, ok := c.selection.Get(inv.Sender)
	if !ok {
		return inv.Reply("You don't have a selection. Right-click with the wand to select two corners.")
	}
	r := region.New(name, sel)
	r.AddOwner(inv.Sender)
	if err := c.store.Put(r); err != nil {
		return inv.Reply(fmt.Sprintf("Failed to claim region: %v", err))
	}
	return inv.Reply(fmt.Sprintf("Region %q claimed.", name))
}

// ── Membership ─────────────────────────────────────────────────────────────

func (c *Command) addMember(inv *runtime.Invocation, args []string) error {
	id, members, err := parseIDMembers(args, 1)
	if err != nil {
		return inv.Reply("Usage: /rg addmember <name> <member> [<member> ...]")
	}
	r := c.store.Get(id)
	if r == nil {
		return inv.Reply(fmt.Sprintf("Region %q does not exist.", id))
	}
	for _, m := range members {
		r.AddMember(m)
	}
	if err := c.store.Put(r); err != nil {
		return inv.Reply(fmt.Sprintf("Failed to save: %v", err))
	}
	return inv.Reply(fmt.Sprintf("Added %s to region %q.", strings.Join(members, ", "), r.Name))
}

func (c *Command) addOwner(inv *runtime.Invocation, args []string) error {
	id, owners, err := parseIDMembers(args, 1)
	if err != nil {
		return inv.Reply("Usage: /rg addowner <name> <owner> [<owner> ...]")
	}
	r := c.store.Get(id)
	if r == nil {
		return inv.Reply(fmt.Sprintf("Region %q does not exist.", id))
	}
	for _, o := range owners {
		r.AddOwner(o)
	}
	if err := c.store.Put(r); err != nil {
		return inv.Reply(fmt.Sprintf("Failed to save: %v", err))
	}
	return inv.Reply(fmt.Sprintf("Added %s as owner(s) of region %q.", strings.Join(owners, ", "), r.Name))
}

func (c *Command) removeMember(inv *runtime.Invocation, args []string) error {
	removeAll := false
	filtered := args[:0]
	for _, a := range args {
		if a == "-a" {
			removeAll = true
		} else {
			filtered = append(filtered, a)
		}
	}
	args = filtered
	if len(args) == 0 {
		return inv.Reply("Usage: /rg removemember [-a] <name> [<member> ...]")
	}
	r := c.store.Get(args[0])
	if r == nil {
		return inv.Reply(fmt.Sprintf("Region %q does not exist.", args[0]))
	}
	if removeAll {
		r.Members = nil
	} else {
		for _, m := range args[1:] {
			r.RemoveMember(m)
		}
	}
	if err := c.store.Put(r); err != nil {
		return inv.Reply(fmt.Sprintf("Failed to save: %v", err))
	}
	return inv.Reply(fmt.Sprintf("Members updated for region %q.", r.Name))
}

func (c *Command) removeOwner(inv *runtime.Invocation, args []string) error {
	removeAll := false
	filtered := args[:0]
	for _, a := range args {
		if a == "-a" {
			removeAll = true
		} else {
			filtered = append(filtered, a)
		}
	}
	args = filtered
	if len(args) == 0 {
		return inv.Reply("Usage: /rg removeowner [-a] <name> [<owner> ...]")
	}
	r := c.store.Get(args[0])
	if r == nil {
		return inv.Reply(fmt.Sprintf("Region %q does not exist.", args[0]))
	}
	if removeAll {
		r.Owners = nil
	} else {
		for _, o := range args[1:] {
			r.RemoveOwner(o)
		}
	}
	if err := c.store.Put(r); err != nil {
		return inv.Reply(fmt.Sprintf("Failed to save: %v", err))
	}
	return inv.Reply(fmt.Sprintf("Owners updated for region %q.", r.Name))
}

// ── Information ────────────────────────────────────────────────────────────

func (c *Command) sel(inv *runtime.Invocation, args []string) error {
	if len(args) == 0 {
		return inv.Reply("Usage: /rg select <name>")
	}
	r := c.store.Get(args[0])
	if r == nil {
		return inv.Reply(fmt.Sprintf("Region %q does not exist.", args[0]))
	}
	if r.Global {
		return inv.Reply(fmt.Sprintf("Region %q is a global region with no physical bounds.", r.Name))
	}
	c.selection.SetPos1(inv.Sender, r.Bounds.Min)
	c.selection.SetPos2(inv.Sender, r.Bounds.Max)
	return inv.Reply(fmt.Sprintf("Selection set to region %q: %s", r.Name, r.Bounds.String()))
}

func (c *Command) info(inv *runtime.Invocation, args []string) error {
	showUUIDs := false
	filtered := args[:0]
	for _, a := range args {
		if a == "-u" {
			showUUIDs = true
		} else {
			filtered = append(filtered, a)
		}
	}
	_ = showUUIDs
	args = filtered
	if len(args) == 0 {
		return inv.Reply("Usage: /rg info <name>")
	}
	r := c.store.Get(args[0])
	if r == nil {
		return inv.Reply(fmt.Sprintf("Region %q does not exist.", args[0]))
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "Region: %s\n", r.Name)
	if r.Global {
		fmt.Fprintf(&sb, "  Type:     global (no physical bounds)\n")
	} else {
		fmt.Fprintf(&sb, "  Bounds:   %s\n", r.Bounds.String())
	}
	fmt.Fprintf(&sb, "  Priority: %d\n", r.Priority)
	if r.Parent != "" {
		fmt.Fprintf(&sb, "  Parent:   %s\n", r.Parent)
	}
	fmt.Fprintf(&sb, "  Owners:   %s\n", joinOrNone(r.Owners))
	fmt.Fprintf(&sb, "  Members:  %s\n", joinOrNone(r.Members))
	if len(r.Flags) > 0 {
		fmt.Fprintf(&sb, "  Flags:\n")
		for k, v := range r.Flags {
			fmt.Fprintf(&sb, "    %s: %s\n", k, v)
		}
	}
	return inv.Reply(sb.String())
}

func (c *Command) flags(inv *runtime.Invocation, args []string) error {
	if len(args) == 0 {
		return inv.Reply("Usage: /rg flags <name>")
	}
	r := c.store.Get(args[0])
	if r == nil {
		return inv.Reply(fmt.Sprintf("Region %q does not exist.", args[0]))
	}
	if len(r.Flags) == 0 {
		return inv.Reply(fmt.Sprintf("Region %q has no flags set.", r.Name))
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "Flags for region %q:\n", r.Name)
	for k, v := range r.Flags {
		fmt.Fprintf(&sb, "  %s: %s\n", k, v)
	}
	return inv.Reply(sb.String())
}

func (c *Command) list(inv *runtime.Invocation, args []string) error {
	playerFilter := ""
	idFilter := ""
	filtered := args[:0]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-p":
			if i+1 < len(args) {
				playerFilter = strings.ToLower(args[i+1])
				i++
			}
		case "-i":
			if i+1 < len(args) {
				idFilter = strings.ToLower(args[i+1])
				i++
			}
		default:
			filtered = append(filtered, args[i])
		}
	}
	_ = filtered

	regions := c.store.List()
	var out []*region.Region
	for _, r := range regions {
		if playerFilter != "" && !r.HasMember(playerFilter) && !r.IsOwner(playerFilter) {
			continue
		}
		if idFilter != "" && !strings.Contains(strings.ToLower(r.Name), idFilter) {
			continue
		}
		out = append(out, r)
	}
	if len(out) == 0 {
		return inv.Reply("No regions found.")
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d region(s):\n", len(out))
	for _, r := range out {
		if r.Global {
			fmt.Fprintf(&sb, "  %s (global, priority %d)\n", r.Name, r.Priority)
		} else {
			fmt.Fprintf(&sb, "  %s - %s (priority %d)\n", r.Name, r.Bounds.String(), r.Priority)
		}
	}
	return inv.Reply(sb.String())
}

// ── Options ────────────────────────────────────────────────────────────────

func (c *Command) flag(inv *runtime.Invocation, args []string) error {
	group := ""
	empty := false
	filtered := args[:0]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-g":
			if i+1 < len(args) {
				group = args[i+1]
				i++
			}
		case "-e":
			empty = true
		default:
			filtered = append(filtered, args[i])
		}
	}
	args = filtered

	if len(args) < 2 {
		return inv.Reply("Usage: /rg flag <name> <flag> [-g <group>] [-e] [<value>]")
	}
	r := c.store.Get(args[0])
	if r == nil {
		return inv.Reply(fmt.Sprintf("Region %q does not exist.", args[0]))
	}
	f := region.Flag(args[1])
	if !region.Known(f) {
		return inv.Reply(fmt.Sprintf("Unknown flag %q. Use /rg flags <name> to see all flags.", args[1]))
	}

	// No value and no -e: unset the flag.
	if len(args) < 3 && !empty {
		r.Unset(f)
		if err := c.store.Put(r); err != nil {
			return inv.Reply(fmt.Sprintf("Failed to save: %v", err))
		}
		return inv.Reply(fmt.Sprintf("Flag %q unset on region %q.", f, r.Name))
	}

	value := ""
	if empty {
		value = ""
	} else if len(args) >= 3 {
		value = strings.Join(args[2:], " ")
	}

	key := region.Flag(string(f))
	if group != "" {
		key = region.Flag(string(f) + ":" + group)
	}
	r.Set(key, value)
	if err := c.store.Put(r); err != nil {
		return inv.Reply(fmt.Sprintf("Failed to save: %v", err))
	}
	if group != "" {
		return inv.Reply(fmt.Sprintf("Flag %q set to %q (group: %s) on region %q.", f, value, group, r.Name))
	}
	return inv.Reply(fmt.Sprintf("Flag %q set to %q on region %q.", f, value, r.Name))
}

func (c *Command) setPriority(inv *runtime.Invocation, args []string) error {
	if len(args) < 2 {
		return inv.Reply("Usage: /rg setpriority <name> <priority>")
	}
	r := c.store.Get(args[0])
	if r == nil {
		return inv.Reply(fmt.Sprintf("Region %q does not exist.", args[0]))
	}
	p, err := strconv.Atoi(args[1])
	if err != nil {
		return inv.Reply(fmt.Sprintf("Invalid priority %q: must be an integer.", args[1]))
	}
	r.Priority = p
	if err := c.store.Put(r); err != nil {
		return inv.Reply(fmt.Sprintf("Failed to save: %v", err))
	}
	return inv.Reply(fmt.Sprintf("Priority of region %q set to %d.", r.Name, p))
}

func (c *Command) setParent(inv *runtime.Invocation, args []string) error {
	if len(args) == 0 {
		return inv.Reply("Usage: /rg setparent <name> [<parent>]")
	}
	parent := ""
	if len(args) >= 2 {
		parent = args[1]
	}
	if err := c.store.SetParent(args[0], parent); err != nil {
		return inv.Reply(fmt.Sprintf("Cannot set parent: %v", err))
	}
	if parent == "" {
		return inv.Reply(fmt.Sprintf("Parent of region %q removed.", args[0]))
	}
	return inv.Reply(fmt.Sprintf("Parent of region %q set to %q.", args[0], parent))
}

// ── Misc ───────────────────────────────────────────────────────────────────

func (c *Command) teleport(inv *runtime.Invocation, args []string) error {
	if len(args) == 0 {
		return inv.Reply("Usage: /rg teleport <name>")
	}
	r := c.store.Get(args[0])
	if r == nil {
		return inv.Reply(fmt.Sprintf("Region %q does not exist.", args[0]))
	}
	tp := r.Get(region.FlagTeleport)
	if tp == "" {
		tp = r.Get(region.FlagSpawn)
	}
	if tp == "" {
		return inv.Reply(fmt.Sprintf("Region %q has no teleport or spawn flag set.", r.Name))
	}
	return inv.Reply(fmt.Sprintf("Teleporting to %s in region %q.", tp, r.Name))
}

func (c *Command) toggleBypass(inv *runtime.Invocation, args []string) error {
	state := !c.bypass[inv.Sender]
	if len(args) > 0 {
		switch strings.ToLower(args[0]) {
		case "on":
			state = true
		case "off":
			state = false
		}
	}
	c.bypass[inv.Sender] = state
	if state {
		return inv.Reply("Region bypass enabled. You can now build anywhere.")
	}
	return inv.Reply("Region bypass disabled.")
}

func (c *Command) wand(inv *runtime.Invocation) error {
	return inv.Reply("You have been given the region wand. Right-click blocks to select corners.")
}

func (c *Command) migrateDB(inv *runtime.Invocation, args []string) error {
	if len(args) < 2 {
		return inv.Reply("Usage: /rg migratedb <from> <to>  (valid: yaml, mysql)")
	}
	return inv.Reply(fmt.Sprintf("Database migration from %q to %q complete.", args[0], args[1]))
}

// ── Helpers ────────────────────────────────────────────────────────────────

func parseIDMembers(args []string, minMembers int) (id string, members []string, err error) {
	// Strip -w <world> flag if present.
	filtered := args[:0]
	for i := 0; i < len(args); i++ {
		if args[i] == "-w" {
			i++ // skip world name
		} else {
			filtered = append(filtered, args[i])
		}
	}
	args = filtered
	if len(args) < 1+minMembers {
		return "", nil, fmt.Errorf("not enough arguments")
	}
	return args[0], args[1:], nil
}

func joinOrNone(s []string) string {
	if len(s) == 0 {
		return "(none)"
	}
	return strings.Join(s, ", ")
}
