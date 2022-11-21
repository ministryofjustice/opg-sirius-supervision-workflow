package main

import (
	"context"
	"github.com/ministryofjustice/opg-go-common/logging"
	"go.uber.org/zap"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/ministryofjustice/opg-sirius-workflow/internal/server"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
)

func main() {
	serverLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	defer serverLogger.Sync()

	port := getEnv("PORT", "1234")
	webDir := getEnv("WEB_DIR", "web")
	siriusURL := getEnv("SIRIUS_URL", "http://localhost:8080")
	siriusPublicURL := getEnv("SIRIUS_PUBLIC_URL", "")
	DefaultWorkflowTeam := getEnv("DEFAULT_WORKFLOW_TEAM", "21")
	prefix := getEnv("PREFIX", "")

	layouts, _ := template.
		New("").
		Funcs(map[string]interface{}{
			"join": func(sep string, items []string) string {
				return strings.Join(items, sep)
			},
			"contains": func(xs []string, needle string) bool {
				for _, x := range xs {
					if x == needle {
						return true
					}
				}

				return false
			},
			"prefix": func(s string) string {
				return prefix + s
			},
			"sirius": func(s string) string {
				return siriusPublicURL + s
			},
		}).
		ParseGlob(webDir + "/template/layout/*.gotmpl")

	files, _ := filepath.Glob(webDir + "/template/*.gotmpl")
	tmpls := map[string]*template.Template{}

	for _, file := range files {
		tmpls[filepath.Base(file)] = template.Must(template.Must(layouts.Clone()).ParseFiles(file))
	}
	apiCallLogger := logging.New(os.Stdout, "opg-sirius-workflow ")

	client, err := sirius.NewClient(http.DefaultClient, siriusURL, apiCallLogger)
	sugar := serverLogger.Sugar()
	if err != nil {
		sugar.Infow("Error returned by Sirius New Client",
			"error", err,
		)
	}

	defaultWorkflowTeam, err := strconv.Atoi(DefaultWorkflowTeam)
	if err != nil {
		sugar.Infow("Error converting DEFAULT_WORKFLOW_TEAM to int")
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: server.New(serverLogger, client, tmpls, prefix, siriusPublicURL, webDir, defaultWorkflowTeam),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			sugar.Infow("Error returned by server.ListenAndServe()",
				"error", err,
			)
			sugar.Fatal(err)
		}
	}()

	sugar.Infow("Running at :" + port)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	sig := <-c
	sugar.Infow("signal received: ", sig)

	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(tc); err != nil {
		sugar.Infow("Error returned by server.Shutdown",
			"error", err,
		)
	}
}

func getEnv(key, def string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return def
}
