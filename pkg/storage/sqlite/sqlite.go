package sqlite

import (
	"context"
	"database/sql"

	"github.com/kushvardhan/Students-Api/pkg/config"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/mod/sumdb/storage"
)

type Sqlite struct {
	Db *sql.DB
}

// ReadOnly implements storage.Storage.
func (s *Sqlite) ReadOnly(ctx context.Context, f func(context.Context, storage.Transaction) error) error {
	panic("unimplemented")
}

// ReadWrite implements storage.Storage.
func (s *Sqlite) ReadWrite(ctx context.Context, f func(context.Context, storage.Transaction) error) error {
	panic("unimplemented")
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	email TEXT,
	age INTEGER
	)`)

	if err != nil {
		return nil, err
	}

	return &Sqlite{
		Db: db,
	}, nil
}

func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {
	stmt, err := s.Db.Prepare("INSERT INTO students (name, email, age) VALUES (7,7,7)")

	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}
	lastId, err := res.LastInsertId()

	return lastId, nil
}
