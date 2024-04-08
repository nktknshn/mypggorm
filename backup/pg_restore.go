package backup

import (
	"os"
	"os/exec"

	"github.com/nktknshn/mypggorm"
)

type PostgresBackuperPgRestore struct {
	pgRestorePath string
}

func NewPgRestore(pgRestorePath string) PostgresBackuperPgRestore {
	return PostgresBackuperPgRestore{
		pgRestorePath: pgRestorePath,
	}
}

func (p PostgresBackuperPgRestore) Restore(cfg mypggorm.DatabaseConnectionConfig, backupPath string) error {

	cmd := exec.Cmd{
		Path: p.pgRestorePath,
		Args: []string{
			p.pgRestorePath,
			"-h", cfg.Host,
			"-p", cfg.Port,
			"-U", cfg.User,
			"-d", cfg.Dbname,
			"-1", backupPath},
		Env: []string{
			"PGPASSWORD=" + cfg.Password,
		}}

	println(cmd.String())

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

type PostgresBackuperPsqlRestore struct {
	psqlPath string
}

func NewPsqlRestore(psqlPath string) PostgresBackuperPsqlRestore {
	return PostgresBackuperPsqlRestore{
		psqlPath: psqlPath,
	}
}

func (p PostgresBackuperPsqlRestore) Restore(cfg mypggorm.DatabaseConnectionConfig, backupPath string) error {

	cmd := exec.Cmd{
		Path: p.psqlPath,
		Args: []string{
			p.psqlPath,
			"-h", cfg.Host,
			"-p", cfg.Port,
			"-U", cfg.User,
			"-d", cfg.Dbname,
			"-f", backupPath},
		Env: []string{
			"PGPASSWORD=" + cfg.Password,
		}}

	println(cmd.String())

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
