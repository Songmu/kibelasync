package kibela

import (
	"encoding/json"

	"golang.org/x/xerrors"
)

type group struct {
	ID   `json:"id"`
	Name string `json:"name"`
}

/*
{
  "data": {
    "groups": {
      "totalCount": 353
    }
  }
}
*/
// OK
func (cli *client) getGroupCount() (int, error) {
	gResp, err := cli.Do(&payload{Query: totalGroupCountQuery})
	if err != nil {
		return 0, xerrors.Errorf("failed to cli.getGroupCount: %w", err)
	}
	var res struct {
		Groups struct {
			TotalCount int `json:"totalCount"`
		} `json:"groups"`
	}
	if err := json.Unmarshal(gResp.Data, &res); err != nil {
		return 0, xerrors.Errorf("failed to cli.getNotesCount: %w", err)
	}
	return res.Groups.TotalCount, nil
}

/*
{
  "data": {
    "groups": {
      "nodes": [
        {
          "id": "R3JvdXAvMQ",
          "name": "Home"
        },
        {
          "id": "R3JvdXAvMg",
          "name": "Test"
        }
      ]
    }
  }
}
*/
// OK
func (cli *client) getGroups() ([]*group, error) {
	num, err := cli.getGroupCount()
	if err != nil {
		return nil, xerrors.Errorf("failed to getGroups: %w", err)
	}
	gResp, err := cli.Do(&payload{Query: listGroupQuery(num)})
	if err != nil {
		return nil, xerrors.Errorf("failed to cli.getGroups: %w", err)
	}
	var res struct {
		Groups struct {
			Nodes []*group `json:"nodes"`
		} `json:"groups"`
	}
	if err := json.Unmarshal(gResp.Data, &res); err != nil {
		return nil, xerrors.Errorf("failed to cli.getNotesCount: %w", err)
	}
	return res.Groups.Nodes, nil
}
