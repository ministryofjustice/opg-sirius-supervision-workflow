package server

import (
	"github.com/ministryofjustice/opg-go-common/securityheaders"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/zap"
	"html/template"
	"io"
	"net/http"
	"net/url"
)

type Client interface {
	WorkflowVarsClient
	ClientTasksClient
}

type Template interface {
	Execute(wr io.Writer, data any) error
}

func New(logger *zap.Logger, client Client, templates map[string]*template.Template, prefix, siriusPublicURL, webDir string, defaultTeamId int) http.Handler {
	wrap := wrapHandler(client, logger, templates["error.gotmpl"], prefix, siriusPublicURL, defaultTeamId)

	mux := http.NewServeMux()

	mux.Handle("/", http.RedirectHandler(prefix+"/client-tasks", http.StatusFound))

	mux.Handle("/client-tasks",
		wrap(
			clientTasks(client, templates["client-tasks.gotmpl"])))

	mux.Handle("/health-check", healthCheck())

	static := http.FileServer(http.Dir(webDir + "/static"))
	mux.Handle("/assets/", static)
	mux.Handle("/javascript/", static)
	mux.Handle("/stylesheets/", static)

	return otelhttp.NewHandler(http.StripPrefix(prefix, securityheaders.Use(mux)), "supervision-workflow")
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
