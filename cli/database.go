package cli

import (
	"fmt"

	"github.com/go-faster/errors"
	"github.com/nktknshn/mypggorm"
	"github.com/nktknshn/mypggorm/backup"
	"github.com/nktknshn/mypggorm/helpers"
	"github.com/nktknshn/mypggorm/postgres"
	"github.com/spf13/cobra"
)

var (
	flagDropUser     bool = false
	flagCreateUser   bool = false
	flagSkipDropping bool = false
)

func init() {
	commandDropDatabase.PersistentFlags().BoolVarP(&flagDropUser, "drop-user", "U", false, "drop user")
	commandDropDatabase.PersistentFlags().BoolVarP(&flagSkipDropping, "skip-dropping", "D", false, "skip dropping")
	commandDropDatabase.PersistentFlags().BoolVarP(&flagCreateUser, "create-user", "C", false, "create user")
}

var commandCreate = &cobra.Command{
	Use: "create",
	RunE: runWithConfigs(func(cmd *cobra.Command, rootConfig mypggorm.DatabaseConnectionConfig, userConfig mypggorm.DatabaseConnectionConfig, args []string) error {

		db := postgres.NewPostgresInstance(rootConfig)

		if err := db.CheckRootConnection(); err != nil {
			return errors.Wrap(err, "failed to check root connection")
		}

		fmt.Println("root connection is ok")

		if err := db.CreateUserAndDatabase(userConfig, true); err != nil {
			return errors.Wrap(err, "failed to create user and database")
		}

		if err := db.CheckConnection(userConfig); err != nil {
			return errors.Wrap(err, "failed to check user connection")
		}

		fmt.Println("user connection is ok")

		return nil
	}),
}

var commandCheck = &cobra.Command{
	Use: "check",
	RunE: runWithConfigs(func(cmd *cobra.Command, rootConfig mypggorm.DatabaseConnectionConfig, userConfig mypggorm.DatabaseConnectionConfig, args []string) error {

		db := postgres.NewPostgresInstance(rootConfig)

		if err := db.CheckRootConnection(); err != nil {
			return errors.Wrap(err, "failed to check root connection")
		}

		fmt.Println("root connection is ok")

		if err := db.CheckConnection(userConfig); err != nil {
			return errors.Wrap(err, "failed to check user connection")
		}

		fmt.Println("user connection is ok")

		return nil
	}),
}

var commandDropDatabase = &cobra.Command{
	Use: "drop",
	RunE: runWithConfigs(func(cmd *cobra.Command, rootConfig mypggorm.DatabaseConnectionConfig, userConfig mypggorm.DatabaseConnectionConfig, args []string) error {

		fmt.Println("DROPPING database", userConfig.Dbname)

		if flagDropUser {
			fmt.Println("User WILL be dropped too.")
		} else {
			fmt.Println("User will NOT be dropped.")
		}

		confirm, err := helpers.AskString(fmt.Sprintf("Are you sure you want to DROP database '%s'? type 'yes DROP it': ", userConfig.Dbname))

		if err != nil {
			return errors.Wrap(err, "failed to ask confirmation")
		}

		if confirm != "yes DROP it" {
			fmt.Println("Cancelled.")
			return nil
		}

		db := postgres.NewPostgresInstance(rootConfig)

		if err := db.DropDatabase(&userConfig, flagDropUser); err != nil {
			return errors.Wrap(err, "failed to drop database")
		}

		return nil
	}),
}

var commandRestore = &cobra.Command{
	Use:  "restore <psql_path> <backup_path>",
	Args: cobra.ExactArgs(2),
	RunE: runWithConfigs(func(cmd *cobra.Command, rootConfig mypggorm.DatabaseConnectionConfig, userConfig mypggorm.DatabaseConnectionConfig, args []string) error {
		psqlPath := args[0]
		backupPath := args[1]

		if !helpers.PathExists(psqlPath) {
			return fmt.Errorf("psql path %v does not exist", psqlPath)
		}

		if !helpers.PathExists(backupPath) {
			return fmt.Errorf("backup path %v does not exist", backupPath)
		}

		db := postgres.NewPostgresInstance(rootConfig)

		reply := helpers.AskStringMust("Are you sure you want to RESTORE database. type 'yes RESTORE it': \n")

		if reply != "yes RESTORE it" {
			fmt.Println("Restore canceled")
			return nil
		}

		fmt.Println("RESTORING database", userConfig.Dbname)

		if !flagSkipDropping {

			err := db.CheckConnection(userConfig)

			if err != nil {
				return errors.Wrap(err, "failed to check user connection")
			}

			fmt.Println("Dropping database", userConfig.Dbname)

			err = db.DropDatabase(&userConfig, flagDropUser)

			if err != nil {
				return errors.Wrap(err, "failed to drop database")
			}
		}

		if err := db.CreateUserAndDatabase(userConfig, flagCreateUser); err != nil {
			return errors.Wrap(err, "failed to create user and database")
		}

		if err := backup.NewPsqlRestore(psqlPath).Restore(userConfig, backupPath); err != nil {
			return errors.Wrap(err, "failed to restore database")
		}

		fmt.Println("Database restored.")
		return nil
	}),
}

var commandBackup = &cobra.Command{
	Use:  "backup <pg_dump_path> <backup_path>",
	Args: cobra.ExactArgs(2),
	RunE: runWithConfigs(func(cmd *cobra.Command, rootConfig mypggorm.DatabaseConnectionConfig, userConfig mypggorm.DatabaseConnectionConfig, args []string) error {

		pgdumpPath := args[0]
		backupPath := args[1]

		if !helpers.PathExists(pgdumpPath) {
			return fmt.Errorf("pg_dump path %v does not exist", pgdumpPath)
		}

		if helpers.PathExists(backupPath) {
			return fmt.Errorf("backup path %v already exists", backupPath)
		}

		db := postgres.NewPostgresInstance(rootConfig)

		fmt.Println("Backing up database", userConfig.Dbname)

		if err := db.CheckConnection(userConfig); err != nil {
			return err
		}

		if err := backup.NewPgDump(pgdumpPath).Backup(userConfig, backupPath); err != nil {
			return err
		}

		println("Backup created.")

		return nil
	}),
}

var commandReset = &cobra.Command{
	Use: "reset",
	RunE: runWithConfigs(func(cmd *cobra.Command, rootConfig mypggorm.DatabaseConnectionConfig, userConfig mypggorm.DatabaseConnectionConfig, args []string) error {

		db := postgres.NewPostgresInstance(rootConfig)
		if err := db.CheckRootConnection(); err != nil {
			return errors.Wrap(err, "failed to check root connection")
		}

		println("Root connection is ok.")

		if err := db.DropDatabase(&userConfig, true); err != nil {
			return err
		}

		println("Database and user dropped.")

		if err := db.CreateUserAndDatabase(userConfig, true); err != nil {
			return errors.Wrap(err, "failed to create user and database")
		}

		println("Database and user created.")

		if err := db.CheckConnection(userConfig); err != nil {
			return errors.Wrap(err, "failed to check user connection")
		}

		println("Connection as user is ok.")

		return nil
	}),
}
