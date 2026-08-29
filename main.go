package main

import (
	"log"
	"os"
	"strings"

	"github.com/GoCraft-MC/WorldGuard-GO/config"
	"github.com/GoCraft-MC/WorldGuard-GO/handler"
	"github.com/GoCraft-MC/WorldGuard-GO/region"
	"github.com/GoCraft-MC/WorldGuard-GO/runtime"
)

func main() {
	socketPath := os.Getenv("GOCRAFT_SOCKET")
	if socketPath == "" {
		log.Fatal("GOCRAFT_SOCKET is not set")
	}
	dataDir := os.Getenv("GOCRAFT_DATA")
	if dataDir == "" {
		dataDir = "data"
	}

	cfg, err := config.Load(dataDir + "/config.yml")
	if err != nil {
		log.Fatalf("loading config: %v", err)
	}

	store, err := region.NewStore(dataDir + "/regions.json")
	if err != nil {
		log.Fatalf("loading region store: %v", err)
	}

	p := &plugin{
		cfg:       cfg,
		protect:   handler.NewProtect(store),
		regionCmd: handler.NewCommand(store),
		gameplay:  handler.NewGameplay(),
		emergency: handler.NewEmergency(),
		wg:        handler.NewWG(),
	}

	if cfg.SummaryOnStart {
		log.Printf("WorldGuard-GO v0.1.0 loaded. Regions enabled: %v", cfg.Regions.Enable)
	}

	if err := runtime.Run(socketPath, p); err != nil {
		log.Fatalf("plugin exited: %v", err)
	}
}

type plugin struct {
	cfg       *config.Config
	protect   *handler.Protect
	regionCmd *handler.Command
	gameplay  *handler.Gameplay
	emergency *handler.Emergency
	wg        *handler.WG
}

func (p *plugin) OnEvent(e *runtime.Event) runtime.Verdict {
	if !p.cfg.Regions.Enable {
		return runtime.Allow()
	}
	switch e.Type {
	case "block.break":
		return p.protect.OnBlockBreak(e)
	case "block.place":
		return p.protect.OnBlockPlace(e)
	}
	return runtime.Allow()
}

func (p *plugin) OnCommand(inv *runtime.Invocation) error {
	cmd := strings.ToLower(inv.Command)
	switch {
	case cmd == "region" || cmd == "rg":
		return p.regionCmd.Invoke(inv)
	case cmd == "god" || cmd == "ungod" || cmd == "heal" || cmd == "slay" || cmd == "locate" || cmd == "stack":
		return p.gameplay.Invoke(inv)
	case cmd == "stopfire" || cmd == "allowfire" || cmd == "stoplag" || cmd == "halt-activity" || cmd == "haltactivity":
		return p.emergency.Invoke(inv)
	case cmd == "wg":
		return p.wg.Invoke(inv)
	}
	return nil
}
