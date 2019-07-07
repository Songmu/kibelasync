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
	)
	if err := fs.Parse(argv); err != nil {
		return err
	}
	mdFile := fs.Arg(0)
	ki, err := kibela.New(version)
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

	m, err := kibela.NewMD(mdFile, r, *title, *coEdit)
	if err != nil {
		return err
	}
	return ki.PublishMD(m, *save)
}
