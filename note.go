package kibela

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Songmu/kibela/client"
	"golang.org/x/xerrors"
)

/*
   {
     "id": "QmxvZy8zNjY",
     "title": "APIテストpublic",
     "content": "コンテント!\nコンテント",
     "coediting": true,
     "folderName": "testtop/testsub1",
     "groups": [
       {
         "name": "Home",
         "id": "R3JvdXAvMQ"
       }
     ],
     "author": {
       "account": "Songmu"
     },
     "createdAt": "2019-06-23T16:54:09.447Z",
     "publishedAt": "2019-06-23T16:54:09.444Z",
     "contentUpdatedAt": "2019-06-23T16:54:09.445Z",
     "updatedAt": "2019-06-23T17:22:38.496Z"
   },
*/
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
	gps := make([]string, len(n.Groups))
	for i, g := range n.Groups {
		gps[i] = g.Name
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
			Groups:    gps,
			Author:    n.Author.Account,
		},
	}
}

/*
{
  "data": {
    "notes": {
      "totalCount": 353
    }
  }
}
*/
// OK
func (ki *kibela) getNotesCount() (int, error) {
	gResp, err := ki.cli.Do(&client.Payload{Query: totalCountQuery})
	if err != nil {
		return 0, xerrors.Errorf("failed to ki.getNotesCount: %w", err)
	}
	var res struct {
		Notes struct {
			TotalCount int `json:"totalCount"`
		} `json:"notes"`
	}
	if err := json.Unmarshal(gResp, &res); err != nil {
		return 0, xerrors.Errorf("failed to ki.getNotesCount: %w", err)
	}
	return res.Notes.TotalCount, nil
}

// OK
func (ki *kibela) listNoteIDs() ([]*note, error) {
	num, err := ki.getNotesCount()
	if err != nil {
		return nil, xerrors.Errorf("failed to ki.listNodeIDs: %w", err)
	}
	gResp, err := ki.cli.Do(&client.Payload{Query: listNoteQuery(num)})
	if err != nil {
		return nil, xerrors.Errorf("failed to ki.getGroups: %w", err)
	}
	var res struct {
		Notes struct {
			Nodes []*note `json:"nodes"`
		} `json:"notes"`
	}
	if err := json.Unmarshal(gResp, &res); err != nil {
		return nil, xerrors.Errorf("failed to ki.getNotesCount: %w", err)
	}
	return res.Notes.Nodes, nil
}

// OK
func (ki *kibela) getNote(id ID) (*note, error) {
	gResp, err := ki.cli.Do(&client.Payload{Query: getNoteQuery(id)})
	if err != nil {
		return nil, xerrors.Errorf("failed to ki.getNote: %w", err)
	}
	var res struct {
		Note *note `json:"note"`
	}
	if err := json.Unmarshal(gResp, &res); err != nil {
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
	updatePayload := struct {
		ID       ID         `json:"id"`
		BaseNote *noteInput `json:"baseNote"`
		NewNote  *noteInput `json:"newNote"`
	}{
		ID:       n.ID,
		BaseNote: remoteNote.toNoteInput(),
		NewNote:  n.toNoteInput(),
	}
	data, err := ki.cli.Do(&client.Payload{
		Query:     updateNoteMutation,
		Variables: updatePayload,
	})
	if err != nil {
		return xerrors.Errorf("failed to pushNote while accessing remote: %w", err)
	}
	var res struct {
		Note struct {
			UpdatedAt Time `json:"updatedAt"`
		} `json:"note"`
	}
	if err := json.Unmarshal(data, &res); err != nil {
		return xerrors.Errorf("failed to ki.pushNote while unmarshaling response: %w", err)
	}
	n.UpdatedAt = res.Note.UpdatedAt
	return nil
}

func (n *note) toNoteInput() *noteInput {
	groupIDs := make([]ID, len(n.Groups))
	for i, g := range n.Groups {
		groupIDs[i] = g.ID
	}
	return &noteInput{
		Title:     n.Title,
		Content:   n.Content,
		GroupIDs:  groupIDs,
		Folder:    n.Folder,
		CoEditing: n.CoEditing,
	}
}
