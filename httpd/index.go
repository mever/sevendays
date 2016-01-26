package httpd

import (
	"fmt"
	"net/http"
	"text/template"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("index.html").ParseFiles(AssetsDir + "/index.html")
	if err != nil {
		w.Write([]byte(fmt.Sprintln(err)))
	} else {

		if err := t.Execute(w, nil); err != nil {
			w.Write([]byte(fmt.Sprintln(err)))
		}
	}
}
