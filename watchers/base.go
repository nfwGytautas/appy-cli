package watchers

import (
	"fmt"

	"github.com/nfwGytautas/appy-cli/config"
	"github.com/nfwGytautas/appy-cli/shared"
	watchers_hmd "github.com/nfwGytautas/appy-cli/watchers/hmd"
	watchers_hss "github.com/nfwGytautas/appy-cli/watchers/hss"
	watchers_shared "github.com/nfwGytautas/appy-cli/watchers/shared"
)

func Watch() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	fmt.Println(cfg.Type)
	fmt.Println("Starting watchers...")

	err = watchers_shared.WatchConfig()
	if err != nil {
		return err
	}

	switch cfg.Type {
	case shared.ScaffoldHMD:
		err := watchers_hmd.Watch()
		if err != nil {
			return err
		}
	case shared.ScaffoldHSS:
		err := watchers_hss.Watch()
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown project type: %s", cfg.Type)
	}

	// Block
	<-make(chan struct{})

	return nil
}
