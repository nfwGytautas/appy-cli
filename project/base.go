package project

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/nfwGytautas/appy-cli/config"
	project_hmd "github.com/nfwGytautas/appy-cli/project/hmd"
	project_shared "github.com/nfwGytautas/appy-cli/project/shared"
	"github.com/nfwGytautas/appy-cli/shared"
	"github.com/nfwGytautas/appy-cli/utils"
)

var projectTypes = map[string]func(context.Context, *sync.WaitGroup) error{
	shared.ScaffoldHMD: project_hmd.Watch,
}

func Watch(ctx context.Context) error {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(ctx)

	defer wg.Wait()
	defer cancel()

	cfg := config.GetConfig()

	utils.Console.DebugLn(cfg.Type)
	utils.Console.DebugLn("Starting watchers...")

	err := project_shared.WatchConfig(ctx, &wg)
	if err != nil {
		return err
	}

	projectRun, exists := projectTypes[cfg.Type]
	if !exists {
		return fmt.Errorf("unknown project type: %s", cfg.Type)
	}

	err = projectRun(ctx, &wg)
	if err != nil {
		return err
	}

	// Block until a signal is received
	<-stop
	utils.Console.ClearLines(1)
	utils.Console.DebugLn("Received signal, shutting down...")

	return nil
}
