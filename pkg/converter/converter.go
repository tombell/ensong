package converter

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/tombell/ensong/pkg/config"
)

type Converter struct {
	cfg      *config.Config
	filepath string
}

func New(cfg *config.Config, filepath string) *Converter {
	return &Converter{cfg, filepath}
}

func (c Converter) Convert(format string) (string, error) {
	file := strings.TrimSuffix(c.filepath, filepath.Ext(c.filepath))
	output := fmt.Sprintf("%s.%s", file, format)

	cmd := exec.Command("xld", "-f", format, "-o", output, c.filepath)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("convert failed: %w", err)
	}

	return output, nil
}
