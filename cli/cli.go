package cli

import (
	"fmt"

	"github.com/nktknshn/mypggorm"
	"github.com/spf13/cobra"
)

// type DatabaseSchema interface {
// 	Migrate(db *gorm.DB) error
// 	GetSchema(db *gorm.DB) (DatabaseSchema, error)
// }

type DatabaseConfig interface {
	GetRootConfig() (mypggorm.DatabaseConnectionConfig, error)
	GetUserConfig() (mypggorm.DatabaseConnectionConfig, error)
}

type DatabaseMethods interface {
	DatabaseConfig
	// DatabaseSchema
}

var (
	DbConfig DatabaseMethods
)

func DatabaseCommand(config DatabaseMethods) *cobra.Command {
	DbConfig = config
	return commandDatabase
}

var commandDatabase = &cobra.Command{
	Use:  "database",
	Args: cobra.MinimumNArgs(1),
}

func init() {

	commandDatabase.AddCommand(commandCreate)
	commandDatabase.AddCommand(commandCheck)
	commandDatabase.AddCommand(commandDropDatabase)
	commandDatabase.AddCommand(commandBackup)
	commandDatabase.AddCommand(commandReset)
	commandDatabase.AddCommand(commandRestore)

	// commandDatabase.AddCommand(commandMigrage)
	// commandDatabase.AddCommand(commandCreateTables)
	// commandDatabase.AddCommand(commandPrintSchema)
}

func runWithDBConfig(f func(cmd *cobra.Command, dbConfig DatabaseMethods, args []string) error) func(cmd *cobra.Command, args []string) error {

	return func(cmd *cobra.Command, args []string) error {
		if DbConfig == nil {
			return fmt.Errorf("dbConfig is nil")
		}

		return f(cmd, DbConfig, args)
	}
}

func RunWithConfigs(f func(cmd *cobra.Command, rootConfig mypggorm.DatabaseConnectionConfig, userConfig mypggorm.DatabaseConnectionConfig, args []string) error) func(cmd *cobra.Command, args []string) error {

	return func(cmd *cobra.Command, args []string) error {

		if DbConfig == nil {
			return fmt.Errorf("dbConfig is nil")
		}

		rootConfig, err := DbConfig.GetRootConfig()

		if err != nil {
			return fmt.Errorf("failed to get root config: %w", err)
		}

		userConfig, err := DbConfig.GetUserConfig()

		if err != nil {
			return fmt.Errorf("failed to get user config: %w", err)
		}

		return f(cmd, rootConfig, userConfig, args)
	}
}
