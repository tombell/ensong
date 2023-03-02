package metadata

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/tombell/ensong/internal/config"
	"github.com/tombell/ensong/internal/tracklists"
)

type Metadata struct {
	cfg *config.Config

	Date   *time.Time
	Tags   []string
	Tracks [][]string
}

func New(cfg *config.Config, filename string) (*Metadata, error) {
	file := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
	parts := strings.Split(file, "|")

	date, err := time.Parse("2006-01-02", parts[0])
	if err != nil {
		return nil, fmt.Errorf("time parse failed: %w", err)
	}

	tracklistFile := fmt.Sprintf("%s/%s.txt", filepath.Dir(filename), date.Format("2006-01-02"))

	tracks, err := tracklists.Read(tracklistFile)
	if err != nil {
		return nil, fmt.Errorf("tracklist read failed: %w", err)
	}

	return &Metadata{cfg, &date, parts[1:], tracks}, nil
}

func (m Metadata) Name() string {
	var name string

	switch m.Date.Weekday() {
	case time.Monday:
		name += m.cfg.Metadata.Names.Monday
	case time.Tuesday:
		name += m.cfg.Metadata.Names.Tuesday
	case time.Wednesday:
		name += m.cfg.Metadata.Names.Wednesday
	case time.Thursday:
		name += m.cfg.Metadata.Names.Thursday
	case time.Friday:
		name += m.cfg.Metadata.Names.Friday
	case time.Saturday:
		name += m.cfg.Metadata.Names.Saturday
	case time.Sunday:
		name += m.cfg.Metadata.Names.Sunday
	}

	name += fmt.Sprintf(" (%s)", m.Date.Format("02/01/2006"))

	return name
}

func (m Metadata) Picture() string {
	switch m.Date.Weekday() {
	case time.Monday:
		return m.cfg.Metadata.Pictures.Monday
	case time.Tuesday:
		return m.cfg.Metadata.Pictures.Tuesday
	case time.Wednesday:
		return m.cfg.Metadata.Pictures.Wednesday
	case time.Thursday:
		return m.cfg.Metadata.Pictures.Thursday
	case time.Friday:
		return m.cfg.Metadata.Pictures.Friday
	case time.Saturday:
		return m.cfg.Metadata.Pictures.Saturday
	case time.Sunday:
		return m.cfg.Metadata.Pictures.Sunday
	default:
		return ""
	}
}
