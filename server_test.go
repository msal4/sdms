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
	subjects []sdms.Subject
}

func (s *stubStore) GetSubjects() []sdms.Subject {
	return s.subjects
}

func (s *stubStore) AddSubject(subject sdms.Subject) {
	s.subjects = append(s.subjects, subject)
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

		if len(store.subjects) != initialSubjectsLength+1 {
			t.Fatalf("got length of subjects %d, want %d", len(store.subjects), initialSubjectsLength+1)
		}

		assertHasSubject(t, store, subject)

		assertStatusCode(t, res, http.StatusAccepted)
	})

}

func assertHasSubject(t testing.TB, store *stubStore, want sdms.Subject) {
	t.Helper()
	var found bool
	for _, s := range store.subjects {
		if s.Name == want.Name && s.Details == want.Details && s.Lecturer.Name == want.Lecturer.Name &&
			s.Semester == want.Semester && s.Syllabus == want.Syllabus {
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
