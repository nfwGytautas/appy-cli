package watchers_hss

import (
	watchers_shared "github.com/nfwGytautas/appy-cli/watchers/shared"
)

func Watch() error {

	_, err := watchers_shared.WatchDomain("domain/")
	if err != nil {
		return err
	}

	return nil
}
