package server

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"

	"github.com/ministryofjustice/opg-go-common/securityheaders"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
)

type Logger interface {
	Request(*http.Request, error)
}

type Client interface {
	WorkflowInformation
}

type Template interface {
	ExecuteTemplate(io.Writer, string, interface{}) error
}

func New(logger Logger, client Client, templates map[string]*template.Template, prefix, siriusPublicURL, webDir string, defaultWorkflowTeam int) http.Handler {
	wrap := errorHandler(logger, templates["error.gotmpl"], prefix, siriusPublicURL)

	mux := http.NewServeMux()
	mux.Handle("/",
		wrap(
			loggingInfoForWorkflow(client, templates["workflow.gotmpl"], defaultWorkflowTeam)))

	mux.Handle("/health-check", healthCheck())

	static := http.FileServer(http.Dir(webDir + "/static"))
	mux.Handle("/assets/", static)
	mux.Handle("/javascript/", static)
	mux.Handle("/stylesheets/", static)

	return http.StripPrefix(prefix, securityheaders.Use(mux))
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

type Handler func(w http.ResponseWriter, r *http.Request) error

type errorVars struct {
	Firstname string
	Surname   string
	SiriusURL string
	Path      string
	Code      int
	Error     string
}

func errorHandler(logger Logger, tmplError Template, prefix, siriusURL string) func(next Handler) http.Handler {
	return func(next Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			err := next(w, r)

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
					Firstname: "",
					Surname:   "",
					SiriusURL: siriusURL,
					Path:      "",
					Code:      code,
					Error:     err.Error(),
				})

				if err != nil {
					logger.Request(r, err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
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
