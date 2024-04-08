package backup

import (
	"os"
	"os/exec"

	"github.com/nktknshn/mypggorm"
)

type PostgresBackuperPgDump struct {
	pgDumpPath string
}

func NewPgDump(pgDumpPath string) PostgresBackuperPgDump {
	return PostgresBackuperPgDump{
		pgDumpPath: pgDumpPath,
	}
}

func (p PostgresBackuperPgDump) Backup(cfg mypggorm.DatabaseConnectionConfig, backupPath string) error {

	cmd := exec.Cmd{
		Path: p.pgDumpPath,
		Args: []string{
			p.pgDumpPath,
			"-h", cfg.Host,
			"-p", cfg.Port,
			"-U", cfg.User,
			"-f", backupPath,
			cfg.Dbname},
		Env: []string{
			"PGPASSWORD=" + cfg.Password,
		}}

	println(cmd.String())

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
