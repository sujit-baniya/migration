package migrations

import (
	"database/sql"
	"fmt"
)

func ApplyMigrations(dialect string, curBD *sql.DB, dir MigrationDirection, dryrun bool, limit int) error {
	source := EmbedFileSystemMigrationSource{
		FileSystem: DefaultConfig.EmbeddedFS,
		Root:       DefaultConfig.Dir,
	}
	if dryrun {
		migrations, _, err := PlanMigration(curBD, dialect, source, dir, limit)
		if err != nil {
			return fmt.Errorf("Cannot plan migration: %s", err)
		}

		for _, m := range migrations {
			PrintMigration(m, dir)
		}
	} else {
		n, err := ExecMax(curBD, dialect, source, dir, limit)
		if err != nil {
			return fmt.Errorf("Migration failed: %s", err)
		}

		if n == 1 {
			ui.Output("Applied 1 migration")
		} else {
			ui.Output(fmt.Sprintf("Applied %d migrations", n))
		}
	}

	return nil
}

func PrintMigration(m *PlannedMigration, dir MigrationDirection) {
	if dir == Up {
		ui.Output(fmt.Sprintf("==> Would apply migration %s (up)", m.Id))
		for _, q := range m.Up {
			ui.Output(q)
		}
	} else if dir == Down {
		ui.Output(fmt.Sprintf("==> Would apply migration %s (down)", m.Id))
		for _, q := range m.Down {
			ui.Output(q)
		}
	} else {
		panic("Not reached")
	}
}
