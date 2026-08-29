package handler

import (
	"fmt"
	"strings"

	"github.com/GoCraft-MC/WorldGuard-GO/runtime"
)

// Gameplay handles /god, /ungod, /heal, /slay, /locate, /stack.
type Gameplay struct{}

func NewGameplay() *Gameplay { return &Gameplay{} }

func (g *Gameplay) Invoke(inv *runtime.Invocation) error {
	switch strings.ToLower(inv.Command) {
	case "god":
		return g.god(inv, false)
	case "ungod":
		return g.god(inv, true)
	case "heal":
		return g.heal(inv)
	case "slay":
		return g.slay(inv)
	case "locate":
		return g.locate(inv)
	case "stack":
		return g.stack(inv)
	}
	return nil
}

func (g *Gameplay) god(inv *runtime.Invocation, remove bool) error {
	silent, target := parseSilentTarget(inv.Args, inv.Sender)
	action := "granted invincibility"
	if remove {
		action = "removed invincibility"
	}
	if !silent {
		if target == inv.Sender {
			return inv.Reply(fmt.Sprintf("You have been %s.", action))
		}
		return inv.Reply(fmt.Sprintf("%s has been %s.", target, action))
	}
	return nil
}

func (g *Gameplay) heal(inv *runtime.Invocation) error {
	silent, target := parseSilentTarget(inv.Args, inv.Sender)
	if !silent {
		if target == inv.Sender {
			return inv.Reply("You have been healed.")
		}
		return inv.Reply(fmt.Sprintf("%s has been healed.", target))
	}
	return nil
}

func (g *Gameplay) slay(inv *runtime.Invocation) error {
	silent, target := parseSilentTarget(inv.Args, inv.Sender)
	if !silent {
		if target == inv.Sender {
			return inv.Reply("You have been slain.")
		}
		return inv.Reply(fmt.Sprintf("%s has been slain.", target))
	}
	return nil
}

func (g *Gameplay) locate(inv *runtime.Invocation) error {
	if len(inv.Args) == 0 {
		return inv.Reply("Compass now points to spawn.")
	}
	return inv.Reply(fmt.Sprintf("Compass now points to %s.", inv.Args[0]))
}

func (g *Gameplay) stack(inv *runtime.Invocation) error {
	return inv.Reply("Inventory organised and items stacked.")
}

// parseSilentTarget parses optional -s flag and optional player name from args.
func parseSilentTarget(args []string, defaultSender string) (silent bool, target string) {
	target = defaultSender
	for _, a := range args {
		if a == "-s" {
			silent = true
		} else {
			target = a
		}
	}
	return
}
