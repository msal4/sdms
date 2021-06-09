// Package sdms ...
package sdms

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type AppStore interface {
	GetSubjects() ([]Subject, error)
	AddSubject(*Subject) error
	UpdateSubject(Subject) error
	GetSubjectByID(id int) (*Subject, error)
	RemoveSubject(id int) error
}

type Server struct {
	store AppStore
	http.Handler
}

func NewServer(store AppStore) *Server {
	s := &Server{store: store}

	router := mux.NewRouter()

	// subjects
	router.Handle("/subjects", http.HandlerFunc(s.handleSubjects)).Methods(http.MethodGet)
	router.Handle("/subjects", http.HandlerFunc(s.handleAddSubject)).Methods(http.MethodPost)
	router.Handle("/subjects/{id:[0-9]+}", http.HandlerFunc(s.handleGetSubjectByID)).Methods(http.MethodGet)
	router.Handle("/subjects/{id:[0-9]+}", http.HandlerFunc(s.handleRemoveSubject)).Methods(http.MethodDelete)
	router.Handle("/subjects/{id:[0-9]+}", http.HandlerFunc(s.handleUpdateSubject)).Methods(http.MethodPut)

	s.Handler = router

	return s
}

func (s *Server) handleGetSubjectByID(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("content-type", "application/json")

	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "subject id is not valid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	subject, err := s.store.GetSubjectByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, _ := json.Marshal(subject)

	fmt.Fprintf(w, string(b))
}

func (s *Server) handleSubjects(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("content-type", "application/json")

	subjects, err := s.store.GetSubjects()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(subjects); err != nil {
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

	if err = s.store.AddSubject(&subject); err != nil {
		fmt.Fprintf(w, "problem adding subject: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (s *Server) handleRemoveSubject(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = s.store.RemoveSubject(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleUpdateSubject(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var subject Subject
	if err := json.NewDecoder(req.Body).Decode(&subject); err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	subject.ID = id

	if err = s.store.UpdateSubject(subject); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

const (
	FirstSemester = iota + 1
	SecondSemester
)

type Subject struct {
	ID       int
	Name     string
	Details  string
	Stage    int
	Lecturer *Lecturer
	Semester int
	Syllabus string
}

type Lecturer struct {
	ID      int
	Name    string
	Image   string
	About   string
	Subject *Subject
}
