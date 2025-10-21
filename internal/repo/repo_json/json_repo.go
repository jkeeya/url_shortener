package repo_json

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/jkeeya/url_shortener/internal/service"
)

func NewJsonRepo(dataSource string) service.Repo {
	r := &RepoJson{
		dataFilePath: dataSource,
		content:      make(map[string]string),
	}

	b, err := os.ReadFile(dataSource)
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
	return "", nil
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

func (j *RepoJson) AddNewAlias(ctx context.Context, url, shortLink string) error {
	j.mu.Lock()
	if j.content == nil {
		j.content = make(map[string]string)
	}
	j.content[shortLink] = url

	data, err := json.MarshalIndent(j.content, "", "  ")
	j.mu.Unlock()
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(j.dataFilePath), 0o755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	tmp := j.dataFilePath + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("write tmp: %w", err)
	}
	if err := os.Rename(tmp, j.dataFilePath); err != nil {
		return fmt.Errorf("rename: %w", err)
	}
	fmt.Printf("writing to %s", j.dataFilePath) // или fmt.Println

	return nil
}
