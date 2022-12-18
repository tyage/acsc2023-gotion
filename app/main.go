package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
)

func main() {
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
	http.Handle("/", http.FileServer(http.Dir("/tmp/")))
	http.ListenAndServe(":3000", nil)
}
