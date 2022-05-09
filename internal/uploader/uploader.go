package uploader

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/go-github/v42/github"
	"golang.org/x/oauth2"

	"github.com/tombell/ensong/internal/config"
	"github.com/tombell/ensong/internal/metadata"
)

type Uploader struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Uploader {
	return &Uploader{cfg}
}

func (u Uploader) UploadMix(path string, meta *metadata.Metadata) error {
	body := &bytes.Buffer{}

	w := multipart.NewWriter(body)
	defer w.Close()

	w.WriteField("name", meta.Name())

	f1, _ := os.Open(path)
	defer f1.Close()

	mp3, _ := w.CreateFormFile("mp3", "file.mp3")
	io.Copy(mp3, f1)

	f2, _ := os.Open(meta.Picture())
	defer f2.Close()

	picture, _ := w.CreateFormFile("picture", "file.png")
	io.Copy(picture, f2)

	for idx, tag := range meta.Tags {
		w.WriteField(fmt.Sprintf("tags-%d-tag", idx), tag)
	}

	url := fmt.Sprintf("https://api.mixcloud.com/upload/?access_token=%s", u.cfg.Mixcloud.Token)

	req, _ := http.NewRequest("POST", url, body)
	req.Header.Add("Content-Type", w.FormDataContentType())

	client := &http.Client{}
	if _, err := client.Do(req); err != nil {
		return fmt.Errorf("http request failed: %w", err)
	}

	return nil
}

func (u Uploader) UploadTracklist(filePath string) error {
	filename := filepath.Base(filePath)
	file := strings.TrimSuffix(filename, path.Ext(filename))

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("os read file failed: %w", err)
	}

	date, err := time.Parse("2006-01-02", file)
	if err != nil {
		return fmt.Errorf("time parse failed: %w", err)
	}

	opts := &github.RepositoryContentFileOptions{
		Message: github.String(fmt.Sprintf("Add %s", filename)),
		Content: data,
		Branch:  github.String(u.cfg.GitHub.Repo.Branch),
		Committer: &github.CommitAuthor{
			Name:  github.String(u.cfg.GitHub.Committer.Name),
			Email: github.String(u.cfg.GitHub.Committer.Email),
		},
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: u.cfg.GitHub.Token})
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)

	_, _, err = client.Repositories.CreateFile(
		context.Background(),
		u.cfg.GitHub.Repo.Owner,
		u.cfg.GitHub.Repo.Name,
		fmt.Sprintf("%d/%s", date.Year(), filename),
		opts,
	)
	if err != nil {
		return fmt.Errorf("repositories create file failed: %w", err)
	}

	return nil
}
