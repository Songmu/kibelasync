package kibela

import (
	"flag"
	"fmt"
	"io"
	"log"

	"github.com/Songmu/kibela/client"
	"golang.org/x/xerrors"
)

const cmdName = "kibela"

// Run the kibela
func Run(argv []string, outStream, errStream io.Writer) error {
	log.SetOutput(errStream)
	fs := flag.NewFlagSet(
		fmt.Sprintf("%s (v%s rev:%s)", cmdName, version, revision), flag.ContinueOnError)
	fs.SetOutput(errStream)
	ver := fs.Bool("version", false, "display version")
	if err := fs.Parse(argv); err != nil {
		return err
	}
	if *ver {
		return printVersion(outStream)
	}
	return nil
}

func printVersion(out io.Writer) error {
	_, err := fmt.Fprintf(out, "%s v%s (rev:%s)\n", cmdName, version, revision)
	return err
}

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
