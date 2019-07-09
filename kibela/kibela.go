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

	folders     map[string]ID
	foldersErr  error
	foldersOnce sync.Once
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
