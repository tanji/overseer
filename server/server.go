package server

import (
	"errors"

	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
)

type Server struct {
	Name         string
	Host         string
	Port         string
	User         string
	Pass         string
	FailureCount uint
	StatusDown   bool
	conn         *sqlx.DB
}

type Servers []*Server

func (s *Server) buildDSN() string {
	opt := "?timeout=1s"
	dsn := s.User + ":" + s.Pass + "@tcp(" + s.Host + ":" + s.Port + ")/" + opt
	return dsn
}

func (s *Server) Open() error {
	var err error
	if s.Port == "" {
		s.Port = "3306"
	}
	s.conn, err = sqlx.Open("mysql", s.buildDSN())
	return err
}

func (s *Server) Ping() error {
	err := s.conn.Ping()
	return err
}

func (s *Server) Close() error {
	err := s.conn.Close()
	return err
}

func (s *Server) Check() error {
	if err := s.Ping(); err != nil {
		return err
	}
	return nil
}

func (s *Server) GetStatus() (map[string]int64, error) {
	type status struct {
		Name  string
		Value int64
	}
	vars := make(map[string]int64)
	rows, err := s.conn.Queryx("SELECT lower(Variable_name) AS name, Variable_Value AS value FROM information_schema.global_status")
	if err != nil {
		return nil, errors.New("Database error: could not get status")
	}
	for rows.Next() {
		v := status{}
		rows.Scan(&v.Name, &v.Value)
		vars[v.Name] = v.Value
	}
	return vars, nil
}
