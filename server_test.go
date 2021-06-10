package sdms_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/msal4/sdms"
)

type stubStore struct {
	subjects  []sdms.Subject
	lecturers []sdms.Lecturer
}

func (s *stubStore) GetSubjects() ([]sdms.Subject, error) {
	return s.subjects, nil
}

func (s *stubStore) AddSubject(subject *sdms.Subject) error {
	if subject != nil {
		s.subjects = append(s.subjects, *subject)
	}
	return nil
}

func (s *stubStore) RemoveSubject(id int) error {
	return nil
}

func (s *stubStore) UpdateSubject(subject sdms.Subject) error {
	return nil
}

func (s *stubStore) GetSubjectByID(id int) (*sdms.Subject, error) {
	return nil, nil
}

func (s *stubStore) GetLecturers() ([]sdms.Lecturer, error) {
	return s.lecturers, nil
}

func (s *stubStore) AddLecturer(lecturer *sdms.Lecturer) error {
	if lecturer != nil {
		s.lecturers = append(s.lecturers, *lecturer)
	}
	return nil
}

func (s *stubStore) RemoveLecturer(id int) error {
	return nil
}

func (s *stubStore) UpdateLecturer(lecturer sdms.Lecturer) error {
	return nil
}

func (s *stubStore) GetLecturerByID(id int) (*sdms.Lecturer, error) {
	return nil, nil
}

func TestGETSubjects(t *testing.T) {
	t.Run("get empty subjects", func(t *testing.T) {
		want := []sdms.Subject{}

		store := &stubStore{
			subjects: want,
		}

		srv := sdms.NewServer(store)
		req := createGetSubjectsRequest()
		res := httptest.NewRecorder()

		srv.ServeHTTP(res, req)

		got := decodeSubjectsResponseBody(res)

		assertSubjectsEqual(t, got, want)
		assertStatusCode(t, res, http.StatusOK)
	})

	t.Run("get subjects", func(t *testing.T) {
		wantedSubjects := []sdms.Subject{
			{
				Name:     "Math",
				Details:  "Math Subject",
				Semester: sdms.FirstSemester,
				Syllabus: "math syllabus",
				Lecturer: &sdms.Lecturer{
					Name:  "Mohammed",
					Image: "/images/1.png",
					About: "I'm a doctor",
				},
			},
			{
				Name:     "Physics",
				Details:  "Physics Subject",
				Semester: sdms.FirstSemester,
				Syllabus: "physics syllabus",
				Lecturer: &sdms.Lecturer{
					Name:  "Salman",
					Image: "/images/2.png",
					About: "I'm a professor",
				},
			},
		}

		store := &stubStore{
			subjects: wantedSubjects,
		}

		srv := sdms.NewServer(store)
		req := createGetSubjectsRequest()
		res := httptest.NewRecorder()

		srv.ServeHTTP(res, req)

		got := decodeSubjectsResponseBody(res)

		assertSubjectsEqual(t, got, wantedSubjects)
		assertStatusCode(t, res, http.StatusOK)
	})
}

func decodeSubjectsResponseBody(res *httptest.ResponseRecorder) (subjects []sdms.Subject) {
	json.NewDecoder(res.Body).Decode(&subjects)
	return
}

func createGetSubjectsRequest() *http.Request {
	req, _ := http.NewRequest("GET", "/subjects", nil)
	return req
}

func TestPOSTSubject(t *testing.T) {
	subjects := []sdms.Subject{
		{
			Name:     "Math",
			Details:  "Math Subject",
			Semester: sdms.FirstSemester,
			Syllabus: "math syllabus",
			Lecturer: &sdms.Lecturer{
				Name:  "Mohammed",
				Image: "/images/1.png",
				About: "I'm a doctor",
			},
		},
		{
			Name:     "Physics",
			Details:  "Physics Subject",
			Semester: sdms.FirstSemester,
			Syllabus: "physics syllabus",
			Lecturer: &sdms.Lecturer{
				Name:  "Salman",
				Image: "/images/2.png",
				About: "I'm a professor",
			},
		},
	}

	store := &stubStore{
		subjects: subjects,
	}

	srv := sdms.NewServer(store)

	t.Run("add subject", func(t *testing.T) {
		initialSubjectsLength := len(store.subjects)

		subject := sdms.Subject{
			Name:     "Geography",
			Details:  "Meta",
			Semester: sdms.FirstSemester,
			Syllabus: "Hi",
			Lecturer: &sdms.Lecturer{
				Name:  "Sub",
				Image: "/images/2.png",
				About: "I'm a professor",
			},
		}

		b := bytes.NewBuffer(nil)
		json.NewEncoder(b).Encode(subject)

		req := createPOSTSubjectsRequest(b)

		res := httptest.NewRecorder()

		srv.ServeHTTP(res, req)

		assertStatusCode(t, res, http.StatusAccepted)

		assertSubjectsLength(t, store.subjects, initialSubjectsLength+1)

		assertHasSubject(t, store.subjects, subject)
	})

}

func assertSubjectsLength(t testing.TB, s []sdms.Subject, want int) {
	t.Helper()
	if len(s) != want {
		t.Fatalf("got length of subjects %d, want %d", len(s), want)
	}
}

func assertHasSubject(t testing.TB, subjects []sdms.Subject, want sdms.Subject) {
	t.Helper()
	var found bool
	for _, s := range subjects {
		if s.Name == want.Name && s.Details == want.Details && s.Semester == want.Semester && s.Syllabus == want.Syllabus {
			found = true
		}
	}

	if !found {
		t.Fatalf("store does not contian subject %#v", want)
	}

}

func createPOSTSubjectsRequest(body io.Reader) *http.Request {
	req, _ := http.NewRequest("POST", "/subjects", body)
	return req
}

func assertSubjectsEqual(t testing.TB, got, want []sdms.Subject) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func assertStatusCode(t testing.TB, res *httptest.ResponseRecorder, want int) {
	t.Helper()
	if res.Code != want {
		t.Fatalf("got status %d, want %d", res.Code, want)
	}
}
