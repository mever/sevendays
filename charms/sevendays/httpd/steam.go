package httpd

import (
	"fmt"
	"github.com/juju/errors"
	"github.com/mever/sevendaystodie"
	"github.com/mever/steam"
	"github.com/mever/steam/cmd"
	"net/http"
	"sync"
	"text/template"
)

var (
	ErrAlreadyInstalling = errors.New("We're aleady installing")
)

func init() {
	cmd.AddQuestion("Steam Guard code:", "What is your Steam Guard code?", true)
}

type question struct {
	Sensitive bool
	Value     string
}

type Installer struct {
	mu          sync.Mutex
	installing  bool
	interviewer cmd.Interviewer

	Questions chan *question
	Answers   chan string
}

func (i *Installer) Installing() bool {
	i.mu.Lock()
	defer i.mu.Unlock()
	return i.installing
}

func (i *Installer) Install(appId steam.AppId) error {
	i.mu.Lock()
	installing := i.installing
	if !installing {
		i.installing = true
	}
	i.mu.Unlock()
	if installing {
		return ErrAlreadyInstalling
	}

	i.Questions = make(chan *question, 1)
	i.Answers = make(chan string)
	i.interviewer = func(q string, sensitive bool) string {
		if q == "" {
			fmt.Println("No further questions...")
			close(i.Questions)
			close(i.Answers)
			i.mu.Lock()
			i.installing = false
			i.mu.Unlock()
			return ""
		} else {
			i.Questions <- &question{Value: q, Sensitive: sensitive}
		}
		return <-i.Answers
	}

	c := cmd.Client{}
	go func() {

		// TODO: move into Steam CMD, not all apps require authentication
		c.AuthUser = i.interviewer("What is your Steam username?", false)
		c.AuthPw = i.interviewer("What is your Steam password?", true)

		err := c.InstallApp(appId, i.interviewer)
		if err != nil {
			fmt.Println(err)
		}
	}()

	return nil
}

func NewInstaller() *Installer {
	i := Installer{}
	return &i
}

var defaultInstaller = NewInstaller()

type SteamApp struct {
	AppId    steam.AppId
	Name     string
	Status   string
	Action   string
	Question *question
}

func steamHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		switch r.Form.Get("action") {
		case "install":
			defaultInstaller.Install(steam.AppIdFromString(r.Form.Get("appId")))

		case "remove":
			c := cmd.Client{}
			if app := c.GetApp(sevendaystodie.AppId); app != nil {
				app.Remove()
			}

		case "answer":
			defaultInstaller.Answers <- r.Form.Get("answer")
			<-defaultInstaller.Questions
		}

		w.Header().Set("Location", r.URL.String())
		w.WriteHeader(http.StatusSeeOther)
	}

	refresh := false
	app := SteamApp{AppId: sevendaystodie.AppId, Name: "7 days to die"}
	if defaultInstaller.Installing() {
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

	if len(defaultInstaller.Questions) > 0 {
		app.Question = <-defaultInstaller.Questions
		defaultInstaller.Questions <- app.Question
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
