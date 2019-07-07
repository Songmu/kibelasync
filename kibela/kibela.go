package kibela

import (
	"fmt"
	"os"
	"sync"

	"github.com/Songmu/kibelasync/client"
	"golang.org/x/xerrors"
)

const (
	envKibelaDIR   = "KIBELA_DIR"
	envKibelaTEAM  = "KIBELA_TEAM"
	envKibelaTOKEN = "KIBELA_TOKEN"
)

var defaultDir = "notes"

func init() {
	d := os.Getenv(envKibelaDIR)
	if d != "" {
		defaultDir = d
	}
}

type Kibela struct {
	cli *client.Client

	team string

	groups     map[string]ID
	groupsErr  error
	groupsOnce sync.Once
}

func New(ver string) (*Kibela, error) {
	token := os.Getenv(envKibelaTOKEN)
	if token == "" {
		return nil, fmt.Errorf("set token by KIBELA_TOKEN env value")
	}
	team := os.Getenv(envKibelaTEAM)
	if team == "" {
		return nil, fmt.Errorf("set team name by KIBELA_TEAM env value")
	}
	cli, err := client.New(ver, team, token)
	if err != nil {
		return nil, xerrors.Errorf("failed to kibela.New: %w", err)
	}
	return &Kibela{
		cli:  cli,
		team: team,
	}, nil
}

func (ki *Kibela) fetchGroups() (map[string]ID, error) {
	ki.groupsOnce.Do(func() {
		if ki.groups != nil {
			return
		}
		groups, err := ki.getGroups()
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

func (ki *Kibela) fetchGroupID(name string) (ID, error) {
	groups, err := ki.fetchGroups()
	if err != nil {
		return "", xerrors.Errorf("failed to fetchGroupID while setGroupID: %w", err)
	}
	id, ok := groups[name]
	if !ok {
		return "", fmt.Errorf("group %q doesn't exists", name)
	}
	return id, nil
}
