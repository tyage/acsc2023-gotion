package main

import (
	"errors"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

type NoteParam struct {
	Id    string
	Title string
	Body  string
}

const (
	PublicDir    = "./public"
	NoteBaseDir  = "./notes"
	NoteTemplate = "./templates/note.html"
)

func WriteNote(file *os.File, body NoteParam) {
	tmpl, err := template.ParseFiles(NoteTemplate)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(file, body)
	if err != nil {
		panic(err)
	}
}

func GenerateNoteId(title string) string {
	prefix := uuid.New().String()
	id := prefix + "-" + strings.ReplaceAll(title, " ", "-")
	return id
}

func GetNotePath(noteId string) (string, string) {
	notePublicPath := filepath.Join(NoteBaseDir, noteId)
	noteFilePath := filepath.Join(PublicDir, notePublicPath)

	return noteFilePath, notePublicPath
}

func ValidateNoteBody(body string) error {
	if len(body) > 1024 {
		return errors.New("note is too long")
	}
	return nil
}
func ValidateNoteTitle(title string) error {
	if len(title) > 20 {
		return errors.New("title is too long")
	}

	validTitleRegExp := regexp.MustCompile("^[a-zA-Z0-9 ]+$")
	if !validTitleRegExp.MatchString(title) {
		return errors.New("title is invalid")
	}
	return nil
}
func ValidateNoteId(id string) error {
	if len(id) > 60 {
		return errors.New("id is too long")
	}

	validTitleRegExp := regexp.MustCompile("^[a-zA-Z0-9-]+$")
	if !validTitleRegExp.MatchString(id) {
		return errors.New("id is invalid")
	}
	return nil
}

func main() {

	http.HandleFunc("/new-note", func(w http.ResponseWriter, r *http.Request) {
		body := r.FormValue("body")
		title := r.FormValue("title")

		// validate body
		err := ValidateNoteBody(body)
		if err != nil {
			http.Error(w, "invalid body", http.StatusInternalServerError)
			return
		}

		// validate title
		err = ValidateNoteTitle(title)
		if err != nil {
			http.Error(w, "invalid title", http.StatusInternalServerError)
			return
		}

		noteId := GenerateNoteId(title)
		noteFilePath, noteURL := GetNotePath(noteId)

		f, err := os.Create(noteFilePath)
		if err != nil {
			http.Error(w, "failed note creation", http.StatusInternalServerError)
			return
		}

		WriteNote(f, NoteParam{Id: noteId, Title: title, Body: body})

		http.Redirect(w, r, noteURL, http.StatusFound)
	})

	http.HandleFunc("/update-note", func(w http.ResponseWriter, r *http.Request) {
		noteId := r.FormValue("noteId")
		title := r.FormValue("title")
		body := r.FormValue("body")

		// validate noteId
		err := ValidateNoteId(noteId)
		if err != nil {
			http.Error(w, "invalid id", http.StatusInternalServerError)
			return
		}

		// validate body
		err = ValidateNoteTitle(title)
		if err != nil {
			http.Error(w, "invalid title", http.StatusInternalServerError)
			return
		}

		// validate body
		err = ValidateNoteBody(body)
		if err != nil {
			http.Error(w, "invalid body", http.StatusInternalServerError)
			return
		}

		noteFilePath, noteURL := GetNotePath(noteId)

		f, err := os.OpenFile(noteFilePath, os.O_WRONLY, 0644)
		if err != nil {
			http.Error(w, "invalid note", http.StatusInternalServerError)
			return
		}

		WriteNote(f, NoteParam{Id: noteId, Title: title, Body: body})

		http.Redirect(w, r, noteURL, http.StatusFound)
	})

	http.Handle("/", http.FileServer(http.Dir(PublicDir)))

	http.ListenAndServe(":3000", nil)
}
