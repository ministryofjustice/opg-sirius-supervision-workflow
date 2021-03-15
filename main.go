package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ministryofjustice/opg-sirius-workflow/internal/server"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
)

func main() {
	logger := log.New(os.Stdout, "opg-sirius-workflow ", log.LstdFlags)

	port := getEnv("PORT", "9001")
	webDir := getEnv("WEB_DIR", "web")
	siriusURL := getEnv("SIRIUS_URL", "http://localhost:8080")

	templates, err := template.New("").Funcs(map[string]interface{}{
		"join": func(sep string, items []string) string {
			return strings.Join(items, sep)
		},
	}).ParseGlob(webDir + "/template/*.gotmpl")
	if err != nil {
		logger.Fatalln(err)
	}

	client, err := sirius.NewClient(http.DefaultClient, siriusURL)
	if err != nil {
		logger.Fatalln(err)
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: server.New(logger, client, templates, webDir),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Fatalln(err)
		}
	}()

	logger.Println("Running at :" + port)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	sig := <-c
	logger.Println("Received terminate, graceful shutdown", sig)

	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(tc); err != nil {
		logger.Println(err)
	}
}

func getEnv(key, def string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return def
}
