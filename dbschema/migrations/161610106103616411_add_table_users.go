package migrations

import "database/sql"

func init() {
	migrator.AddMigration(&Migration{
		Version: "161610106103616411",
		Up:      mig_161610106103616411_add_table_users_up,
		Down:    mig_161610106103616411_add_table_users_down,
	})
}

func mig_161610106103616411_add_table_users_up(tx *sql.Tx) error {
	_, err := tx.Exec(`CREATE TABLE users (
		id bigserial not null primary key,
		email varchar not null unique,
		encrypted_password varchar not null
	);`)

	if err != nil {
		return err
	}
	return nil
}

func mig_161610106103616411_add_table_users_down(tx *sql.Tx) error {
	_, err := tx.Exec("DROP TABLE users;")
	if err != nil {
		return err
	}
	return nil
}
