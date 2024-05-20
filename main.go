package main

import (
	"context"
	"github.com/ministryofjustice/opg-go-common/env"
	"github.com/ministryofjustice/opg-go-common/telemetry"
	"github.com/ministryofjustice/opg-go-common/paginate"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/util"
	"go.opentelemetry.io/contrib/detectors/aws/ecs"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/ministryofjustice/opg-sirius-workflow/internal/server"
	"github.com/ministryofjustice/opg-sirius-workflow/internal/sirius"
)

func initTracerProvider(ctx context.Context, logger *zap.SugaredLogger) func() {
	resource, err := ecs.NewResourceDetector().Detect(ctx)
	if err != nil {
		logger.Fatal(err)
	}

	traceExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint("0.0.0.0:4317"),
	)
	if err != nil {
		logger.Fatal(err)
	}

	idg := xray.NewIDGenerator()
	tp := trace.NewTracerProvider(
		trace.WithResource(resource),
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(traceExporter),
		trace.WithIDGenerator(idg),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(xray.Propagator{})

	return func() {
		if err := tp.Shutdown(ctx); err != nil {
			logger.Fatal(err)
		}
	}
}

func main() {
	logger := zap.Must(zap.NewProduction(zap.Fields(zap.String("service_name", "opg-sirius-workflow")))).Sugar()
	apiCallLogger := telemetry.NewLogger("opg-sirius-workflow")

	defer func() { _ = logger.Sync() }()

	if env.Get("TRACING_ENABLED", "0") == "1" {
		shutdown := initTracerProvider(context.Background(), logger)
		defer shutdown()
	}

	httpClient := http.DefaultClient
	httpClient.Transport = otelhttp.NewTransport(httpClient.Transport)

	envVars, err := server.NewEnvironmentVars()
	if err != nil {
		logger.Fatalw("Error creating EnvironmentVars", "error", err)
	}

	client, err := sirius.NewApiClient(http.DefaultClient, envVars.SiriusURL, apiCallLogger)
	if err != nil {
		logger.Fatalw("Error returned by Sirius New ApiClient", "error", err)
	}

	templates := createTemplates(envVars)

	server := &http.Server{
		Addr:    ":" + envVars.Port,
		Handler: server.New(logger, client, templates, envVars),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Infow("Error returned by server.ListenAndServe()",
				"error", err,
			)
			logger.Fatal(err)
		}
	}()

	logger.Infow("Running at :" + envVars.Port)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	sig := <-c
	logger.Infow("signal received: ", sig)

	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(tc); err != nil {
		logger.Infow("Error returned by server.Shutdown",
			"error", err,
		)
	}
}

func createTemplates(envVars server.EnvironmentVars) map[string]*template.Template {
	templates := map[string]*template.Template{}
	templateFunctions := map[string]interface{}{
		"contains": func(xs []string, needle string) bool {
			for _, x := range xs {
				if x == needle {
					return true
				}
			}

			return false
		},
		"prefix": func(s string) string {
			return envVars.Prefix + s
		},
		"sirius": func(s string) string {
			return envVars.SiriusPublicURL + s
		},
		"is_last": util.IsLast,
	}

	templateDirPath := envVars.WebDir + "/template"
	templateDir, _ := os.Open(templateDirPath)
	templateDirs, _ := templateDir.Readdir(0)
	_ = templateDir.Close()

	mainTemplates, _ := filepath.Glob(templateDirPath + "/*.gotmpl")

	for _, file := range mainTemplates {
		tmpl := template.New(filepath.Base(file)).Funcs(templateFunctions)
		for _, dir := range templateDirs {
			if dir.IsDir() {
				tmpl, _ = tmpl.ParseGlob(templateDirPath + "/" + dir.Name() + "/*.gotmpl")
			}
		}
		tmpl, _ = tmpl.Parse(paginate.Template)
		templates[tmpl.Name()] = template.Must(tmpl.ParseFiles(file))
	}

	return templates
}
