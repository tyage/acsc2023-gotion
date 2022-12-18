package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

func main() {
	publicDir := "./public"
	siteBaseDir := "./sites"

	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		siteId := uuid.New().String()
		newSiteDir := filepath.Join(publicDir, siteBaseDir, siteId)

		err := os.Mkdir(newSiteDir, 0755)
		if err != nil {
			panic(err)
		}

		http.Redirect(w, r, filepath.Join(siteBaseDir, siteId), http.StatusFound)
	})

	http.HandleFunc("/save", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.New("test").Parse("<p>{{.}}</p>")
		if err != nil {
			panic(err)
		}

		f, err := os.Create("/tmp/memo.html")
		if err != nil {
			panic(err)
		}

		err = tmpl.Execute(f, r.FormValue("memo"))
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(w, "done")
	})

	http.Handle("/", http.FileServer(http.Dir(publicDir)))

	http.ListenAndServe(":3000", nil)
}
