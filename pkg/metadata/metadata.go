package metadata

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/tombell/ensong/pkg/config"
)

type Metadata struct {
	cfg      *config.Config
	filename string
	date     *time.Time

	Tags []string
}

func New(cfg *config.Config, filename string) (*Metadata, error) {
	parts := strings.Split(filename, "|")

	date, err := time.Parse("2006-01-02", parts[0])
	if err != nil {
		return nil, fmt.Errorf("time parse failed: %w", err)
	}

	meta := &Metadata{
		cfg:      cfg,
		filename: filename,
		date:     &date,

		Tags: parts[1:],
	}

	return meta, nil
}

func (m Metadata) Name() (string, error) {
	var name string

	switch m.date.Weekday() {
	case time.Friday:
		name += "IAMDJRIFF pres. The Weekend Warmup"
	case time.Saturday, time.Sunday:
		name += "IAMDJRIFF pres. The Weekender"
	default:
		return "", errors.New("unsupported day of the week")
	}

	name += fmt.Sprintf(" (%s)", m.date.Format("02/01/2006"))

	return name, nil
}
