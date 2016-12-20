package server

import "testing"

func TestOpenServer(t *testing.T) {
	s := Server{
		Host: "127.0.0.1",
		User: "root",
		Pass: "admin",
	}
	t.Log("DSN:", s.buildDSN())
	err := s.Open()
	if err != nil {
		t.Fatal("Could not open test db")
	}
	t.Logf("Connection pointer is %v", s.conn)
	err = s.Ping()
	if err != nil {
		t.Fatal("Could not ping test db")
	}
}
