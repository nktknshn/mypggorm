package main

import (
	"github.com/nktknshn/mypggorm/cli"
	"github.com/spf13/cobra"
)

var (
	flagDBConfig string = "database-config.yaml"
)

func init() {
	cfg := cli.NewDatabaseConfigProvider(func() string { return flagDBConfig })
	commandDB := cli.DatabaseCommand(cfg)
	commandDB.PersistentFlags().StringVarP(&flagDBConfig, "config", "c", flagDBConfig, "database config file")
	cmd.AddCommand(commandDB)
}

var cmd = &cobra.Command{
	Use:  "cli",
	Args: cobra.MinimumNArgs(1),
}

func main() {
	err := cmd.Execute()

	if err != nil {
		panic(err)
	}
}
