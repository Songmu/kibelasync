package kibela

import (
	"encoding/json"

	"github.com/Songmu/kibela/client"
	"golang.org/x/xerrors"
)

type group struct {
	ID   `json:"id"`
	Name string `json:"name"`
}

func (ki *kibela) getGroupCount() (int, error) {
	gResp, err := ki.cli.Do(&client.Payload{Query: totalGroupCountQuery})
	if err != nil {
		return 0, xerrors.Errorf("failed to ki.getGroupCount: %w", err)
	}
	var res struct {
		Groups struct {
			TotalCount int `json:"totalCount"`
		} `json:"groups"`
	}
	if err := json.Unmarshal(gResp, &res); err != nil {
		return 0, xerrors.Errorf("failed to ki.getNotesCount: %w", err)
	}
	return res.Groups.TotalCount, nil
}

func (ki *kibela) getGroups() ([]*group, error) {
	num, err := ki.getGroupCount()
	if err != nil {
		return nil, xerrors.Errorf("failed to getGroups: %w", err)
	}
	gResp, err := ki.cli.Do(&client.Payload{Query: listGroupQuery(num)})
	if err != nil {
		return nil, xerrors.Errorf("failed to ki.getGroups: %w", err)
	}
	var res struct {
		Groups struct {
			Nodes []*group `json:"nodes"`
		} `json:"groups"`
	}
	if err := json.Unmarshal(gResp, &res); err != nil {
		return nil, xerrors.Errorf("failed to ki.getNotesCount: %w", err)
	}
	return res.Groups.Nodes, nil
}
