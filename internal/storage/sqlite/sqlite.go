package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/d1vyanshu-kumar/students-api/internal/config"
	"github.com/d1vyanshu-kumar/students-api/internal/types"
	_ "github.com/mattn/go-sqlite3" // sqlite driver, we are using this underhood thus we are using blank identifier. remember that
)

type SQLite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*SQLite, error) {

	db, err := sql.Open("sqlite3", cfg.StoragePath) // here we are opening the sqlite database and the path is provided in the config file.
	if err != nil {
		return nil, err
	}

	// we need to first create table here:

	_, er := db.Exec(`
		CREATE TABLE IF NOT EXISTS students (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT NOT NULL,
			age INTEGER NOT NULL
		);
	`)

	if er != nil {
		return nil, er // we are usinf nill here cause we there is no sqlite instance to return
	}

	return &SQLite{Db: db}, nil
}

// we are going to implement the CreateStudent method here so that we can satisfy the Storage interface. and use that into the main.go file.
func (s *SQLite) CreateStudent(name string, email string, age int) (int64, error) {

	// here we are going to insert the student data into the database. write here the sql query to insert the data into the database.

	stmt, err := s.Db.Prepare(`INSERT INTO students (name, email, age) VALUES (?, ?, ?)`)

	if err != nil {
		return 0, err
	}
	// we need to close that statment  here we go:

	defer stmt.Close()
	// now as above we can see that we have prepared the statement now we need to execute that statement.
	// here we are using EXEC because we are not expecting any rows to be returned from the database.

	result, err := stmt.Exec(name, email, age) // bind here in same order as the above query

	if err != nil {
		return 0, err
	}

	lastID, err := result.LastInsertId() // this will return the last inserted id

	if err != nil {
		return 0, err // empty value for int64 is 0
	}

	// and if there is not error then here  then last id is got so we simply return that id

	return lastID, nil
}

// implementing the GetStudentByID method here

func (s *SQLite) GetStudentByID(id int64) (types.Student, error) {
	// here we are going to see how can we run the database query to get the student by ID

	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students WHERE id = ? LIMIT 1")

	// 	Performance difference: If your users table has 1 million records and the matching user is the 10th record:
	// Without LIMIT 1: Database checks all 1,000,000 records
	// With LIMIT 1: Database stops at record #10

	if err != nil {
		return types.Student{}, err
	}
	// now close the statment

	defer stmt.Close()

	// now the data which is coming out from the database we need to serilized that and put that inside a structs.
	// using QueryRow method and then use scan method

	var student types.Student

	// we are passing here reference directly
	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)

	if err != nil {
		// what if ther is no user found.
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no student found with id %s", fmt.Sprint(id)) // conver id into string.
		}
		return types.Student{}, fmt.Errorf("query Error: %w", err)
	}

	return student, nil

}
