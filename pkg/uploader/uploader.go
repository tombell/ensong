package uploader

import (
	"github.com/tombell/ensong/pkg/config"
)

type Uploader struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Uploader {
	return &Uploader{cfg}
}

func (u Uploader) UploadMix() error {
	return nil
}

func (u Uploader) UploadTracklist() error {
	return nil
}
