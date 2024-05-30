package goravel

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func (g *Goravel) MigrateUp(dsn string) error {

	m, err := migrate.New(fmt.Sprintf("file://%s/migrations", g.RootPath), dsn)

	if err != nil {
		return err
	}

	defer m.Close()

	err = m.Up()

	if err != nil {
		g.ErrorLog.Println("Error running migration:", err)
		return err
	}

	return nil

}

func (g *Goravel) MigrateDownAll(dsn string) error {

	m, err := migrate.New(fmt.Sprintf("file://%s/migrations", g.RootPath), dsn)

	if err != nil {
		return err
	}

	defer m.Close()

	err = m.Down()

	if err != nil {
		g.ErrorLog.Println("Error running migration:", err)
		return err
	}

	return nil
}

// Steps runs n steps of migrations. If n > 0, it will run n steps of "up" migrations.
// If n < 0, it will run n steps of "down" migrations.
func (g *Goravel) MigrateSteps(n int, dsn string) error {
	m, err := migrate.New("file://"+g.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Steps(n); err != nil {
		g.ErrorLog.Println("Error running migration:", err)
		return err
	}

	return nil
}

// MigrateForce forces the migration to the immediate last version
func (g *Goravel) MigrateForce(dsn string) error {
	m, err := migrate.New("file://"+g.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Force(-1); err != nil {
		g.ErrorLog.Println("Error running migration:", err)
		return err
	}

	return nil
}
