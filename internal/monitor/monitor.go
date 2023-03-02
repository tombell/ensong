package monitor

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"

	"github.com/tombell/ensong/internal/config"
	"github.com/tombell/ensong/internal/converter"
	"github.com/tombell/ensong/internal/metadata"
	"github.com/tombell/ensong/internal/uploader"
)

type Monitor struct {
	cfg     *config.Config
	logger  *log.Logger
	watcher *fsnotify.Watcher
}

func New(cfg *config.Config, logger *log.Logger) (*Monitor, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("new watcher failed: %w", err)
	}

	monitor := &Monitor{cfg, logger, watcher}

	for _, dir := range monitor.cfg.Directories {
		logger.Printf("watching directory %s...\n", dir)

		if err := monitor.watcher.Add(dir); err != nil {
			return nil, fmt.Errorf("watcher add failed: %w", err)
		}
	}

	return monitor, nil
}

func (m Monitor) Run() {
	for {
		select {
		case event, ok := <-m.watcher.Events:
			if !ok || event.Op&fsnotify.Create != fsnotify.Create {
				continue
			}

			switch strings.ToLower(filepath.Ext(event.Name)) {
			case ".wav":
				m.handleMix(event.Name)
			case ".txt":
				m.handleTracklist(event.Name)
			}
		case err, ok := <-m.watcher.Errors:
			if !ok {
				continue
			}

			m.logger.Printf("watcher error: %s", err)
		}
	}
}

func (m Monitor) handleMix(name string) {
	m.logger.Printf("handling %s...\n", name)

	conv := converter.New(m.cfg, name)

	m.logger.Println("converting from wav to aiff...")

	aiff, err := conv.Convert("aiff")
	if err != nil {
		m.logger.Printf("converter error: failed to convert to aiff: %s", err)
		return
	}

	m.logger.Println("moving aiff file to backup directory...")

	if err = os.Rename(aiff, fmt.Sprintf("%s/%s", m.cfg.Backup, filepath.Base(aiff))); err != nil {
		m.logger.Printf("error: failed to rename file: %s", err)
		return
	}

	m.logger.Println("converting from wav to mp3...")

	mp3, err := conv.Convert("mp3")
	if err != nil {
		m.logger.Printf("converter error: failed to convert to mp3: %s", err)
		return
	}

	m.logger.Println("parsing metadata...")

	meta, err := metadata.New(m.cfg, name)
	if err != nil {
		m.logger.Printf("metadata error: failed to parse metadata: %s", err)
		return
	}

	m.logger.Println("uploading mix to mixcloud...")

	u := uploader.New(m.cfg)
	if err := u.UploadMix(mp3, meta); err != nil {
		m.logger.Printf("uploader error: failed to upload mix: %s", err)
		return
	}

	m.logger.Println("cleaning up files...")

	for _, file := range []string{name, mp3} {
		if err := os.Remove(file); err != nil {
			m.logger.Printf("error: failed to remove file: %s", err)
			return
		}
	}

	m.logger.Printf("finished handling %s...\n", name)
}

func (m Monitor) handleTracklist(name string) {
	m.logger.Printf("handling %s...\n", name)
	m.logger.Printf("uploading tracklist to github: %s...\n", name)

	u := uploader.New(m.cfg)
	if err := u.UploadTracklist(name); err != nil {
		m.logger.Printf("uploader error: failed to upload tracklist: %s", err)
		return
	}

	m.logger.Printf("finished handling %s...\n", name)
}
