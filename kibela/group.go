package kibela

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Songmu/kibelasync/client"
	"golang.org/x/xerrors"
)

// Group represents groups of Kibela
type Group struct {
	ID   `json:"id"`
	Name string `json:"name"`
}

func (ki *Kibela) getGroupCount(ctx context.Context) (int, error) {
	data, err := ki.cli.Do(ctx, &client.Payload{Query: totalGroupCountQuery})
	if err != nil {
		return 0, xerrors.Errorf("failed to ki.getGroupCount: %w", err)
	}
	var res struct {
		Groups struct {
			TotalCount int `json:"totalCount"`
		} `json:"groups"`
	}
	if err := json.Unmarshal(data, &res); err != nil {
		return 0, xerrors.Errorf("failed to ki.getNotesCount: %w", err)
	}
	return res.Groups.TotalCount, nil
}

func (ki *Kibela) getGroups(ctx context.Context) ([]*Group, error) {
	num, err := ki.getGroupCount(ctx)
	if err != nil {
		return nil, xerrors.Errorf("failed to getGroups: %w", err)
	}
	data, err := ki.cli.Do(ctx, &client.Payload{Query: listGroupQuery(num)})
	if err != nil {
		return nil, xerrors.Errorf("failed to ki.getGroups: %w", err)
	}
	var res struct {
		Groups struct {
			Nodes []*Group `json:"nodes"`
		} `json:"groups"`
	}
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, xerrors.Errorf("failed to ki.getNotesCount: %w", err)
	}
	return res.Groups.Nodes, nil
}

func (ki *Kibela) fetchGroups(ctx context.Context) (map[string]ID, error) {
	ki.groupsOnce.Do(func() {
		if ki.groups != nil {
			return
		}
		groups, err := ki.getGroups(ctx)
		if err != nil {
			ki.groupsErr = xerrors.Errorf("failed to ki.setGroups: %w", err)
			return
		}
		groupMap := make(map[string]ID, len(groups))
		for _, g := range groups {
			groupMap[g.Name] = g.ID
		}
		ki.groups = groupMap
	})
	return ki.groups, ki.groupsErr
}

func (ki *Kibela) fetchGroupID(ctx context.Context, name string) (ID, error) {
	groups, err := ki.fetchGroups(ctx)
	if err != nil {
		return "", xerrors.Errorf("failed to fetchGroupID while setGroupID: %w", err)
	}
	id, ok := groups[name]
	if !ok {
		return "", fmt.Errorf("group %q doesn't exists", name)
	}
	return id, nil
}
