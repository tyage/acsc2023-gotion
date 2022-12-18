package main

import (
	"errors"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type NoteParam struct {
	Id   uuid.UUID
	Body string
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

func GetNotePath(noteId uuid.UUID) (string, string) {
	noteFile := noteId.String() + ".html"
	notePublicPath := filepath.Join(NoteBaseDir, noteFile)
	noteFilePath := filepath.Join(PublicDir, notePublicPath)

	return noteFilePath, notePublicPath
}

func ValidateNoteBody(body string) error {
	if len(body) > 1024 {
		return errors.New("note is too long")
	}
	return nil
}

func main() {

	http.HandleFunc("/new-note", func(w http.ResponseWriter, r *http.Request) {
		body := r.FormValue("body")

		// validate body
		err := ValidateNoteBody(body)
		if err != nil {
			http.Error(w, "invalid body", http.StatusInternalServerError)
			return
		}

		noteId := uuid.New()
		noteFilePath, noteURL := GetNotePath(noteId)

		f, err := os.Create(noteFilePath)
		if err != nil {
			http.Error(w, "failed note creation", http.StatusInternalServerError)
			return
		}

		WriteNote(f, NoteParam{Id: noteId, Body: body})

		http.Redirect(w, r, noteURL, http.StatusFound)
	})

	http.HandleFunc("/update-note", func(w http.ResponseWriter, r *http.Request) {
		noteId := r.FormValue("noteId")
		body := r.FormValue("body")

		// validate noteId
		uuidNoteId, err := uuid.Parse(noteId)
		if err != nil {
			http.Error(w, "invalid noteId", http.StatusInternalServerError)
			return
		}

		// validate body
		err = ValidateNoteBody(body)
		if err != nil {
			http.Error(w, "invalid body", http.StatusInternalServerError)
			return
		}

		noteFilePath, noteURL := GetNotePath(uuidNoteId)

		f, err := os.OpenFile(noteFilePath, os.O_WRONLY, 0644)
		if err != nil {
			http.Error(w, "invalid note", http.StatusInternalServerError)
			return
		}

		WriteNote(f, NoteParam{Id: uuidNoteId, Body: body})

		http.Redirect(w, r, noteURL, http.StatusFound)
	})

	http.Handle("/", http.FileServer(http.Dir(PublicDir)))

	http.ListenAndServe(":3000", nil)
}
