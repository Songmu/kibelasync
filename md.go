package kibela

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Songmu/kibela/client"
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
	Folder    string   `json:"folder,omitempty"`
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
	if m.filepath == "" {
		basePath := m.dir
		if basePath == "" {
			basePath = "notes"
		}
		m.filepath = filepath.Join(basePath, fmt.Sprintf("%d.md", idNum))
	}
	if err := os.MkdirAll(filepath.Dir(m.filepath), 0755); err != nil {
		return xerrors.Errorf("failed to save Markdown: %w", err)
	}
	if err := ioutil.WriteFile(m.filepath, []byte(m.fullContent()), 0644); err != nil {
		return xerrors.Errorf("failed to save Markdown: %w", err)
	}
	if !m.UpdatedAt.IsZero() {
		if err := os.Chtimes(m.filepath, m.UpdatedAt, m.UpdatedAt); err != nil {
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
	num, err := strconv.Atoi(stuffs[0])
	if err != nil {
		return nil, fmt.Errorf("invalid filename (must be [0-9]+.md): %s", fname)
	}

	f, err := os.Open(fpath)
	if err != nil {
		return nil, xerrors.Errorf("failed to load md: %s, %w", fpath, err)
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return nil, xerrors.Errorf("failed to load md: %w", err)
	}
	m := &md{
		ID:        newID("Blog", num),
		UpdatedAt: fi.ModTime(),
		filepath:  fpath,
	}
	if err := m.loadContentFromReader(f, true); err != nil {
		return nil, xerrors.Errorf("failed to loadMD: %w", err)
	}
	return m, nil
}

func (m *md) loadContentFromReader(r io.Reader, forceFrontmatter bool) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return xerrors.Errorf("failed to load md: %w", err)
	}
	if m.FrontMatter == nil {
		m.FrontMatter = &meta{}
	}
	contents := strings.SplitN(string(b), "---\n", 3)
	if len(contents) == 3 && contents[0] == "" {
		if err := yaml.Unmarshal([]byte(contents[1]), m.FrontMatter); err != nil {
			if forceFrontmatter {
				return xerrors.Errorf("invalid frontmatter: %w", err)
			}
			m.FrontMatter.Title, m.Content = detectTitle(string(b))
		} else {
			if m.FrontMatter.Title == "" {
				m.FrontMatter.Title, m.Content = detectTitle(contents[2])
			} else {
				m.Content = strings.TrimSpace(contents[2]) + "\n"
			}
		}
	} else if !forceFrontmatter {
		m.FrontMatter.Title, m.Content = detectTitle(string(b))
	} else {
		return fmt.Errorf("invalid contents of md: %s", string(b))
	}
	return nil
}

var detectTitleReg = regexp.MustCompile(`\A` +
	`(?:` +
	`([^\r\n]+)\r?\n={2,}` + // underlined title: ex. "Title Content\n====="
	`|` +
	`#\s+([^\r\n]+)` + // hashed title: ex. "# Title Content"
	`)` +
	`\r?\n`)

func detectTitle(rawContent string) (title, content string) {
	rawContent = strings.TrimSpace(rawContent) + "\n"
	m := detectTitleReg.FindStringSubmatch(rawContent)
	if len(m) < 3 {
		return "", rawContent
	}
	title = m[1]
	if title == "" {
		title = m[2]
	}
	content = strings.TrimSpace(strings.TrimPrefix(rawContent, m[0])) + "\n"
	return title, content
}

func (m *md) toNote() *note {
	groups := make([]*group, len(m.FrontMatter.Groups))
	for i, g := range m.FrontMatter.Groups {
		groups[i] = &group{Name: g}
	}
	return &note{
		ID:        m.ID,
		Title:     m.FrontMatter.Title,
		Content:   m.Content,
		CoEditing: m.FrontMatter.CoEditing,
		Folder:    m.FrontMatter.Folder,
		Groups:    groups,
		Author: struct {
			Account string `json:"account"`
		}{
			Account: m.FrontMatter.Author,
		},
		UpdatedAt: Time{Time: m.UpdatedAt},
	}
}

func (ki *kibela) pushMD(m *md) error {
	n := m.toNote()
	if err := ki.pushNote(n); err != nil {
		return xerrors.Errorf("failed to pushMD: %w", err)
	}
	return os.Chtimes(m.filepath, n.UpdatedAt.Time, n.UpdatedAt.Time)
}

func (ki *kibela) publishMD(m *md, save bool) error {
	if m.FrontMatter == nil {
		m.FrontMatter = &meta{}
	}
	groupIDs := make([]string, len(m.FrontMatter.Groups))
	for i, g := range m.FrontMatter.Groups {
		id, err := ki.fetchGroupID(g)
		if err != nil {
			return xerrors.Errorf("failed to publishMD: %w", err)
		}
		groupIDs[i] = string(id)
	}
	sort.Strings(groupIDs)
	input := &noteInput{
		Title:     m.FrontMatter.Title,
		Content:   m.Content,
		Folder:    m.FrontMatter.Folder,
		CoEditing: m.FrontMatter.CoEditing,
		GroupIDs:  groupIDs,
	}
	data, err := ki.cli.Do(&client.Payload{
		Query: createNoteMutation,
		Variables: struct {
			Input *noteInput `json:"input"`
		}{
			Input: input,
		},
	})
	if err != nil {
		return xerrors.Errorf("failed to publishNote while accessing remote: %w", err)
	}
	var res struct {
		CreateNote struct {
			Note *note `json:"note"`
		} `json:"createNote"`
	}
	if err := json.Unmarshal(data, &res); err != nil {
		return xerrors.Errorf("failed to ki.publishNote while unmarshaling response: %w", err)
	}
	if !save {
		return nil
	}
	n := res.CreateNote.Note
	groups := make([]string, len(n.Groups))
	for i, g := range n.Groups {
		groups[i] = g.Name
	}
	m.FrontMatter.Groups = groups
	m.ID = n.ID
	m.UpdatedAt = n.UpdatedAt.Time
	m.FrontMatter.Author = n.Author.Account
	origFilePath := m.filepath
	m.filepath = ""
	if err := m.save(); err != nil {
		return xerrors.Errorf("failed to publishMD. publish succeeded but failed to store file: %w", err)
	}
	if origFilePath != "" {
		if err := os.RemoveAll(origFilePath); err != nil {
			return xerrors.Errorf("failed to publishMD while cleanup orginal MD: %w", err)
		}
	}
	return nil
}
