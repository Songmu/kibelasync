package kibela

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Songmu/kibelasync/client"
	"golang.org/x/xerrors"
	"gopkg.in/yaml.v2"
)

type MD struct {
	ID          ID
	FrontMatter *Meta
	Content     string
	UpdatedAt   time.Time

	dir, filepath string
}

func NewMD(fpath string, r io.Reader, title string, coEdit bool) (*MD, error) {
	m := &MD{
		filepath: fpath,
	}
	if err := m.loadContentFromReader(r, false); err != nil {
		return nil, xerrors.Errorf("failed to NewMDForPublish: %w", err)
	}
	if title != "" {
		m.FrontMatter.Title = title
	}
	if m.FrontMatter.Title == "" {
		return nil, xerrors.New("title required")
	}
	if !coEdit && m.FrontMatter.Author == "" {
		m.FrontMatter.Author = "dummy"
	}
	return m, nil
}

type Meta struct {
	Title  string   `yaml:"title"`
	Author string   `yaml:"author,omitempty"`
	Groups []string `yaml:"groups,flow"`
	Folder string   `yaml:"folder,omitempty"`
}

func (me *Meta) coediting() bool {
	return me.Author == ""
}

func (m *MD) fullContent() string {
	fm, _ := yaml.Marshal(m.FrontMatter)

	c := strings.Join([]string{"---", string(fm) + "---", "", m.Content}, "\n")
	if !strings.HasSuffix(c, "\n") {
		// fill newline for suppressing warning of "No newline at end of file"
		c += "\n"
	}
	return c
}

func (m *MD) save() error {
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
	log.Printf("saved to %q", m.filepath)
	return nil
}

func LoadMD(fpath string) (*MD, error) {
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
	m := &MD{
		ID:        newID(idTypeBlog, num),
		UpdatedAt: fi.ModTime(),
		filepath:  fpath,
	}
	if err := m.loadContentFromReader(f, true); err != nil {
		return nil, xerrors.Errorf("failed to LoadMD: %w", err)
	}
	return m, nil
}

func (m *MD) loadContentFromReader(r io.Reader, forceFrontmatter bool) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return xerrors.Errorf("failed to load md: %w", err)
	}
	if m.FrontMatter == nil {
		m.FrontMatter = &Meta{}
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

func (m *MD) toNote() *Note {
	groups := make([]*Group, len(m.FrontMatter.Groups))
	for i, g := range m.FrontMatter.Groups {
		groups[i] = &Group{Name: g}
	}
	return &Note{
		ID:        m.ID,
		Title:     m.FrontMatter.Title,
		Content:   m.Content,
		CoEditing: m.FrontMatter.coediting(),
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

func (ki *Kibela) PushMD(m *MD) error {
	n := m.toNote()
	if err := ki.pushNote(n); err != nil {
		return xerrors.Errorf("failed to pushMD: %w", err)
	}
	return os.Chtimes(m.filepath, n.UpdatedAt.Time, n.UpdatedAt.Time)
}

func (ki *Kibela) PublishMD(m *MD, save bool) error {
	groupIDs := make([]string, len(m.FrontMatter.Groups))
	for i, g := range m.FrontMatter.Groups {
		id, err := ki.fetchGroupID(g)
		if err != nil {
			return xerrors.Errorf("failed to publishMD: %w", err)
		}
		groupIDs[i] = string(id)
	}
	sort.Strings(groupIDs)
	data, err := ki.cli.Do(&client.Payload{
		Query: createNoteMutation,
		Variables: struct {
			Input *noteInput `json:"input"`
		}{
			Input: &noteInput{
				Title:     m.FrontMatter.Title,
				Content:   m.Content,
				Folder:    m.FrontMatter.Folder,
				CoEditing: m.FrontMatter.coediting(),
				GroupIDs:  groupIDs,
			},
		},
	})
	if err != nil {
		return xerrors.Errorf("failed to publishNote while accessing remote: %w", err)
	}
	var res struct {
		CreateNote struct {
			Note *Note `json:"note"`
		} `json:"createNote"`
	}
	if err := json.Unmarshal(data, &res); err != nil {
		return xerrors.Errorf("failed to ki.publishNote while unmarshaling response: %w", err)
	}
	if res.CreateNote.Note == nil {
		return xerrors.New("failed to publish to kibela on any reason. null createNote was returned")
	}
	n := res.CreateNote.Note
	n.CoEditing = m.FrontMatter.coediting()
	log.Printf("published %s", ki.noteURL(n))
	if !save {
		return nil
	}
	groups := make([]string, len(n.Groups))
	for i, g := range n.Groups {
		groups[i] = g.Name
	}
	m.FrontMatter.Groups = groups
	m.ID = n.ID
	m.UpdatedAt = n.UpdatedAt.Time
	if !n.CoEditing {
		m.FrontMatter.Author = n.Author.Account
	}
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
