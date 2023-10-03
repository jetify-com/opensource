package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"go.jetpack.io/pkg/sandbox/vm"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	vm := vm.VM{}
	flag.StringVar(&vm.HostDataDir, "datadir", ".devbox/vm", "`path` to the directory for saving VM state")
	flag.BoolVar(&vm.Install, "install", false, "mount NixOS install image")
	flag.Parse()

	if vm.Install {
		slog.Debug("downloading the NixOS installer, this make take a few minutes")
	}

	err := vm.Start(ctx)
	if err != nil {
		slog.Error("start virtual machine", "err", err)
		os.Exit(1)
	}

	<-ctx.Done()
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := vm.Stop(ctx); err != nil {
		slog.Error("stop virtual machine", "err", err)
	}
}
