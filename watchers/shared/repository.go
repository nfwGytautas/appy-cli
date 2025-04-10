package watchers_shared

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/nfwGytautas/appy-cli/utils"
)

func WatchRepositories() error {
	repositoryWatcher, err := utils.NewWatcher("repositories/", onRepositoryEvent)
	if err != nil {
		return err
	}

	repositoryWatcher.Start()

	return nil
}

func onRepositoryEvent(event fsnotify.Event) {
	fmt.Println("Repository:", event.Name, event.Op)
}
