package watchers

import (
	"fmt"

	"github.com/nfwGytautas/appy-cli/shared"
	watchers_hmd "github.com/nfwGytautas/appy-cli/watchers/hmd"
	watchers_hss "github.com/nfwGytautas/appy-cli/watchers/hss"
)

func Watch() error {
	cfg, err := shared.LoadConfig()
	if err != nil {
		return err
	}

	fmt.Println(cfg.Type)
	fmt.Println("Starting watcher...")

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
