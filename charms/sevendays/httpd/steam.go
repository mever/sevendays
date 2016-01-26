package httpd

import (
	"fmt"
	"github.com/mever/sevendaystodie"
	"github.com/mever/steam"
	"github.com/mever/steam/cmd"
	"net/http"
	"text/template"
)

var installer = &cmd.Installer{}

func init() {
	cmd.AddQuestion("Steam Guard code:", "What is your Steam Guard code?", true)
}

type SteamApp struct {
	AppId    steam.AppId
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
			installer.Install(steam.AppIdFromString(r.Form.Get("appId")))

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

	refresh := false
	app := SteamApp{AppId: sevendaystodie.AppId, Name: "7 days to die"}
	if installer.Installing() {
		app.Status = "installing..."
		refresh = true
	} else {
		c := cmd.Client{}
		if nil == c.GetApp(sevendaystodie.AppId) {
			app.Status = "not installed"
			app.Action = "install"
		} else {
			app.Status = "installed"
			app.Action = "remove"
		}
	}

	if len(installer.Questions) > 0 {
		app.Question = <-installer.Questions
		installer.Questions <- app.Question
		refresh = false
	}

	t, err := template.New("steam.html").ParseFiles(AssetsDir + "/steam.html")
	if err != nil {
		w.Write([]byte(fmt.Sprintln(err)))
	} else {

		data := struct {
			Refresh bool
			Apps    []SteamApp
		}{
			Refresh: refresh,
			Apps:    []SteamApp{app},
		}

		if err := t.Execute(w, data); err != nil {
			w.Write([]byte(fmt.Sprintln(err)))
		}
	}
}
