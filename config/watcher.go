package config

import (
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/nfwGytautas/appy-cli/utils"
)

type Watcher struct {
	Fs          []string `yaml:"fs"`
	EventFilter []string `yaml:"eventFilter"`
	Actions     []string `yaml:"actions"`

	provider *Provider         `yaml:"-"`
	watcher  *fsnotify.Watcher `yaml:"-"`
}

func (w *Watcher) Configure(opts RepositoryConfigureOpts) error {
	for it, fs := range w.Fs {
		w.Fs[it] = w.provider.ApplyStringSubstitution(fs)
	}

	for it, action := range w.Actions {
		w.Actions[it] = w.provider.ApplyStringSubstitution(action)
	}

	return nil
}

func (w *Watcher) Stop() error {
	if w.watcher != nil {
		return w.watcher.Close()
	}

	return nil
}

func (w *Watcher) Start() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	for _, fs := range w.Fs {
		err = watcher.Add(fs)
		if err != nil {
			return err
		}
		utils.Console.DebugLn("attached watcher to %s", fs)
	}

	w.watcher = watcher

	go func() {
		for event := range w.watcher.Events {
			w.onFileEvent(event)
		}
	}()

	return nil
}

func (w *Watcher) Restart() error {
	err := w.Stop()
	if err != nil {
		return err
	}

	err = w.Start()
	if err != nil {
		return err
	}

	return nil
}

func (w *Watcher) onFileEvent(event fsnotify.Event) {
	if event.Op.Has(fsnotify.Chmod) {
		return
	}

	if len(w.EventFilter) > 0 {
		matched := false

		for _, filter := range w.EventFilter {
			if filter == event.Op.String() {
				matched = true
				break
			}
		}

		if !matched {
			return
		}
	}

	utils.Console.DebugLn("%s", event.String())

	w.runActions()
}

func (w *Watcher) runActions() {
	cwd, err := os.Getwd()
	if err != nil {
		utils.Console.ErrorLn("failed to get current working directory: %v", err)
		return
	}

	for _, action := range w.Actions {
		err := utils.RunCommand(cwd, action)
		if err != nil {
			utils.Console.ErrorLn("failed to run hook action `%s`: %v", action, err)
			continue
		}
	}
}
