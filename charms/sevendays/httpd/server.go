package httpd

import (
	"net/http"
	"golang.org/x/net/context"
)

type State struct{
	User string
	Password string
}

var AssetsDir string

func Serve(ctx context.Context) error {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/steam", steamHandler)
	return http.ListenAndServe(":80", nil)
}
