package server

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gorilla/sessions"
	"log/slog"
	"net/http"
	"time"

	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
)

type ErrorVars struct {
	Code  int
	Error string
	EnvironmentVars
}

type Redirect struct {
	Path           string
	SuccessMessage string
}

func (e Redirect) Error() string {
	return "redirect to " + string(e.Path)
}

func (e Redirect) To() string {
	return string(e.Path)
}

func (e Redirect) GetHeaderMessage() string { return string(e.SuccessMessage) }

type StatusError int

func (e StatusError) Error() string {
	code := e.Code()

	return fmt.Sprintf("%d %s", code, http.StatusText(code))
}

func (e StatusError) Code() int {
	return int(e)
}

type Handler func(app WorkflowVars, w http.ResponseWriter, r *http.Request) error

func wrapHandler(client ApiClient, logger *slog.Logger, tmplError Template, envVars EnvironmentVars, cookieStore sessions.CookieStore) func(next Handler) http.Handler {
	return func(next Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			vars, err := NewWorkflowVars(client, r, envVars)
			if err == nil {
				err = next(*vars, w, r)
			}

			logger.Info(
				"Application Request",
				"method", r.Method,
				"uri", r.URL.RequestURI(),
				"duration", time.Since(start),
			)

			if err != nil {
				if errors.Is(err, context.Canceled) {
					w.WriteHeader(499)
					return
				}

				if err == sirius.ErrUnauthorized {
					http.Redirect(w, r, envVars.SiriusURL+"/auth", http.StatusFound)
					return
				}

				if redirect, ok := err.(Redirect); ok {
					fmt.Println("adding a new cookie")

					if redirect.SuccessMessage != "" {
						err := createSuccessMessage(r, w, cookieStore, redirect.SuccessMessage)
						if err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
							return
						}
					}

					fmt.Println("redirecting now")
					http.Redirect(w, r, envVars.Prefix+"/"+redirect.To(), http.StatusFound)
					return
				}

				logger.Error("Error handler", "error", err)

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
					logger.Error("Failed to render error template", "error", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}
		})
	}
}

func createSuccessMessage(r *http.Request, w http.ResponseWriter, cookieStore sessions.CookieStore, redirectSuccessMessage string) error {
	session, err := cookieStore.Get(r, "successMessageStore")
	if err != nil {
		return err
	}

	// Get the previous flashes, if any.
	if flashes := session.Flashes(); len(flashes) > 0 {
		// Use the flash values. This should remove the old flash if one exists unread.
		err = session.Save(r, w)
		if err != nil {
			return err
		}
	}

	// Set a new flash.
	encodedContent := base64.StdEncoding.EncodeToString([]byte(redirectSuccessMessage))
	session.AddFlash(encodedContent)

	err = session.Save(r, w)
	if err != nil {
		return err
	}
	return nil
}
