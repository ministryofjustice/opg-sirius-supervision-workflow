package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type ErrorVars struct {
	Code  int
	Error string
	EnvironmentVars
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

type Handler func(app WorkflowVars, w http.ResponseWriter, r *http.Request) error

func wrapHandler(client ApiClient, logger *zap.SugaredLogger, tmplError Template, envVars EnvironmentVars) func(next Handler) http.Handler {
	return func(next Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			vars, err := NewWorkflowVars(client, r, envVars)
			if err == nil {
				err = next(*vars, w, r)
			}

			logger.Infow(
				"Application Request",
				"method", r.Method,
				"uri", r.URL.RequestURI(),
				"duration", time.Since(start),
			)

			if err != nil {
				if err == sirius.ErrUnauthorized {
					fmt.Println("unauthorised error ")
					http.Redirect(w, r, envVars.SiriusURL+"/auth", http.StatusFound)
					return
				}

				if redirect, ok := err.(RedirectError); ok {
					fmt.Println("redirect error ")

					http.Redirect(w, r, envVars.Prefix+"/"+redirect.To(), http.StatusFound)
					return
				}
				fmt.Println("error")
				fmt.Println(err.Error())

				logger.Errorw("Error handler", err)

				code := http.StatusInternalServerError
				if serverStatusError, ok := err.(StatusError); ok {
					code = serverStatusError.Code()
				}
				if siriusStatusError, ok := err.(sirius.StatusError); ok {
					code = siriusStatusError.Code
				}

				w.WriteHeader(code)
				errVars := ErrorVars{
					Code:            code,
					Error:           err.Error(),
					EnvironmentVars: envVars,
				}
				err = tmplError.Execute(w, errVars)

				if err != nil {
					logger.Errorw("Failed to render error template", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}
		})
	}
}
