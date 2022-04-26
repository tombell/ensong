package converter

import (
	"fmt"
	"os/exec"

	"github.com/tombell/ensong/pkg/config"
)

type Converter struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Converter {
	return &Converter{cfg}
}

func (c Converter) ToAIFF() (string, error) {
	cmd := exec.Command("xld", "-f", "aif", "-o", "OUTPUT", "INPUT")

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("cmd run failed: %w", err)
	}

	return "OUTPUT", nil
}

func (c Converter) ToMP3() (string, error) {
	cmd := exec.Command("xld", "-f", "mp3", "-o", "OUTPUT", "INPUT")

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("cmd run failed: %w", err)
	}

	return "OUTPUT", nil
}
