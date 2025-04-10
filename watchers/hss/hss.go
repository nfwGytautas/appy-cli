package watchers_hss

import (
	watchers_shared "github.com/nfwGytautas/appy-cli/watchers/shared"
)

func Watch() error {
	err := watchers_shared.WatchRepositories()
	if err != nil {
		return err
	}

	domainWatcher, err := watchers_shared.WatchDomain("domain/")
	if err != nil {
		return err
	}

	domainWatcher.Start()

	return nil
}
