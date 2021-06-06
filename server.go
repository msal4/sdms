// Package sdms ...
package sdms

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type AppStore interface {
	GetSubjects() []Subject
	AddSubject(Subject)
}

type Server struct {
	store AppStore
	http.Handler
}

func NewServer(store AppStore) *Server {
	s := &Server{store: store}

	router := mux.NewRouter()

	router.Handle("/subjects", http.HandlerFunc(s.handleSubjects)).Methods(http.MethodGet)
	router.Handle("/subjects", http.HandlerFunc(s.handleAddSubject)).Methods(http.MethodPost)

	s.Handler = router

	return s
}

func (s *Server) handleSubjects(w http.ResponseWriter, req *http.Request) {
	err := json.NewEncoder(w).Encode(s.store.GetSubjects())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleAddSubject(w http.ResponseWriter, req *http.Request) {
	var subject Subject
	err := json.NewDecoder(req.Body).Decode(&subject)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.store.AddSubject(subject)

	w.WriteHeader(http.StatusAccepted)
}

const (
	FirstSemester = iota + 1
	SecondSemester
)

type Subject struct {
	Name     string
	Details  string
	Lecturer *Lecturer
	Semester int
	Syllabus string
}

type Lecturer struct {
	Name    string
	Image   string
	About   string
	Subject *Subject
}
