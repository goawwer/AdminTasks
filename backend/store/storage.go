package store

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/goawwer/admintasks/config"
)

type Storage struct {
	databaseURI    string
	db             *sql.DB
	taskRepository *TaskRepository
}

func New(cfg *config.Config) *Storage {
	return &Storage{
		databaseURI: cfg.DatabaseURI(),
	}
}

func (s *Storage) Open() error {
	db, err := sql.Open("postgres", s.databaseURI)
	if err != nil {
		return fmt.Errorf("cannot validate database arguments: %w", err)
	}

	if err = db.Ping(); err != nil {
		return fmt.Errorf("cannot connect to database: %w", err)
	}

	s.db = db
	log.Println("database connection created successfully")
	return nil
}

func (s *Storage) Close() {
	s.db.Close()
}

func (s *Storage) TaskRepository() *TaskRepository {
	if s.taskRepository != nil {
		return s.taskRepository
	}

	s.taskRepository = &TaskRepository{
		storage: s,
	}

	return s.taskRepository
}
