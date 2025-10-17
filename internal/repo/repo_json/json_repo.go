package repo_json

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"

	"github.com/jkeeya/url_shortener/internal/service"
)

func NewJsonRepo() service.Repo {
	exe, _ := os.Executable()
	dir := filepath.Dir(exe)
	dataFile := filepath.Join(dir, "internal/repo/repo_json/data.json")

	r := &RepoJson{
		dataFilePath: dataFile,
		content:      make(map[string]string),
	}

	b, err := os.ReadFile(dataFile)
	if err == nil && len(b) > 0 {
		_ = json.Unmarshal(b, &r.content)
	}

	return r
}

type RepoJson struct {
	dataFilePath string
	content      map[string]string
	mu           sync.RWMutex
}

func (j *RepoJson) FindByURL(ctx context.Context, url string) (string, error) {
	j.mu.RLock()
	defer j.mu.RUnlock()

	for short, checkedURL := range j.content {
		if checkedURL == url {
			return short, nil
		}
	}
	return "", errors.New("запрашиваемого URL нет в базе")
}

func (j *RepoJson) FindByShortLink(ctx context.Context, shortLink string) (string, error) {
	j.mu.RLock()
	defer j.mu.RUnlock()

	url, ok := j.content[shortLink]
	if !ok {
		return "", errors.New("что-то пошло не так: краткой ссылки для данного URL не существует")
	} else {
		return url, nil
	}
}

func (j *RepoJson) AddNewAlias(ctx context.Context, url string, shortLink string) error {
	j.mu.Lock()
	j.content[shortLink] = url
	file, _ := os.Open(j.dataFilePath)
	defer file.Close()

	data, err := json.MarshalIndent(j.content, "", "  ")
	j.mu.Unlock()
	if err != nil {
		return err
	}

	tmp := j.dataFilePath + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, j.dataFilePath)
}
