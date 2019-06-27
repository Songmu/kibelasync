package kibela

import (
	"strings"

	"github.com/ghodss/yaml"
)

type md struct {
	FrontMatter *meta
	Content     string
}

type meta struct {
	ID        ID       `json:"-"`
	Title     string   `json:"title"`
	CoEditing bool     `json:"coediting"`
	Folder    string   `json:"folder"`
	Groups    []string `json:"groups"`
	Author    string   `json:"author"`
	UpdatedAt Time     `json:"updatedAt"` // XXX may be removed in future
}

func (m *md) fullContent() string {
	fm, _ := yaml.Marshal(m.FrontMatter)

	return strings.Join([]string{"---", string(fm) + "---", m.Content}, "\n")
}
