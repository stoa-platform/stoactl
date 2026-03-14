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

	"github.com/stoa-platform/stoa-go/internal/connect"
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
		stoaCtx, ctxErr := cfg.GetCurrentContext()
		if ctxErr == nil && stoaCtx != nil {
			log.Printf("using context %q (server: %s)", cfg.CurrentContext, stoaCtx.Context.Server)
		}
	}

	port := os.Getenv("STOA_CONNECT_PORT")
	if port == "" {
		port = "8090"
	}

	// Set up CP registration agent
	agent := connect.New(connect.ConfigFromEnv(Version))

	// Root context — cancelled on SIGINT/SIGTERM
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Register with Control Plane (if configured)
	if agent.IsConfigured() {
		if err := agent.Register(ctx, port); err != nil {
			log.Printf("warning: CP registration failed: %v", err)
		} else {
			agent.StartHeartbeat(ctx)
		}
	} else {
		log.Println("CP registration skipped (STOA_CONTROL_PLANE_URL or STOA_GATEWAY_API_KEY not set)")
	}

	// Health endpoint
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		gatewayID := agent.GatewayID()
		if gatewayID == "" {
			gatewayID = "unregistered"
		}
		fmt.Fprintf(w, `{"status":"ok","version":"%s","commit":"%s","gateway_id":"%s"}`,
			Version, Commit, gatewayID)
	})

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("stoa-connect %s (%s) listening on :%s", Version, Commit, port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down...")
	cancel() // Stop heartbeat goroutine

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("shutdown: %v", err)
	}
}
