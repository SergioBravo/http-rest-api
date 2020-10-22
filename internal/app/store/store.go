package store

import (
	"database/sql"

	_ "github.com/lib/pq" // ...
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

type Store struct {
	config *Config
	db     *sql.DB
}

func New(config *Config) *Store {
	return &Store{
		config: config,
	}
}

func (s *Store) Open() error {
	db, err := sql.Open("postgres", s.config.DataBaseUrl)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	s.db = db

	return nil
}

func (s *Store) Close() {
	s.db.Close()
}

func NewDB() *sql.DB {
	logger.Info("Connecting to database")
	db, err := sql.Open("postgres", "postgres://sergio:abkjcjabz24@localhost/restapi_dev?sslmode=disable")
	if err != nil {
		logger.WithError(err).Error("Unable to connect to database")
		return nil
	}

	if err := db.Ping(); err != nil {
		logger.WithError(err).Error("Unable to connect to database")
		return nil
	}

	logger.Info("Database connected!")

	return db
}
