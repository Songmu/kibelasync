package kibela

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
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

	dir, filepath string
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
	savePath := m.filepath
	if savePath == "" {
		basePath := m.dir
		if basePath == "" {
			basePath = "notes"
		}
		savePath = filepath.Join(basePath, fmt.Sprintf("%d.md", idNum))
	}
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

func loadMD(fpath string) (*md, error) {
	fname := filepath.Base(fpath)
	stuffs := strings.Split(fname, ".")
	if len(stuffs) != 2 {
		return nil, fmt.Errorf("invalid filename (must be [0-9]+.md): %s", fname)
	}
	if stuffs[1] != "md" {
		return nil, fmt.Errorf("invalid filename (must be [0-9]+.md): %s", fname)
	}
	if _, err := strconv.Atoi(stuffs[0]); err != nil {
		return nil, fmt.Errorf("invalid filename (must be [0-9]+.md): %s", fname)
	}
	mdNumStr := stuffs[0]

	f, err := os.Open(fpath)
	if err != nil {
		return nil, xerrors.Errorf("failed to load md: %s, %w", fpath, err)
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, xerrors.Errorf("failed to load md: %s, %w", fpath, err)
	}
	contents := strings.SplitN(string(b), "---\n", 3)
	if len(contents) != 3 || contents[0] != "" {
		return nil, fmt.Errorf("invalid contents of md: %s", string(b))
	}
	var me meta
	if err := yaml.Unmarshal([]byte(contents[1]), &me); err != nil {
		return nil, xerrors.Errorf("invalid frontmatter: %w", err)
	}
	fi, err := f.Stat()
	if err != nil {
		return nil, xerrors.Errorf("failed to load md: %w", err)
	}
	return &md{
		ID:          newID(fmt.Sprintf("Blog/%s", mdNumStr)),
		FrontMatter: &me,
		Content:     strings.TrimSpace(contents[2]) + "\n",
		UpdatedAt:   fi.ModTime(),
		filepath:    fpath,
	}, nil
}
