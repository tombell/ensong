package monitor

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"

	"github.com/tombell/ensong/pkg/config"
)

type Monitor struct {
	cfg     *config.Config
	watcher *fsnotify.Watcher
}

func New(cfg *config.Config) (*Monitor, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("new watcher failed: %w", err)
	}

	m := &Monitor{cfg, watcher}

	if err := m.watcher.Add("WATCH"); err != nil {
		return nil, fmt.Errorf("watcher add failed: %w", err)
	}

	return m, nil
}

func (m Monitor) Run(ch chan error) {
	for {
		select {
		case event, ok := <-m.watcher.Events:
			if !ok || event.Op&fsnotify.Create != fsnotify.Create {
				continue
			}

			switch strings.ToLower(filepath.Ext(event.Name)) {
			case ".wav":
				// TODO: handle .wav file
			case ".txt":
				// TODO: handle .txt file
			}
		case err, ok := <-m.watcher.Errors:
			if !ok {
				continue
			}
		}
	}
}
