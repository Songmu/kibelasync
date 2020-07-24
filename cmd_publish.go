package kibelasync

import (
	"context"
	"flag"
	"io"
	"os"

	"github.com/Songmu/kibelasync/kibela"
)

type cmdPublish struct{}

func (cp *cmdPublish) name() string {
	return "publish"
}

func (cp *cmdPublish) description() string {
	return "publish new markdown"
}

func (cp *cmdPublish) run(ctx context.Context, argv []string, outStream io.Writer, errStream io.Writer) error {
	fs := flag.NewFlagSet("kibelasync publish", flag.ContinueOnError)
	fs.SetOutput(errStream)
	var (
		title  = fs.String("title", "", "title of the note")
		save   = fs.Bool("save", false, "save file after published the note")
		coEdit = fs.Bool("co-edit", false, "co-editing on")
		dir    = fs.String("dir", "notes", "sync directory")
	)
	if err := fs.Parse(argv); err != nil {
		return err
	}
	mdFile := fs.Arg(0)
	ki, err := kibela.New(version)
	if err != nil {
		return err
	}

	var r io.Reader = os.Stdin
	if mdFile != "" {
		var err error
		f, err := os.Open(mdFile)
		if err != nil {
			return err
		}
		defer f.Close()
		r = f
	}

	m, err := kibela.NewMD(mdFile, r, *title, *coEdit, *dir)
	if err != nil {
		return err
	}
	return ki.PublishMD(ctx, m, *save)
}
