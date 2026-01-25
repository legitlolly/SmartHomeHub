package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/legitlolly/SmartHomeHub/internal/api"
	"github.com/legitlolly/SmartHomeHub/internal/device"
	"github.com/legitlolly/SmartHomeHub/internal/providers/hue"
	"github.com/legitlolly/SmartHomeHub/internal/providers/simulator"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	registry := device.NewRegistry()

	//temp test device
	tempDevice := simulator.NewSimulatedDevice("temp-light-1")
	if err := registry.Register(tempDevice); err != nil {
		log.Printf("Failed to register temp device: %v", err)
	}
	log.Println("Registered temp device: temp-light-1")

	// Load Hue configuration from environment
	hueIP := os.Getenv("HUE_BRIDGE_IP")
	hueUsername := os.Getenv("HUE_USERNAME")

	if hueIP != "" && hueUsername != "" {
		log.Println("Hue configuration detected, discovering Hue devices...")

		// Discover and register all lights on the bridge
		if err := hue.DiscoverAndRegisterLights(ctx, registry, hueIP, hueUsername); err != nil {
			log.Printf("Failed to discover Hue devices: %v", err)
		}
	} else {
		log.Println("Hue configuration not found (HUE_BRIDGE_IP and HUE_USERNAME not set)")
	}

	handler := api.NewHandler(registry)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("HTTP server listening on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()
	log.Println("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}
