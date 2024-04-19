package cli

import (
	"fmt"

	"github.com/go-faster/errors"
	"github.com/nktknshn/mypggorm"
	"github.com/spf13/cobra"
)

// type DatabaseSchema interface {
// 	Migrate(db *gorm.DB) error
// 	GetSchema(db *gorm.DB) (DatabaseSchema, error)
// }

type DatabaseRootConfigGetter interface {
	GetRootConfig() (mypggorm.DatabaseConnectionConfig, error)
}

type DatabaseUserConfigGetter interface {
	GetUserConfig() (mypggorm.DatabaseConnectionConfig, error)
}

type DatabaseConfigGetter interface {
	DatabaseRootConfigGetter
	DatabaseUserConfigGetter
}

var (
	DbConfig DatabaseConfigGetter
)

func DatabaseCommand(config DatabaseConfigGetter) *cobra.Command {
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

func runWithDBConfig(f func(cmd *cobra.Command, dbConfig DatabaseConfigGetter, args []string) error) func(cmd *cobra.Command, args []string) error {

	return func(cmd *cobra.Command, args []string) error {
		if DbConfig == nil {
			return fmt.Errorf("dbConfig is nil")
		}

		return f(cmd, DbConfig, args)
	}
}

type RunWithConfigsFunc func(cmd *cobra.Command, rootConfig mypggorm.DatabaseConnectionConfig, userConfig mypggorm.DatabaseConnectionConfig, args []string) error

func RunWithConfigs(f RunWithConfigsFunc) func(cmd *cobra.Command, args []string) error {

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

// print database schema
var CmdPrintSchema = &cobra.Command{
	Use:  "print-schema",
	RunE: RunWithConfigs(runPrintSchema),
}

func runPrintSchema(cmd *cobra.Command, rootConfig mypggorm.DatabaseConnectionConfig, userConfig mypggorm.DatabaseConnectionConfig, args []string) error {

	conn, err := userConfig.Connect()
	if err != nil {
		return errors.Wrap(err, "connect to database")
	}
	schema, err := mypggorm.GetSchemaPublic(conn)
	if err != nil {
		return errors.Wrap(err, "get schema")
	}
	fmt.Println(schema.String())
	return nil
}

type Migrator interface{ Migrate() error }
type MigratorGetter func(rootConfig, userConfig mypggorm.DatabaseConnectionConfig) (Migrator, error)

var CmdMigrate = func(migrator MigratorGetter) *cobra.Command {
	return &cobra.Command{
		Use:  "migrate",
		RunE: RunWithConfigs(runMigrate(migrator)),
	}
}

func runMigrate(migrator MigratorGetter) RunWithConfigsFunc {
	return func(cmd *cobra.Command, rootConfig, userConfig mypggorm.DatabaseConnectionConfig, args []string) error {

		fmt.Println("Migrating database", userConfig.Dbname)
		conn, err := userConfig.Connect()
		if err != nil {
			return errors.Wrap(err, "connect to database")
		}
		m, err := migrator(rootConfig, userConfig)
		if err != nil {
			return errors.Wrap(err, "migrate database")
		}
		if err := m.Migrate(); err != nil {
			return errors.Wrap(err, "migrate database")
		}
		fmt.Println("Database migrated")
		schema, err := mypggorm.GetSchemaPublic(conn)
		if err != nil {
			return errors.Wrap(err, "get schema")
		}
		fmt.Println(schema.String())
		return nil
	}
}
