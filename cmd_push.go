package kibela

import (
	"flag"
	"io"

	"golang.org/x/xerrors"
)

type cmdPush struct {
}

func (cp *cmdPush) run(argv []string, outStream io.Writer, errStream io.Writer) error {
	fs := flag.NewFlagSet("kibela push", flag.ContinueOnError)
	fs.SetOutput(errStream)

	if err := fs.Parse(argv); err != nil {
		return err
	}

	ki, err := newKibela()
	if err != nil {
		return err
	}
	if fs.NArg() < 1 {
		return xerrors.New("usage: kibela pull [md files]")
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
