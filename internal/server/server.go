package server

import (
	"io"
	"log"
	"net/http"
)

type Logger interface {
	Request(*http.Request, error)
}

type Client interface {
	userDetailsClient
}

type AuthenticateClient interface {
	Authenticate(http.ResponseWriter, *http.Request)
}

type Templates interface {
	ExecuteTemplate(io.Writer, string, interface{}) error
}

func New(logger *log.Logger, client Client, templates Templates, webDir string) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", http.RedirectHandler("/workflow", http.StatusFound))
	mux.Handle("/workflow", loggingInfoForWorflow(logger, client, templates))

	static := http.FileServer(http.Dir(webDir + "/static"))
	mux.Handle("/assets/", static)
	mux.Handle("/javascript/", static)
	mux.Handle("/stylesheets/", static)

	return mux
}
