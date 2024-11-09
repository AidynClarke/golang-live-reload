package watcher

import (
	"log"
	"os"
	"path/filepath"

	"github.com/AidynClarke/golang-live-reload/buffer"
	"github.com/AidynClarke/golang-live-reload/utils"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	watcher *fsnotify.Watcher
	buffer *buffer.Buffer

	dir string
	recursive bool
	eventHandler func(event fsnotify.Event)
	errorHandler func(err error)

	exclude []string
}

func NewWatcher(dir string, recursive bool, eventHandler func(event fsnotify.Event), errorHandler func(err error), exclude []string) (*Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	excludePatterns := make([]string, len(exclude))
	for i, pattern := range exclude {
		excludePatterns[i] = utils.Normalise(pattern)
	}

	return &Watcher{
		watcher: watcher,

		dir: dir,
		recursive: recursive,
		eventHandler: eventHandler,
		errorHandler: errorHandler,
		buffer: buffer.NewBuffer(eventHandler),

		exclude: excludePatterns,
	}, nil
}

func (w *Watcher) Watch() error {
	if w.recursive {
		err := w.watchRecursive()
		if err != nil {
			return err
		}
	} else {
		err := w.watcher.Add(w.dir)
		if err != nil {
			return err
		}
	}
	// Start listening for events.
	go func() {
			for {
					select {
					case event, ok := <-w.watcher.Events:
							if !ok {
									return
							}

							w.buffer.NewEvent(event)
					case err, ok := <-w.watcher.Errors:
							if !ok {
									return
							}
							w.errorHandler(err)
					}
			}
	}()

	return nil
}

func (w *Watcher) watchRecursive() error {
	globMatcher := NewGlobMatcher(w.exclude)

	return filepath.Walk(w.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		normalisedPath := utils.Normalise(path)

		if globMatcher.Match(normalisedPath) {
			return filepath.SkipDir
		}

		if info.IsDir() {
			err := w.watcher.Add(path)
			if err != nil {
				return err
			}
			log.Println("Watching directory: ", normalisedPath)
		}
		return nil
	})
}

