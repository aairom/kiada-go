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
)

const appName = "kiada-go"
const version = "1.0"

func main() {
	listenPort := getEnv("LISTEN_PORT", "8080")
	addr := fmt.Sprintf(":%s", listenPort)

	logStartup(addr)

	srv := &http.Server{
		Addr:         addr,
		Handler:      newRouter(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Run server in goroutine so shutdown handling works
	go func() {
		log.Printf("%s v%s listening on %s", appName, version, addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	// Graceful shutdown on SIGTERM / SIGINT (matches Ch 6 SIGTERM handler pattern)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Printf("Received shutdown signal. Shutting down %s...", appName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
	log.Printf("%s shut down cleanly.", appName)
}

func logStartup(addr string) {
	hostname, _ := os.Hostname()
	log.Printf("-------------------------------------------")
	log.Printf("%s v%s — Kubernetes in Action Demo (Go)", appName, version)
	log.Printf("-------------------------------------------")
	log.Printf("Pod name   : %s", getEnv("POD_NAME", hostname))
	log.Printf("Pod IP     : %s", getEnv("POD_IP", "0.0.0.0"))
	log.Printf("Node name  : %s", getEnv("NODE_NAME", "unknown"))
	log.Printf("Node IP    : %s", getEnv("NODE_IP", "0.0.0.0"))
	log.Printf("QUOTE_URL  : %s", getEnv("QUOTE_URL", "(not set)"))
	log.Printf("QUIZ_URL   : %s", getEnv("QUIZ_URL", "(not set)"))
	log.Printf("Listen addr: %s", addr)
}
