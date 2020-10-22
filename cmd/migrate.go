package cmd

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/SergioBravo/http-rest-api/dbschema/migrations"
	"github.com/SergioBravo/http-rest-api/internal/app/apiserver"
	"github.com/SergioBravo/http-rest-api/internal/app/store"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	configPath string
)

var logger = logrus.New()

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "database migrations tool",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "starting app",
	Run: func(cmd *cobra.Command, args []string) {
		config := apiserver.NewConfig()
		_, err := toml.DecodeFile(configPath, config)
		if err != nil {
			log.Fatal(err)
		}

		db := store.NewDB()

		migrator, err := migrations.Init(db)
		if err != nil {
			logger.WithError(err).Error("Unable to fetch migrator")
			return
		}

		if err = migrator.Up(0); err != nil {
			logger.WithError(err).Error("Unable to run `up` migrations")
			return
		}

		s := apiserver.New(config)
		if err := s.Start(); err != nil {
			log.Fatal(err)
		}
	},
}

var migrateCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create a new empty migrations file",
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			logger.WithError(err).Error("Unable to read flag `name`")
			return
		}

		if err := migrations.Create(name); err != nil {
			logger.WithError(err).Error("Unable to create migration")
			return
		}
	},
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "run up migrations",
	Run: func(cmd *cobra.Command, args []string) {
		step, err := cmd.Flags().GetInt("step")
		if err != nil {
			logger.WithError(err).Error("Unable to read flag `step`")
			return
		}

		db := store.NewDB()

		migrator, err := migrations.Init(db)
		if err != nil {
			logger.WithError(err).Error("Unable to fetch migrator")
			return
		}

		if err = migrator.Up(step); err != nil {
			logger.WithError(err).Error("Unable to run `up` migrations")
			return
		}
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "run down migrations",
	Run: func(cmd *cobra.Command, args []string) {
		step, err := cmd.Flags().GetInt("step")
		if err != nil {
			logger.WithError(err).Error("Unable to read flag `step`")
			return
		}

		db := store.NewDB()

		migrator, err := migrations.Init(db)
		if err != nil {
			logger.WithError(err).Error("Unable to fetch migrator")
			return
		}

		if err = migrator.Down(step); err != nil {
			logger.WithError(err).Error("Unable to run `down` migrations")
			return
		}
	},
}

var migrateStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "display status of each migrations",
	Run: func(cmd *cobra.Command, args []string) {
		db := store.NewDB()

		migrator, err := migrations.Init(db)
		if err != nil {
			logger.WithError(err).Error("Unable to fetch migrator")
			return
		}

		migrator.MigrationStatus()

		return
	},
}

func init() {
	flag.StringVar(&configPath, "config-path",
		"configs/apiserver.toml",
		"path to config file")
	// Add "--name" flag to "create" command
	migrateCreateCmd.Flags().StringP("name", "n", "", "Name for the migration")

	// Add "--step" flag to both "up" and "down" command
	migrateUpCmd.Flags().IntP("step", "s", 0, "Number of migrations to execute")
	migrateDownCmd.Flags().IntP("step", "s", 0, "Number of migrations to execute")

	// Add "create", "up" and "down" commands to the "migrate" command
	migrateCmd.AddCommand(migrateUpCmd, migrateDownCmd, migrateCreateCmd, migrateStatusCmd)

	// Add "migrate" command to the root command
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(serveCmd)
}
