package kibelasync

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/Songmu/kibelasync/client"
	"golang.org/x/xerrors"
)

type note struct {
	ID        `json:"id"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	CoEditing bool     `json:"coediting"`
	Folder    string   `json:"folderName"`
	Groups    []*group `json:"groups"`
	Author    struct {
		Account string `json:"account"`
	}
	UpdatedAt Time `json:"updatedAt"`
}

func (n *note) toMD(dir string) *md {
	groups := make([]string, len(n.Groups))
	for i, g := range n.Groups {
		groups[i] = g.Name
	}
	return &md{
		ID:        n.ID,
		Content:   n.Content,
		UpdatedAt: n.UpdatedAt.Time,
		dir:       dir,
		FrontMatter: &meta{
			Title:     n.Title,
			CoEditing: n.CoEditing,
			Folder:    n.Folder,
			Groups:    groups,
			Author:    n.Author.Account,
		},
	}
}

func (ki *kibela) getNotesCount() (int, error) {
	data, err := ki.cli.Do(&client.Payload{Query: totalCountQuery})
	if err != nil {
		return 0, xerrors.Errorf("failed to ki.getNotesCount: %w", err)
	}
	var res struct {
		Notes struct {
			TotalCount int `json:"totalCount"`
		} `json:"notes"`
	}
	if err := json.Unmarshal(data, &res); err != nil {
		return 0, xerrors.Errorf("failed to ki.getNotesCount: %w", err)
	}
	return res.Notes.TotalCount, nil
}

const (
	// max query cost per request is 10,000
	// so adjust limit size to not exceed the limit
	// 100 (base) = 2 (id, updatedAt) * 4900 = 9900
	bundleLimit = 4900
	// 100 (base) = 3 (id, updatedAt, cursor) * 3200 = 9700
	pageLimit = 3200
)

// OK
func (ki *kibela) listNoteIDs() ([]*note, error) {
	num, err := ki.getNotesCount()
	if err != nil {
		return nil, xerrors.Errorf("failed to ki.listNodeIDs: %w", err)
	}
	if num > bundleLimit {
		nextCursor := ""
		rest := num
		notes := make([]*note, 0, num)
		for rest > 0 {
			take := pageLimit
			if take > rest {
				take = rest
			}
			rest = rest - take
			data, err := ki.cli.Do(&client.Payload{Query: listNotePaginateQuery(take, nextCursor)})
			if err != nil {
				return nil, xerrors.Errorf("failed to ki.getGroups: %w", err)
			}
			var res struct {
				Notes struct {
					Edges []struct {
						Node   *note  `json:"node"`
						Cursor string `json:"cursor"`
					} `json:"edges"`
				} `json:"notes"`
			}
			if err := json.Unmarshal(data, &res); err != nil {
				return nil, xerrors.Errorf("failed to ki.listNoteIDs: %w", err)
			}
			if len(res.Notes.Edges) > 0 {
				nextCursor = res.Notes.Edges[len(res.Notes.Edges)-1].Cursor
			}
			for _, n := range res.Notes.Edges {
				notes = append(notes, n.Node)
			}
		}
		return notes, nil
	}
	data, err := ki.cli.Do(&client.Payload{Query: listNoteQuery(num)})
	if err != nil {
		return nil, xerrors.Errorf("failed to ki.listNoteIDs: %w", err)
	}
	var res struct {
		Notes struct {
			Nodes []*note `json:"nodes"`
		} `json:"notes"`
	}
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, xerrors.Errorf("failed to ki.getNotesCount: %w", err)
	}
	return res.Notes.Nodes, nil
}

// OK
func (ki *kibela) getNote(id ID) (*note, error) {
	data, err := ki.cli.Do(&client.Payload{Query: getNoteQuery(id)})
	if err != nil {
		return nil, xerrors.Errorf("failed to ki.getNote: %w", err)
	}
	var res struct {
		Note *note `json:"note"`
	}
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, xerrors.Errorf("failed to ki.getNote: %w", err)
	}
	res.Note.ID = id
	return res.Note, nil
}

func (ki *kibela) pullNotes(dir string) error {
	notes, err := ki.listNoteIDs()
	if err != nil {
		return xerrors.Errorf("failed to pullNotes: %w", err)
	}
	for _, n := range notes {
		localT := time.Time{}
		idNum, err := n.ID.Number()
		if err != nil {
			return xerrors.Errorf("failed to pullNotes: %w", err)
		}
		mdFilePath := filepath.Join(dir, fmt.Sprintf("%d.md", idNum))
		_, err = os.Stat(mdFilePath)
		if err == nil {
			localMD, err := loadMD(mdFilePath)
			if err != nil {
				return xerrors.Errorf("failed to pullNotes: %w", err)
			}
			localT = localMD.UpdatedAt
		}
		if n.UpdatedAt.After(localT) {
			allNote, err := ki.getNote(n.ID)
			if err != nil {
				return xerrors.Errorf("failed to pullNotes: %w", err)
			}
			if err := allNote.toMD(dir).save(); err != nil {
				return xerrors.Errorf("failed to pullNotes: %w", err)
			}
		}
	}
	return nil
}

const pullBundleLimit = 50

func (ki *kibela) pullFullNotes(dir string) error {
	num, err := ki.getNotesCount()
	if err != nil {
		return xerrors.Errorf("failed to ki.pullFullNotes: %w", err)
	}
	nextCursor := ""
	rest := num
	for rest > 0 {
		take := pullBundleLimit
		if take > rest {
			take = rest
		}
		rest = rest - take
		data, err := ki.cli.Do(&client.Payload{Query: listFullNotePaginateQuery(take, nextCursor)})
		if err != nil {
			return xerrors.Errorf("failed to ki.pullFullNotes: %w", err)
		}
		var res struct {
			Notes struct {
				Edges []struct {
					Node   *note  `json:"node"`
					Cursor string `json:"cursor"`
				} `json:"edges"`
			} `json:"notes"`
		}
		if err := json.Unmarshal(data, &res); err != nil {
			return xerrors.Errorf("failed to ki.pullFullNotes: %w", err)
		}
		if len(res.Notes.Edges) > 0 {
			nextCursor = res.Notes.Edges[len(res.Notes.Edges)-1].Cursor
		}
		for _, e := range res.Notes.Edges {
			if err := e.Node.toMD(dir).save(); err != nil {
				return xerrors.Errorf("failed to pullFullNotes while saving md: %w", err)
			}
		}
	}
	return nil
}

func (ki *kibela) pushNote(n *note) error {
	remoteNote, err := ki.getNote(n.ID)
	if err != nil {
		return xerrors.Errorf("failed to pushNote: %w", err)
	}
	groupMap := make(map[string]ID)
	for _, g := range remoteNote.Groups {
		groupMap[g.Name] = g.ID
	}

	// fill group IDs
	for _, g := range n.Groups {
		if string(g.ID) == "" {
			id, ok := groupMap[g.Name]
			if ok {
				g.ID = id
			} else {
				id, err := ki.fetchGroupID(g.Name)
				if err != nil {
					return xerrors.Errorf("failed to fetch group name: %q, %w", g.Name, err)
				}
				g.ID = id
			}
		}
	}
	baseNote := remoteNote.toNoteInput()
	newNote := n.toNoteInput()
	if reflect.DeepEqual(*baseNote, *newNote) {
		// no update defferences
		n.UpdatedAt = remoteNote.UpdatedAt
		return nil
	}
	data, err := ki.cli.Do(&client.Payload{
		Query: updateNoteMutation,
		Variables: struct {
			ID       ID         `json:"id"`
			BaseNote *noteInput `json:"baseNote"`
			NewNote  *noteInput `json:"newNote"`
		}{
			ID:       n.ID,
			BaseNote: baseNote,
			NewNote:  newNote,
		},
	})
	if err != nil {
		return xerrors.Errorf("failed to pushNote while accessing remote: %w", err)
	}
	var res struct {
		UpdateNote struct {
			Note *note `json:"note"`
		} `json:"updateNote"`
	}
	if err := json.Unmarshal(data, &res); err != nil {
		return xerrors.Errorf("failed to ki.pushNote while unmarshaling response: %w", err)
	}
	if res.UpdateNote.Note == nil {
		return xerrors.New("failed to update kibela on any reason. null updateNote was returned")
	}
	log.Printf("updated %s", ki.noteURL(n))
	n.UpdatedAt = res.UpdateNote.Note.UpdatedAt
	return nil
}

func (n *note) toNoteInput() *noteInput {
	groupIDs := make([]string, len(n.Groups))
	for i, g := range n.Groups {
		groupIDs[i] = string(g.ID)
	}
	sort.Strings(groupIDs)
	return &noteInput{
		Title:     n.Title,
		Content:   strings.TrimSpace(n.Content) + "\n",
		GroupIDs:  groupIDs,
		Folder:    n.Folder,
		CoEditing: n.CoEditing,
	}
}

func (ki *kibela) noteURL(n *note) string {
	basePath := "notes"
	if !n.CoEditing {
		basePath = "@" + n.Author.Account
	}
	num, _ := n.ID.Number()
	return fmt.Sprintf("https://%s.kibe.la/%s/%d", ki.team, basePath, num)
}
