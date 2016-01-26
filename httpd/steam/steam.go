package steam

import (
	"fmt"
	"github.com/mever/sevendaystodie"
	"github.com/mever/steam/cmd"
	"net/http"
	"text/template"
)

var installer = &cmd.Installer{}
var assetsDir string

func Setup(assets string) {
	assetsDir = assets
	cmd.AddQuestion("Steam Guard code:", "What is your Steam Guard code?", true)
	http.HandleFunc("/steam", steamHandler)
}

type page struct {
	Refresh  bool
	Name     string
	Status   string
	Action   string
	Question *cmd.Question
}

func steamHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		switch r.Form.Get("action") {
		case "install":
			installer.Install(sevendaystodie.AppId)

		case "remove":
			c := cmd.Client{}
			if app := c.GetApp(sevendaystodie.AppId); app != nil {
				app.Remove()
			}

		case "answer":
			installer.Answers <- r.Form.Get("answer")
			<-installer.Questions
		}

		w.Header().Set("Location", r.URL.String())
		w.WriteHeader(http.StatusSeeOther)
	}

	pageData := page{}
	if installer.Installing() {
		pageData.Status = "installing..."
		pageData.Refresh = true
	} else {
		c := cmd.Client{}
		if nil == c.GetApp(sevendaystodie.AppId) {
			pageData.Status = "not installed"
			pageData.Action = "install"
		} else {
			pageData.Status = "installed"
			pageData.Action = "remove"
		}
	}

	if len(installer.Questions) > 0 {
		pageData.Question = <-installer.Questions
		installer.Questions <- pageData.Question
		pageData.Refresh = false
	}

	t, err := template.New("steam.html").ParseFiles(assetsDir + "/steam.html")
	if err != nil {
		w.Write([]byte(fmt.Sprintln(err)))
	} else {

		if err := t.Execute(w, pageData); err != nil {
			w.Write([]byte(fmt.Sprintln(err)))
		}
	}
}
