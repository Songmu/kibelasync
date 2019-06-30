package kibela

import (
	"encoding/json"

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
	Title     string `json:"title"`
	Content   string `json:"content"`
	CoEditing bool   `json:"coediting"`
	Folder    string `json:"folderName"`
	Groups    []struct {
		ID   `json:"id"`
		Name string `json:"name"`
	}
	Author struct {
		Account string `json:"account"`
	}
	UpdatedAt Time `json:"updatedAt"`
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
func (cli *client) getNotesCount() (int, error) {
	gResp, err := cli.Do(&payload{Query: totalCountQuery})
	if err != nil {
		return 0, xerrors.Errorf("failed to cli.getNotesCount: %w", err)
	}
	var res struct {
		Notes struct {
			TotalCount int `json:"totalCount"`
		} `json:"notes"`
	}
	if err := json.Unmarshal(gResp.Data, &res); err != nil {
		return 0, xerrors.Errorf("failed to cli.getNotesCount: %w", err)
	}
	return res.Notes.TotalCount, nil
}

// OK
func (cli *client) listNoteIDs() ([]*note, error) {
	num, err := cli.getNotesCount()
	if err != nil {
		return nil, xerrors.Errorf("failed to cli.listNodeIDs: %w", err)
	}
	gResp, err := cli.Do(&payload{Query: listNoteQuery(num)})
	if err != nil {
		return nil, xerrors.Errorf("failed to cli.getGroups: %w", err)
	}
	var res struct {
		Notes struct {
			Nodes []*note `json:"nodes"`
		} `json:"notes"`
	}
	if err := json.Unmarshal(gResp.Data, &res); err != nil {
		return nil, xerrors.Errorf("failed to cli.getNotesCount: %w", err)
	}
	return res.Notes.Nodes, nil
}
