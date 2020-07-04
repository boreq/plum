package commands

import (
	"encoding/json"
	"fmt"
	"github.com/boreq/guinea"
	"github.com/boreq/plum/config"
	"github.com/pkg/errors"
)

var defaultConfigCmd = guinea.Command{
	Run:              runDefaultConfig,
	ShortDescription: "prints default config to stdout",
}

func runDefaultConfig(c guinea.Context) error {
	conf := config.Default()

	j, err := json.MarshalIndent(conf, "", "    ")
	if err != nil {
		return errors.Wrap(err, "could not marshal the config")
	}

	fmt.Println(string(j))

	return nil
}
