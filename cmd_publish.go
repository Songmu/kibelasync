package kibela

import (
	"flag"
	"io"
	"os"

	"golang.org/x/xerrors"
)

type cmdPublish struct {
}

func (cp *cmdPublish) run(argv []string, outStream io.Writer, errStream io.Writer) error {
	fs := flag.NewFlagSet("kibela pull", flag.ContinueOnError)
	fs.SetOutput(errStream)
	var (
		title = fs.String("title", "", "title of the note")
		save  = fs.Bool("save", true, "save file after published the note")
	)
	if err := fs.Parse(argv); err != nil {
		return err
	}
	mdFile := fs.Arg(0)
	ki, err := newKibela()
	if err != nil {
		return err
	}

	var r io.ReadCloser = os.Stdin
	if mdFile != "" {
		var err error
		if r, err = os.Open(mdFile); err != nil {
			return err
		}
	}
	defer r.Close()

	m := &md{
		filepath: mdFile,
	}
	m.loadContentFromReader(r, false)
	if *title != "" {
		if m.FrontMatter == nil {
			m.FrontMatter = &meta{}
		}
		m.FrontMatter.Title = *title
	}
	if m.FrontMatter == nil || m.FrontMatter.Title == "" {
		// XXX detect title from markdown?
		return xerrors.New("title required")
	}
	return ki.publishMD(m, *save)
}
