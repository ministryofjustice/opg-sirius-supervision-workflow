package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type Client interface {
	myDetailsClient
}

type AuthenticateClient interface {
	Authenticate(http.ResponseWriter, *http.Request)
}

type Templates interface {
	ExecuteTemplate(io.Writer, string, interface{}) error
}

func New(logger *log.Logger, client Client, templates Templates, webDir string) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", http.RedirectHandler("/my-details", http.StatusFound))
	mux.Handle("/my-details", myDetails(logger, client, templates))

	static := http.FileServer(http.Dir(webDir + "/static"))
	mux.Handle("/assets/", static)
	mux.Handle("/javascript/", static)
	mux.Handle("/stylesheets/", static)

	return mux
}

type RedirectError string

func (e RedirectError) Error() string {
	return "redirect to " + string(e)
}

func (e RedirectError) To() string {
	return string(e)
}

type StatusError int

func (e StatusError) Error() string {
	code := e.Code()

	return fmt.Sprintf("%d %s", code, http.StatusText(code))
}

func (e StatusError) Code() int {
	return int(e)
}

type Handler func(perm sirius.PermissionSet, w http.ResponseWriter, r *http.Request) error

type errorVars struct {
	SiriusURL string
	Path      string

	Code  int
	Error string
}

type ErrorHandlerClient interface {
	MyPermissions(sirius.Context) (sirius.PermissionSet, error)
}

func errorHandler(logger Logger, client ErrorHandlerClient, tmplError Template, prefix, siriusURL string) func(next Handler) http.Handler {
	return func(next Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			myPermissions, err := client.MyPermissions(getContext(r))

			if err == nil {
				err = next(myPermissions, w, r)
			}

			if err != nil {
				if err == sirius.ErrUnauthorized {
					http.Redirect(w, r, siriusURL+"/auth", http.StatusFound)
					return
				}

				if redirect, ok := err.(RedirectError); ok {
					http.Redirect(w, r, prefix+redirect.To(), http.StatusFound)
					return
				}

				logger.Request(r, err)

				code := http.StatusInternalServerError
				if status, ok := err.(StatusError); ok {
					if status.Code() == http.StatusForbidden || status.Code() == http.StatusNotFound {
						code = status.Code()
					}
				}

				w.WriteHeader(code)
				err = tmplError.ExecuteTemplate(w, "page", errorVars{
					SiriusURL: siriusURL,
					Path:      "",
					Code:      code,
					Error:     err.Error(),
				})

				if err != nil {
					logger.Request(r, err)
					http.Error(w, "Could not generate error template", http.StatusInternalServerError)
				}
			}
		})
	}
}

func getContext(r *http.Request) sirius.Context {
	token := ""

	if r.Method == http.MethodGet {
		if cookie, err := r.Cookie("XSRF-TOKEN"); err == nil {
			token, _ = url.QueryUnescape(cookie.Value)
		}
	} else {
		token = r.FormValue("xsrfToken")
	}

	return sirius.Context{
		Context:   r.Context(),
		Cookies:   r.Cookies(),
		XSRFToken: token,
	}
}
