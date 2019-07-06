package kibelasync

import (
	"context"
	"flag"
	"io"

	"golang.org/x/xerrors"
)

type cmdPush struct{}

func (cp *cmdPush) name() string {
	return "push"
}

func (cp *cmdPush) description() string {
	return "push markdown"
}

func (cp *cmdPush) run(ctx context.Context, argv []string, outStream io.Writer, errStream io.Writer) error {
	fs := flag.NewFlagSet("kibelasync push", flag.ContinueOnError)
	fs.SetOutput(errStream)

	if err := fs.Parse(argv); err != nil {
		return err
	}

	ki, err := newKibela()
	if err != nil {
		return err
	}
	if fs.NArg() < 1 {
		return xerrors.New("usage: kibelasync pull [md files]")
	}
	for _, f := range fs.Args() {
		md, err := loadMD(f)
		if err != nil {
			return err
		}
		if err := ki.pushMD(md); err != nil {
			return err
		}
	}
	return nil
}
