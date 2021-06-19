// Package sdms ...
package sdms

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

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
	ID       int
	Name     string
	Image    string
	Password string
	About    string
	Subject  *Subject
}

type AppStore interface {
	// Subject
	GetSubjects() ([]Subject, error)
	AddSubject(*Subject) error
	UpdateSubject(Subject) error
	GetSubjectByID(id int) (*Subject, error)
	RemoveSubject(id int) error

	// Lecturer
	GetLecturers() ([]Lecturer, error)
	AddLecturer(*Lecturer) error
	UpdateLecturer(Lecturer) error
	GetLecturerByID(id int) (*Lecturer, error)
	RemoveLecturer(id int) error
}

type Server struct {
	store AppStore
	http.Handler
}

func NewServer(store AppStore) *Server {
	s := &Server{store: store}

	router := mux.NewRouter()

	apiRouter := router.PathPrefix("/api").Subrouter()

	// subjects
	apiRouter.Handle("/subjects", http.HandlerFunc(s.handleSubjects)).Methods(http.MethodGet)
	apiRouter.Handle("/subjects", http.HandlerFunc(s.handleAddSubject)).Methods(http.MethodPost)
	apiRouter.Handle("/subjects/{id:[0-9]+}", http.HandlerFunc(s.handleGetSubjectByID)).Methods(http.MethodGet)
	apiRouter.Handle("/subjects/{id:[0-9]+}", http.HandlerFunc(s.handleRemoveSubject)).Methods(http.MethodDelete)
	apiRouter.Handle("/subjects/{id:[0-9]+}", http.HandlerFunc(s.handleUpdateSubject)).Methods(http.MethodPut)

	// lecturers
	apiRouter.Handle("/lecturers", http.HandlerFunc(s.handleLecturers)).Methods(http.MethodGet)
	apiRouter.Handle("/lecturers", http.HandlerFunc(s.handleAddLecturer)).Methods(http.MethodPost)
	apiRouter.Handle("/lecturers/{id:[0-9]+}", http.HandlerFunc(s.handleGetLecturerByID)).Methods(http.MethodGet)
	apiRouter.Handle("/lecturers/{id:[0-9]+}", http.HandlerFunc(s.handleRemoveLecturer)).Methods(http.MethodDelete)
	apiRouter.Handle("/lecturers/{id:[0-9]+}", http.HandlerFunc(s.handleUpdateLecturer)).Methods(http.MethodPut)

	s.Handler = handlers.CORS()(router)

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

func (s *Server) handleGetLecturerByID(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("content-type", "application/json")

	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "lecturer id is not valid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	lecturer, err := s.store.GetLecturerByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, _ := json.Marshal(lecturer)

	fmt.Fprintf(w, string(b))
}

func (s *Server) handleLecturers(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("content-type", "application/json")

	lecturers, err := s.store.GetLecturers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(lecturers); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleAddLecturer(w http.ResponseWriter, req *http.Request) {
	var lecturer Lecturer
	err := json.NewDecoder(req.Body).Decode(&lecturer)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = s.store.AddLecturer(&lecturer); err != nil {
		fmt.Fprintf(w, "problem adding lecturer: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (s *Server) handleRemoveLecturer(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = s.store.RemoveLecturer(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleUpdateLecturer(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var lecturer Lecturer
	if err := json.NewDecoder(req.Body).Decode(&lecturer); err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	lecturer.ID = id

	if err = s.store.UpdateLecturer(lecturer); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
