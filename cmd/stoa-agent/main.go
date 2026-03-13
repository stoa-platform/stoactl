package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/stoa-platform/stoa-go/pkg/config"
)

// Version and Commit are set via ldflags at build time.
var (
	Version = "dev"
	Commit  = "none"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Printf("warning: could not load config: %v", err)
	} else {
		ctx, ctxErr := cfg.GetCurrentContext()
		if ctxErr == nil && ctx != nil {
			log.Printf("using context %q (server: %s)", cfg.CurrentContext, ctx.Context.Server)
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok","version":"%s","commit":"%s"}`, Version, Commit)
	})

	port := os.Getenv("STOA_AGENT_PORT")
	if port == "" {
		port = "8090"
	}

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("stoa-agent %s (%s) listening on :%s", Version, Commit, port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown: %v", err)
	}
}
