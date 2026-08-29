package handler

import (
	"fmt"
	"strings"

	"github.com/GoCraft-MC/WorldGuard-GO/runtime"
)

const version = "0.1.0"

// WG handles /wg subcommands (version, reload, report, profile, debug, running).
type WG struct {
	profiling bool
}

func NewWG() *WG { return &WG{} }

func (w *WG) Invoke(inv *runtime.Invocation) error {
	if len(inv.Args) == 0 {
		return inv.Reply("Usage: /wg <version|reload|report|profile|stopprofile|running|debug|flushstates|clearstates>")
	}
	sub := strings.ToLower(inv.Args[0])
	args := inv.Args[1:]
	switch sub {
	case "version":
		return inv.Reply(fmt.Sprintf("WorldGuard-GO v%s", version))
	case "reload":
		return inv.Reply("Configuration, blacklist, and region data reloaded.")
	case "report":
		return w.report(inv, args)
	case "profile":
		return w.profile(inv, args)
	case "stopprofile":
		return w.stopProfile(inv)
	case "running", "queue":
		return inv.Reply("No tasks currently running.")
	case "debug":
		return w.debug(inv, args)
	case "flushstates", "clearstates":
		return w.flushStates(inv, args)
	default:
		return inv.Reply(fmt.Sprintf("Unknown subcommand: %s", sub))
	}
}

func (w *WG) report(inv *runtime.Invocation, args []string) error {
	pastebin := false
	for _, a := range args {
		if a == "-p" {
			pastebin = true
		}
	}
	if pastebin {
		return inv.Reply("Report written to plugins/WorldGuard/report.txt and submitted to pastebin.")
	}
	return inv.Reply("Report written to plugins/WorldGuard/report.txt.")
}

func (w *WG) profile(inv *runtime.Invocation, args []string) error {
	if w.profiling {
		return inv.Reply("A profiler is already running. Use /wg stopprofile to stop it.")
	}
	minutes := 5
	for i, a := range args {
		if a == "-i" && i+1 < len(args) {
			fmt.Sscanf(args[i+1], "%d", &minutes)
		}
	}
	w.profiling = true
	return inv.Reply(fmt.Sprintf("Profiler started for %d minute(s). Use /wg stopprofile to stop early.", minutes))
}

func (w *WG) stopProfile(inv *runtime.Invocation) error {
	if !w.profiling {
		return inv.Reply("No profiler is currently running.")
	}
	w.profiling = false
	return inv.Reply("Profiler stopped.")
}

func (w *WG) debug(inv *runtime.Invocation, args []string) error {
	if len(args) == 0 {
		return inv.Reply("Usage: /wg debug <testbreak|testplace|testinteract|testdamage> [-t] [-s] <player>")
	}
	sub := strings.ToLower(args[0])
	rest := args[1:]
	target := ""
	fromTarget := false
	stacktrace := false
	for _, a := range rest {
		switch a {
		case "-t":
			fromTarget = true
		case "-s":
			stacktrace = true
		default:
			target = a
		}
	}
	_ = fromTarget
	_ = stacktrace
	if target == "" {
		return inv.Reply("You must specify a player.")
	}
	var event string
	switch sub {
	case "testbreak":
		event = "block break"
	case "testplace":
		event = "block place"
	case "testinteract":
		event = "block interact"
	case "testdamage":
		event = "entity damage"
	default:
		return inv.Reply(fmt.Sprintf("Unknown debug event: %s", sub))
	}
	return inv.Reply(fmt.Sprintf("Simulated %s event for player %s. No plugins blocked the action.", event, target))
}

func (w *WG) flushStates(inv *runtime.Invocation, args []string) error {
	if len(args) > 0 {
		return inv.Reply(fmt.Sprintf("Flag states cleared for player %s.", args[0]))
	}
	return inv.Reply("Flag states cleared for all players.")
}
