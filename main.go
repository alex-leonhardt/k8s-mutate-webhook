package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/loadsmart/k8s-mutate-webhook/pkg/web"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	wait             time.Duration
	originalRegistry string
	newRegistry      string
	certFile string
	keyFile string
)

func init() {
	originalRegistry = os.Getenv("ORIGINAL_REGISTRY")
	newRegistry = os.Getenv("NEW_REGISTRY")

	flag.StringVar(&originalRegistry, "original-registry", originalRegistry, "registry to be replaced by new registry")
	flag.StringVar(&newRegistry, "new-registry", newRegistry, "registry to replace the original registry")
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.StringVar(&certFile, "cert-file", certFile, "TLS certificate file")
	flag.StringVar(&keyFile, "key-file", keyFile, "TLS private key")
	flag.Parse()
}

func initializeLog() (*zap.Logger, error) {
	c := zap.NewProductionConfig()
	c.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(time.RFC3339))
	}
	return c.Build()
}

func main() {
	// Initialize Zap Logger
	logger, err := initializeLog()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync() // nolint

	if originalRegistry == "" || newRegistry == "" {
		log.Fatal("original registry and new registry are required to start")
	}

	logger.Info("Starting server ...")

	r := mux.NewRouter()
	// Configuring handlers and routes
	r.HandleFunc("/healthz", web.HandleHealthz)
	web.NewMutatePodsHandler(r, originalRegistry, newRegistry, logger)

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
		ErrorLog:     zap.NewStdLog(logger),
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if certFile == "" || keyFile == "" {
			if err := srv.ListenAndServe(); err != nil {
				log.Println(err)
			}
		} else {
			if err := srv.ListenAndServeTLS(certFile, keyFile); err != nil {
				log.Println(err)
			}
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C) or SIGTERM (Ctrl+/)
	// SIGKILL or SIGQUIT will not be caught.
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx) // nolint
	logger.Info("shutting down")
	os.Exit(0)
}
