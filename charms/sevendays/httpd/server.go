package httpd

import (
	"golang.org/x/net/context"
	"net/http"
)

type State struct {
	User     string
	Password string
}

var AssetsDir string

func Serve(ctx context.Context) error {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/steam", steamHandler)
	return http.ListenAndServe(":80", nil)
}
