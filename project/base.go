package project

import (
	"fmt"

	"github.com/nfwGytautas/appy-cli/config"
	project_hmd "github.com/nfwGytautas/appy-cli/project/hmd"
	project_shared "github.com/nfwGytautas/appy-cli/project/shared"
	"github.com/nfwGytautas/appy-cli/shared"
	"github.com/nfwGytautas/appy-cli/utils"
)

var projectTypes = map[string]func() error{
	shared.ScaffoldHMD: project_hmd.Watch,
}

func Watch() error {
	cfg := config.GetConfig()

	utils.Console.DebugLn(cfg.Type)
	utils.Console.DebugLn("Starting watchers...")

	err := project_shared.WatchConfig()
	if err != nil {
		return err
	}

	projectRun, exists := projectTypes[cfg.Type]
	if !exists {
		return fmt.Errorf("unknown project type: %s", cfg.Type)
	}

	err = projectRun()
	if err != nil {
		return err
	}

	// Block
	<-make(chan struct{})

	return nil
}
