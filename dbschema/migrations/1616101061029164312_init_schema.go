package migrations

import "database/sql"

func init() {
	migrator.AddMigration(&Migration{
		Version: "1616101061029164312",
		Up:      mig_1616101061029164312_init_schema_up,
		Down:    mig_1616101061029164312_init_schema_down,
	})
}

func mig_1616101061029164312_init_schema_up(tx *sql.Tx) error {
	_, err := tx.Exec("CREATE TABLE test (id bigserial not null primary key);")

	if err != nil {
		return err
	}

	return nil
}

func mig_1616101061029164312_init_schema_down(tx *sql.Tx) error {
	_, err := tx.Exec("DROP TABLE test;")
	if err != nil {
		return err
	}
	return nil
}
