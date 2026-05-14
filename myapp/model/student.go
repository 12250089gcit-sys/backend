package model

import (
	"database/sql"
	"errors"

	"myapp/dataStore/postgres"
)

type Student struct {
	StdId     int64  `json:"stdid"`
	FirstName string `json:"fname"`
	LastName  string `json:"lname"`
	Email     string `json:"email"`
}

var (
	ErrStudentIDRequired = errors.New("student id is required")
	ErrFirstNameRequired = errors.New("student first name is required")
	ErrEmailRequired     = errors.New("student email is required")
)

const (
	queryInsertStudentWithID = `INSERT INTO student (stdid, firstname, lastname, email) VALUES ($1, $2, $3, $4)`
	queryInsertStudentAuto   = `INSERT INTO student (firstname, lastname, email) VALUES ($1, $2, $3) RETURNING stdid`
	queryGetStudent          = `SELECT stdid, firstname, lastname, email FROM student WHERE stdid = $1`
	queryUpdateStudent       = `UPDATE student SET firstname=$1, lastname=$2, email=$3 WHERE stdid=$4`
	queryDeleteStudent       = `DELETE FROM student WHERE stdid = $1`
	queryGetAll              = `SELECT stdid, firstname, lastname, email FROM student`
)

func (s *Student) Validate() error {
	if s.FirstName == "" {
		return ErrFirstNameRequired
	}
	if s.Email == "" {
		return ErrEmailRequired
	}
	return nil
}

func (s *Student) ValidateID() error {
	if s.StdId == 0 {
		return ErrStudentIDRequired
	}
	return s.Validate()
}

func (s *Student) Create() error {
	if err := s.Validate(); err != nil {
		return err
	}

	if s.StdId == 0 {
		return postgres.Db.QueryRow(queryInsertStudentAuto, s.FirstName, s.LastName, s.Email).Scan(&s.StdId)
	}

	_, err := postgres.Db.Exec(queryInsertStudentWithID, s.StdId, s.FirstName, s.LastName, s.Email)
	return err
}

func (s *Student) Read() error {
	if err := s.ValidateID(); err != nil {
		return err
	}
	return postgres.Db.QueryRow(queryGetStudent, s.StdId).Scan(&s.StdId, &s.FirstName, &s.LastName, &s.Email)
}

func (s *Student) Update(oldID int64) error {
	if oldID == 0 {
		return ErrStudentIDRequired
	}
	if err := s.Validate(); err != nil {
		return err
	}

	result, err := postgres.Db.Exec(queryUpdateStudent, s.FirstName, s.LastName, s.Email, oldID)
	if err != nil {
		return err
	}
	if rows, err := result.RowsAffected(); err == nil && rows == 0 {
		return sql.ErrNoRows
	}
	return err
}

func (s *Student) Delete() error {
	if err := s.ValidateID(); err != nil {
		return err
	}

	result, err := postgres.Db.Exec(queryDeleteStudent, s.StdId)
	if err != nil {
		return err
	}
	if rows, err := result.RowsAffected(); err == nil && rows == 0 {
		return sql.ErrNoRows
	}
	return err
}

func GetAllStudents() ([]Student, error) {
	rows, err := postgres.Db.Query(queryGetAll)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []Student
	for rows.Next() {
		var s Student
		if err := rows.Scan(&s.StdId, &s.FirstName, &s.LastName, &s.Email); err != nil {
			return nil, err
		}
		students = append(students, s)
	}
	return students, rows.Err()
}
