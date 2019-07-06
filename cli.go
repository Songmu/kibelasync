package kibela

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"

	"golang.org/x/xerrors"
)

const cmdName = "kibela"

// Run the kibela
func Run(argv []string, outStream, errStream io.Writer) error {
	log.SetOutput(errStream)
	log.SetPrefix(fmt.Sprintf("[%s] ", cmdName))
	nameAndVer := fmt.Sprintf("%s (v%s rev:%s)", cmdName, version, revision)
	fs := flag.NewFlagSet(nameAndVer, flag.ContinueOnError)
	fs.SetOutput(errStream)
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage of %s:\n", nameAndVer)
		fs.PrintDefaults()
	}

	ver := fs.Bool("version", false, "display version")
	if err := fs.Parse(argv); err != nil {
		return err
	}
	if *ver {
		return printVersion(outStream)
	}

	argv = fs.Args()
	if len(argv) < 1 {
		return xerrors.New("no subcommand specified")
	}
	rnr, ok := dispatch[argv[0]]
	if !ok {
		return xerrors.Errorf("unknown subcommand: %s", argv[0])
	}
	return rnr.run(context.Background(), argv[1:], outStream, errStream)
}

func printVersion(out io.Writer) error {
	_, err := fmt.Fprintf(out, "%s v%s (rev:%s)\n", cmdName, version, revision)
	return err
}

var dispatch = map[string]runner{
	"publish": &cmdPublish{},
	"pull":    &cmdPull{},
	"push":    &cmdPush{},
}

type runner interface {
	description() string
	run(context.Context, []string, io.Writer, io.Writer) error
}
