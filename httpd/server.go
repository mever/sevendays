package httpd

import (
	"github.com/mever/sevendays/httpd/steam"
	"golang.org/x/net/context"
	"net/http"
)

type State struct {
	User     string
	Password string
}

var AssetsDir string

func Serve(ctx context.Context) error {
	steam.Setup(AssetsDir)
	http.HandleFunc("/", indexHandler)
	return http.ListenAndServe(":80", nil)
}
