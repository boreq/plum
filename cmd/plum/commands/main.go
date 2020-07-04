package commands

import (
	"github.com/boreq/guinea"
	"github.com/boreq/plum/logging"
)

var log = logging.New("main/commands")

var MainCmd = guinea.Command{
	Run: runMain,
	Subcommands: map[string]*guinea.Command{
		"run": &runCmd,
		"default_config": &defaultConfigCmd,
	},
	ShortDescription: "a real-time access log analyser",
	Description: `
Plum analyses web server access logs in real time and allows the user to access
the produced statistics using a web dashboard.
`,
}

func runMain(c guinea.Context) error {
	return guinea.ErrInvalidParms
}
