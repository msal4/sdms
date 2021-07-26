// Package sdms ...
package sdms

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const (
	FirstSemester = iota + 1
	SecondSemester
)

const (
	FirstStage = iota + 1
	SecondStage
	ThirdStage
	FourthStage
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
	Username string
	Password string
	About    string
}

type Announcement struct {
	ID      int
	Title   string
	Image   string
	Details string
}

type AppStore interface {
	GetSubjects() ([]Subject, error)
	AddSubject(*Subject) error
	UpdateSubject(Subject) error
	GetSubjectsByLecturerID(id int) ([]Subject, error)
	GetSubjectByID(id int) (*Subject, error)
	RemoveSubject(id int) error

	GetLecturers() ([]Lecturer, error)
	AddLecturer(*Lecturer) error
	UpdateLecturer(Lecturer) error
	GetLecturerByID(id int) (*Lecturer, error)
	GetLecturerByUsername(username, password string) (*Lecturer, error)
	RemoveLecturer(id int) error

	GetAnnouncements() ([]Announcement, error)
	AddAnnouncement(*Announcement) error
	RemoveAnnouncement(id int) error
}

type Server struct {
	store AppStore
	http.Handler
}

const storagePath = "./storage"
const pdfPath = storagePath + "/pdf"
const imagesPath = storagePath + "/images"

func NewServer(store AppStore) *Server {
	s := &Server{store: store}

	router := mux.NewRouter()

	apiRouter := router.PathPrefix("/api").Subrouter()

	// subjects
	apiRouter.Handle("/subjects", http.HandlerFunc(s.handleSubjects)).Methods(http.MethodGet)
	apiRouter.Handle("/subjects", http.HandlerFunc(s.handleAddSubject)).Methods(http.MethodPost)
	apiRouter.Handle("/subjects/lecturer/{id:[0-9]+}", http.HandlerFunc(s.handleGetSubjectsByLecturerID)).Methods(http.MethodGet)
	apiRouter.Handle("/subjects/{id:[0-9]+}", http.HandlerFunc(s.handleGetSubjectByID)).Methods(http.MethodGet)
	apiRouter.Handle("/subjects/{id:[0-9]+}", http.HandlerFunc(s.handleRemoveSubject)).Methods(http.MethodDelete)
	apiRouter.Handle("/subjects/{id:[0-9]+}", http.HandlerFunc(s.handleUpdateSubject)).Methods(http.MethodPut)

	// lecturers
	apiRouter.Handle("/lecturers", http.HandlerFunc(s.handleLecturers)).Methods(http.MethodGet)
	apiRouter.Handle("/lecturers", http.HandlerFunc(s.handleAddLecturer)).Methods(http.MethodPost)
	apiRouter.Handle("/lecturers/username/{username}/{password}", http.HandlerFunc(s.handleGetLecturerByUsername)).Methods(http.MethodGet)
	apiRouter.Handle("/lecturers/{id:[0-9]+}", http.HandlerFunc(s.handleGetLecturerByID)).Methods(http.MethodGet)
	apiRouter.Handle("/lecturers/{id:[0-9]+}", http.HandlerFunc(s.handleRemoveLecturer)).Methods(http.MethodDelete)
	apiRouter.Handle("/lecturers/{id:[0-9]+}", http.HandlerFunc(s.handleUpdateLecturer)).Methods(http.MethodPut)

	apiRouter.Handle("/announcements", http.HandlerFunc(s.handleAnnouncements)).Methods(http.MethodGet)
	apiRouter.Handle("/announcements", http.HandlerFunc(s.handleAddAnnouncement)).Methods(http.MethodPost)
	apiRouter.Handle("/announcements/{id:[0-9]+}", http.HandlerFunc(s.handleRemoveAnnouncement)).Methods(http.MethodDelete)

	if _, err := os.Stat(storagePath); os.IsNotExist(err) {
		os.MkdirAll(pdfPath, 0777)
		os.MkdirAll(imagesPath, 0777)
	}

	fs := http.FileServer(http.Dir(storagePath))
	router.PathPrefix("/storage").Handler(http.StripPrefix("/storage", fs))

	allowedMethods := handlers.AllowedMethods([]string{http.MethodGet, http.MethodPut, http.MethodHead, http.MethodPost, http.MethodDelete})
	allowedHeaders := handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "Content-Language", "Origin", "Content-Type"})
	s.Handler = handlers.CORS(allowedMethods, allowedHeaders)(router)

	return s
}

func (s *Server) handleGetLecturerByUsername(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("content-type", "application/json")

	vars := mux.Vars(req)
	l, err := s.store.GetLecturerByUsername(vars["username"], vars["password"])
	if err != nil {
		if err == ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, _ := json.Marshal(l)

	fmt.Fprintf(w, string(b))
}

func (s *Server) handleGetSubjectsByLecturerID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	subjects, err := s.store.GetSubjectsByLecturerID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(subjects); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}


}

func (s *Server) handleGetSubjectByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	vars := mux.Vars(r)
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

func (s *Server) handleAddSubject(w http.ResponseWriter, r *http.Request) {
	var subject Subject
	subject.Lecturer = &Lecturer{}
	subject.Name = r.FormValue("Name")
	subject.Details = r.FormValue("Details")
	subject.Lecturer.ID, _ = strconv.Atoi(r.FormValue("Lecturer"))
	subject.Semester, _ = strconv.Atoi(r.FormValue("Semester"))
	subject.Stage, _ = strconv.Atoi(r.FormValue("Stage"))

	if err := s.store.AddSubject(&subject); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "problem adding subject: %v", err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	syllabusFile, _, err := r.FormFile("Syllabus")
	if err != nil {
		log.Printf("failed to get form file: %v\n", err)
		return
	}
	defer syllabusFile.Close()

	contents, err := ioutil.ReadAll(syllabusFile)
	filepath := path.Join(pdfPath, strconv.Itoa(subject.ID)+".pdf")
	os.WriteFile(filepath, contents, 0777)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = s.store.AddLecturer(&lecturer); err != nil {
		http.Error(w, fmt.Sprintf("problem adding lecturer: %v", err), http.StatusInternalServerError)
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

func (s *Server) handleAnnouncements(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("content-type", "application/json")

	announcements, err := s.store.GetAnnouncements()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(announcements); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleAddAnnouncement(w http.ResponseWriter, r *http.Request) {
	announcement := Announcement{
		Title:   r.FormValue("Title"),
		Details: r.FormValue("Details"),
	}

	if err := s.store.AddAnnouncement(&announcement); err != nil {
		fmt.Fprintf(w, "problem adding announcement: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)

	imageFile, _, err := r.FormFile("Image")
	if err != nil {
		log.Printf("announcement: failed to get form file: %v\n", err)
		return
	}
	defer imageFile.Close()

	contents, err := ioutil.ReadAll(imageFile)
	filepath := path.Join(imagesPath, strconv.Itoa(announcement.ID)+".png")
	os.WriteFile(filepath, contents, 0777)
}

func (s *Server) handleRemoveAnnouncement(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = s.store.RemoveAnnouncement(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
