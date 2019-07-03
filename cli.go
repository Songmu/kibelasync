package kibela

import (
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

	argv = fs.Args()
	if len(argv) < 1 {
		return xerrors.New("no subcommand specified")
	}
	rnr, ok := dispatch[argv[0]]
	if !ok {
		return xerrors.Errorf("unknown subcommand: %s", argv[0])
	}
	return rnr.run(argv[1:], outStream, errStream)
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
	run([]string, io.Writer, io.Writer) error
}
