package utils

import "github.com/fsnotify/fsnotify"

type OnFileEvent func(event fsnotify.Event)

type Watcher struct {
	watcher     *fsnotify.Watcher
	onFileEvent OnFileEvent
}

func NewWatcher(directory string, onFileEvent OnFileEvent) (*Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	err = watcher.Add(directory)
	if err != nil {
		return nil, err
	}

	return &Watcher{
		watcher:     watcher,
		onFileEvent: onFileEvent,
	}, nil
}

func (w *Watcher) Start() {
	go func() {
		for event := range w.watcher.Events {
			w.onFileEvent(event)
		}
	}()
}

func (w *Watcher) Stop() {
	w.watcher.Close()
}
