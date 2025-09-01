package server

import (
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/ministryofjustice/opg-go-common/securityheaders"
	"github.com/ministryofjustice/opg-go-common/telemetry"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
)

type ApiClient interface {
	WorkflowVarsClient
	ClientTasksClient
	CaseloadClient
	DeputiesClient
}

type Template interface {
	Execute(wr io.Writer, data any) error
}

func New(logger *slog.Logger, client ApiClient, templates map[string]*template.Template, envVars EnvironmentVars) http.Handler {
	wrap := wrapHandler(client, logger, templates["error.gotmpl"], envVars)

	mux := http.NewServeMux()

	mux.Handle("/", http.RedirectHandler(envVars.Prefix+"/client-tasks", http.StatusFound))

	mux.Handle("/client-tasks",
		wrap(
			clientTasks(client, templates["client-tasks.gotmpl"])))

	mux.Handle("/caseload",
		wrap(
			caseload(client, templates["caseload.gotmpl"])))

	mux.Handle("/deputy-tasks",
		wrap(
			deputyTasks(client, templates["deputy-tasks.gotmpl"])))

	mux.Handle("/deputies",
		wrap(
			deputies(client, templates["deputies.gotmpl"])))

	mux.Handle("/health-check", healthCheck())

	static := http.FileServer(http.Dir(envVars.WebDir + "/static"))
	mux.Handle("/assets/", static)
	mux.Handle("/javascript/", static)
	mux.Handle("/stylesheets/", static)

	return http.StripPrefix(envVars.Prefix, securityheaders.Use(telemetry.Middleware(logger)(mux)))
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
