package kibela

import (
	"fmt"

	"github.com/Songmu/kibela/client"
	"golang.org/x/xerrors"
)

type kibela struct {
	cli    *client.Client
	groups map[string]ID
}

func newKibela() (*kibela, error) {
	cli, err := client.New(version)
	if err != nil {
		return nil, xerrors.Errorf("failed to newKibela: %w", err)
	}
	return &kibela{cli: cli}, nil
}

func (ki *kibela) setGroups() error {
	if ki.groups != nil {
		return nil
	}
	// XXX race
	groups, err := ki.getGroups()
	if err != nil {
		return xerrors.Errorf("failed to ki.setGroups: %w", err)
	}
	groupMap := make(map[string]ID, len(groups))
	for _, g := range groups {
		groupMap[g.Name] = g.ID
	}
	ki.groups = groupMap
	return nil
}

func (ki *kibela) fetchGroupID(name string) (ID, error) {
	if err := ki.setGroups(); err != nil {
		return "", xerrors.Errorf("failed to fetchGroupID while setGroupID: %w", err)
	}
	id, ok := ki.groups[name]
	if !ok {
		return "", fmt.Errorf("group %q doesn't exists", name)
	}
	return id, nil
}
