package kibela

import (
	"flag"
	"io"
)

type cmdPull struct {
}

func (cp *cmdPull) run(argv []string, outStream io.Writer, errStream io.Writer) error {
	fs := flag.NewFlagSet("kibela pull", flag.ContinueOnError)
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
	return ki.pullNotes(dir)
}
