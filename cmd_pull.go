package kibelasync

import (
	"context"
	"flag"
	"io"
)

type cmdPull struct{}

func (cp *cmdPull) name() string {
	return "pull"
}

func (cp *cmdPull) description() string {
	return "sync all markdowns"
}

func (cp *cmdPull) run(ctx context.Context, argv []string, outStream io.Writer, errStream io.Writer) error {
	fs := flag.NewFlagSet("kibelasync pull", flag.ContinueOnError)
	var full = fs.Bool("full", false, "pull every markdowns")
	fs.SetOutput(errStream)

	if err := fs.Parse(argv); err != nil {
		return err
	}
	dir := fs.Arg(0)
	if dir == "" {
		dir = "notes"
	}

	ki, err := newKibela()
	if err != nil {
		return err
	}
	if *full {
		return ki.pullFullNotes(dir)
	}
	return ki.pullNotes(dir)
}
