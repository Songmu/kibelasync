package kibela

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	"golang.org/x/xerrors"
)

type md struct {
	ID          ID
	FrontMatter *meta
	Content     string
	UpdatedAt   time.Time
}

type meta struct {
	Title     string   `json:"title"`
	CoEditing bool     `json:"coediting"`
	Folder    string   `json:"folder"`
	Groups    []string `json:"groups"`
	Author    string   `json:"author"`
}

func (m *md) fullContent() string {
	fm, _ := yaml.Marshal(m.FrontMatter)

	c := strings.Join([]string{"---", string(fm) + "---", "", m.Content}, "\n")
	if !strings.HasSuffix(c, "\n") {
		// fill newline for suppressing warning of "No newline at end of file"
		c += "\n"
	}
	return c
}

func (m *md) save() error {
	stuff := strings.Split(m.ID.String(), "/")
	if len(stuff) != 2 {
		return fmt.Errorf("invalid id: %s", string(m.ID))
	}
	idNum, err := m.ID.Number()
	if err != nil {
		return xerrors.Errorf("failed to save Markdown: %w", err)
	}
	basePath := "." // XXX
	savePath := filepath.Join(basePath, "notes", fmt.Sprintf("%d.md", idNum))
	if err := os.MkdirAll(filepath.Dir(savePath), 0755); err != nil {
		return xerrors.Errorf("failed to save Markdown: %w", err)
	}
	if err := ioutil.WriteFile(savePath, []byte(m.fullContent()), 0644); err != nil {
		return xerrors.Errorf("failed to save Markdown: %w", err)
	}
	if !m.UpdatedAt.IsZero() {
		if err := os.Chtimes(savePath, m.UpdatedAt, m.UpdatedAt); err != nil {
			return xerrors.Errorf("failed to set mtime to Markdown: %w", err)
		}
	}
	return nil
}
