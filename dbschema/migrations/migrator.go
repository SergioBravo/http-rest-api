package migrations

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"text/template"

	migrate "github.com/rubenv/sql-migrate"

	"github.com/sirupsen/logrus"
)

func Migrate() {
	migrations := &migrate.AssetMigrationSource{
		Asset:    Asset,
		AssetDir: AssetDir,
		Dir:      "migrations",
	}
}

type Migration struct {
	Version string
	Up      func(*sql.Tx) error
	Down    func(*sql.Tx) error

	done bool
}

type Migrator struct {
	db         *sql.DB
	Versions   []string
	Migrations map[string]*Migration
}

var logger = logrus.New()

var migrator = &Migrator{
	Versions:   []string{},
	Migrations: map[string]*Migration{},
}

func (m *Migrator) AddMigration(mg *Migration) {
	m.Migrations[mg.Version] = mg

	index := 0
	for index < len(m.Versions) {
		if m.Versions[index] > mg.Version {
			break
		}

		index++
	}

	m.Versions = append(m.Versions, mg.Version)
	copy(m.Versions[index+1:], m.Versions[index:])
	m.Versions[index] = mg.Version
}

func Create(name string) error {
	version := time.Now().Format("001")

	in := struct {
		Version string
		Name    string
	}{
		Version: version,
		Name:    name,
	}

	var out bytes.Buffer

	t := template.Must(template.ParseFiles("./dbschema/migrations/template.txt"))
	if err := t.Execute(&out, in); err != nil {
		return errors.New("Unable to execute template: " + err.Error())
	}

	f, err := os.Create(fmt.Sprintf("./dbschema/migrations/%s_%s.go", version, name))
	if err != nil {
		return errors.New("Unable to create migration file: " + err.Error())
	}

	defer f.Close()

	if _, err := f.WriteString(out.String()); err != nil {
		return errors.New("Unable to write migration file: " + err.Error())
	}

	fmt.Println("Generated new migration files ... ", f.Name())
	return nil
}

func Init(db *sql.DB) (*Migrator, error) {
	migrator.db = db

	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS schema_migrations (version varchar(255));"); err != nil {
		logger.WithError(err).Error("Unable to creatte `schema_migrations` table")
		return migrator, err
	}

	rows, err := db.Query("SELECT version FROM schema_migrations;")
	if err != nil {
		logger.WithError(err).Error("Unable to select `version` from `schema_migrations` table")
		return migrator, err
	}

	defer rows.Close()

	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return migrator, err
		}

		if migrator.Migrations[version] != nil {
			migrator.Migrations[version].done = true
		}
	}

	return migrator, nil
}

func (m *Migrator) Up(step int) error {
	tx, err := m.db.BeginTx(context.TODO(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	count := 0
	for _, v := range m.Versions {
		if step > 0 && count == step {
			break
		}

		mg := m.Migrations[v]

		if mg.done {
			continue
		}

		logger.Info("Running migration", mg.Version)

		if err := mg.Up(tx); err != nil {
			tx.Rollback()
			return err
		}

		if _, err := tx.Exec("INSERT INTO schema_migrations(version) VALUES ($1);", mg.Version); err != nil {
			tx.Rollback()
			return err
		}

		logger.Info("Finished running migration", mg.Version)

		count++
	}

	tx.Commit()

	return nil
}

func (m *Migrator) Down(step int) error {
	tx, err := m.db.BeginTx(context.TODO(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	count := 0
	for _, v := range reverse(m.Versions) {
		if step > 0 && count == step {
			break
		}

		mg := m.Migrations[v]

		if !mg.done {
			continue
		}

		logger.Info("Reverting Migration", mg.Version)

		if err := mg.Down(tx); err != nil {
			tx.Rollback()
			return err
		}

		if _, err := tx.Exec("DELETE FROM schema_migrations WHERE version = $1;", mg.Version); err != nil {
			tx.Rollback()
			return err
		}

		logger.Info("Finished reverting migration", mg.Version)

		count++
	}

	tx.Commit()

	return nil
}

func (m *Migrator) MigrationStatus() {
	for _, v := range m.Versions {
		mg := m.Migrations[v]

		if mg.done {
			logger.Info(fmt.Sprintf("Migration %s... completed", v))
		} else {
			logger.Info(fmt.Sprintf("Migration %s... pending", v))
		}
	}
}

func reverse(arr []string) []string {
	for i := 0; i < len(arr)/2; i++ {
		j := len(arr) - i - 1
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}
