package kibela

import (
	"strings"
	"time"

	"github.com/ghodss/yaml"
)

type md struct {
	FrontMatter *meta
	Content     string
}

type meta struct {
	ID        ID        `json:"-"`
	Title     string    `json:"title"`
	CoEditing bool      `json:"coediting"`
	Folder    string    `json:"folder"`
	Groups    []string  `json:"groups"`
	Author    string    `json:"author"`
	UpdatedAt time.Time `json:"-"`
}

func (m *md) fullContent() string {
	fm, _ := yaml.Marshal(m.FrontMatter)

	return strings.Join([]string{"---", string(fm) + "---", "", m.Content}, "\n")
}
