package handler

import (
	"fmt"
	"strings"

	"github.com/GoCraft-MC/WorldGuard-GO/region"
	"github.com/GoCraft-MC/WorldGuard-GO/runtime"
)

// Command handles /region subcommands.
type Command struct {
	store     *region.Store
	selection *SelectionStore
}

// NewCommand returns a Command handler backed by store.
func NewCommand(store *region.Store) *Command {
	return &Command{store: store, selection: NewSelectionStore()}
}

// Invoke dispatches a /region subcommand invocation.
func (c *Command) Invoke(inv *runtime.Invocation) error {
	if len(inv.Args) == 0 {
		return inv.Reply("Usage: /region <define|redefine|remove|flag|addmember|removemember|info|list>")
	}
	sub := strings.ToLower(inv.Args[0])
	args := inv.Args[1:]
	switch sub {
	case "define":
		return c.define(inv, args, false)
	case "redefine":
		return c.define(inv, args, true)
	case "remove":
		return c.remove(inv, args)
	case "flag":
		return c.flag(inv, args)
	case "addmember":
		return c.addMember(inv, args)
	case "removemember":
		return c.removeMember(inv, args)
	case "info":
		return c.info(inv, args)
	case "list":
		return c.list(inv)
	default:
		return inv.Reply(fmt.Sprintf("Unknown subcommand: %s", sub))
	}
}

func (c *Command) define(inv *runtime.Invocation, args []string, redefine bool) error {
	if len(args) == 0 {
		return inv.Reply("Usage: /region define <name>")
	}
	name := args[0]
	sel, ok := c.selection.Get(inv.Sender)
	if !ok {
		return inv.Reply("You don't have a selection. Use your wand to select two corners first.")
	}
	if c.store.Get(name) != nil && !redefine {
		return inv.Reply(fmt.Sprintf("Region %q already exists. Use /region redefine to update it.", name))
	}
	if err := c.store.Put(region.New(name, sel)); err != nil {
		return inv.Reply(fmt.Sprintf("Failed to save region: %v", err))
	}
	return inv.Reply(fmt.Sprintf("Region %q defined: %s", name, sel.String()))
}

func (c *Command) remove(inv *runtime.Invocation, args []string) error {
	if len(args) == 0 {
		return inv.Reply("Usage: /region remove <name>")
	}
	removed, err := c.store.Remove(args[0])
	if err != nil {
		return inv.Reply(fmt.Sprintf("Failed to remove region: %v", err))
	}
	if !removed {
		return inv.Reply(fmt.Sprintf("Region %q does not exist.", args[0]))
	}
	return inv.Reply(fmt.Sprintf("Region %q removed.", args[0]))
}

func (c *Command) flag(inv *runtime.Invocation, args []string) error {
	if len(args) < 3 {
		return inv.Reply("Usage: /region flag <name> <flag> <value>")
	}
	r := c.store.Get(args[0])
	if r == nil {
		return inv.Reply(fmt.Sprintf("Region %q does not exist.", args[0]))
	}
	f := region.Flag(args[1])
	if !region.Known(f) {
		return inv.Reply(fmt.Sprintf("Unknown flag %q. Known: build, pvp, mob-spawn, greeting, farewell", args[1]))
	}
	r.Set(f, args[2])
	if err := c.store.Put(r); err != nil {
		return inv.Reply(fmt.Sprintf("Failed to save region: %v", err))
	}
	return inv.Reply(fmt.Sprintf("Flag %q set to %q on region %q.", f, args[2], r.Name))
}

func (c *Command) addMember(inv *runtime.Invocation, args []string) error {
	if len(args) < 2 {
		return inv.Reply("Usage: /region addmember <name> <player>")
	}
	r := c.store.Get(args[0])
	if r == nil {
		return inv.Reply(fmt.Sprintf("Region %q does not exist.", args[0]))
	}
	r.AddMember(args[1])
	if err := c.store.Put(r); err != nil {
		return inv.Reply(fmt.Sprintf("Failed to save region: %v", err))
	}
	return inv.Reply(fmt.Sprintf("Added %q to region %q.", args[1], r.Name))
}

func (c *Command) removeMember(inv *runtime.Invocation, args []string) error {
	if len(args) < 2 {
		return inv.Reply("Usage: /region removemember <name> <player>")
	}
	r := c.store.Get(args[0])
	if r == nil {
		return inv.Reply(fmt.Sprintf("Region %q does not exist.", args[0]))
	}
	r.RemoveMember(args[1])
	if err := c.store.Put(r); err != nil {
		return inv.Reply(fmt.Sprintf("Failed to save region: %v", err))
	}
	return inv.Reply(fmt.Sprintf("Removed %q from region %q.", args[1], r.Name))
}

func (c *Command) info(inv *runtime.Invocation, args []string) error {
	if len(args) == 0 {
		return inv.Reply("Usage: /region info <name>")
	}
	r := c.store.Get(args[0])
	if r == nil {
		return inv.Reply(fmt.Sprintf("Region %q does not exist.", args[0]))
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "Region %q\n", r.Name)
	fmt.Fprintf(&sb, "  Bounds:  %s\n", r.Bounds.String())
	fmt.Fprintf(&sb, "  Members: %s\n", strings.Join(r.Members, ", "))
	fmt.Fprintf(&sb, "  Owners:  %s\n", strings.Join(r.Owners, ", "))
	for k, v := range r.Flags {
		fmt.Fprintf(&sb, "  %s = %s\n", k, v)
	}
	return inv.Reply(sb.String())
}

func (c *Command) list(inv *runtime.Invocation) error {
	regions := c.store.List()
	if len(regions) == 0 {
		return inv.Reply("No regions defined.")
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d region(s):\n", len(regions))
	for _, r := range regions {
		fmt.Fprintf(&sb, "  %s - %s\n", r.Name, r.Bounds.String())
	}
	return inv.Reply(sb.String())
}
