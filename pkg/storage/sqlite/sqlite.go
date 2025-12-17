package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"go/types"
	"github.com/kushvardhan/Students-Api/types"
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

func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students WHERE id = ? LIMIT 1")
	if err != nil {
		return types.Student{}, err
	}
	defer stmt.Close()

	var student types.Student

	err = stmt.QueryRow(id).Scan(
		&student.Id,
		&student.Name,
		&student.Email,
		&student.Age,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no student found with id %d", id)
		}
		return types.Student{}, fmt.Errorf("query error: %w", err)
	}

	return student, nil
}

func (s *Sqlite) GetStudents() ([]types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []types.Student

	for rows.Next() {
		var student types.Student
		err := rows.Scan(
			&student.Id,
			&student.Name,
			&student.Email,
			&student.Age,
		)
		if err != nil {
			return nil, err
		}
		students = append(students, student)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}
