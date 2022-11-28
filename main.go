package main

import (
	"context"
	"github.com/ministryofjustice/opg-go-common/env"
	"github.com/ministryofjustice/opg-go-common/logging"
	"go.opentelemetry.io/contrib/detectors/aws/ecs"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"html/template"
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

func initTracerProvider(ctx context.Context, logger *zap.Logger) func() {
	resource, err := ecs.NewResourceDetector().Detect(ctx)
	sugar := logger.Sugar()
	if err != nil {
		sugar.Fatal(err)
	}

	traceExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint("0.0.0.0:4317"),
		otlptracegrpc.WithDialOption(grpc.WithBlock()),
	)
	if err != nil {
		sugar.Fatal(err)
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
			sugar.Fatal(err)
		}
	}
}

func main() {
	serverLogger, err := zap.NewProduction()
	sugar := serverLogger.Sugar()

	if err != nil {
		sugar.Infow("Error creating logger: %v\n", err)
	}

	if err := serverLogger.Sync(); err != nil {
		sugar.Infow("Error syncing logger: %v\n", err)
	}

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

	if env.Get("TRACING_ENABLED", "0") == "1" {
		shutdown := initTracerProvider(context.Background(), serverLogger)
		defer shutdown()
	}

	httpClient := http.DefaultClient
	httpClient.Transport = otelhttp.NewTransport(httpClient.Transport)

	client, err := sirius.NewClient(http.DefaultClient, siriusURL, apiCallLogger)
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
