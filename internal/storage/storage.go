package storage

import "github.com/d1vyanshu-kumar/students-api/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentByID(id int64) (types.Student, error)        // implement this method to get student by ID in sqlite.go file
	GetStudents() ([]types.Student, error) // implement this method to get list of students in sqlite.go file
}