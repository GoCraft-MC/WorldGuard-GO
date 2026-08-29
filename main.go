package main

import (
	"log"
	"os"
	"strings"

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

	store, err := region.NewStore(dataDir + "/regions.json")
	if err != nil {
		log.Fatalf("loading region store: %v", err)
	}

	p := &plugin{
		protect: handler.NewProtect(store),
		cmd:     handler.NewCommand(store),
	}

	if err := runtime.Run(socketPath, p); err != nil {
		log.Fatalf("plugin exited: %v", err)
	}
}

type plugin struct {
	protect *handler.Protect
	cmd     *handler.Command
}

func (p *plugin) OnEvent(e *runtime.Event) runtime.Verdict {
	switch e.Type {
	case "block.break":
		return p.protect.OnBlockBreak(e)
	case "block.place":
		return p.protect.OnBlockPlace(e)
	}
	return runtime.Allow()
}

func (p *plugin) OnCommand(inv *runtime.Invocation) error {
	if strings.HasPrefix(inv.Command, "region") {
		return p.cmd.Invoke(inv)
	}
	return nil
}
