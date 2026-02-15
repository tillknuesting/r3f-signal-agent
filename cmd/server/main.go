package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpcollector "r3f-trends/internal/adapter/driven/collector/http"
	"r3f-trends/internal/adapter/driven/config/yaml"
	"r3f-trends/internal/adapter/driven/storage/markdown"
	"r3f-trends/internal/app/service"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config"
	}

	cfgLoader := yaml.NewConfigLoader(configPath)
	cfg, err := cfgLoader.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	trendRepo := markdown.NewTrendRepositoryAdapter(cfg.Storage.BasePath)
	httpCollector := httpcollector.New()

	collectorSvc := service.NewCollectorService(
		trendRepo,
		map[string]interface{}{
			"http": httpCollector,
		},
	)

	trendSvc := service.NewTrendService(trendRepo)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	mux.HandleFunc("/api/v1/trends", func(w http.ResponseWriter, r *http.Request) {
		trends, total, err := trendSvc.List(r.Context(), service.ListOptions{
			Limit:  100,
			Offset: 0,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var dtos []interface{}
		for _, t := range trends {
			dtos = append(dtos, t.ToDTO())
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"trends": dtos,
			"total":  total,
		})
	})

	mux.HandleFunc("/api/v1/collect", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		sourceLoader := yaml.NewSourceLoader(configPath + "/sources")
		sources, err := sourceLoader.LoadAll(r.Context())
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to load sources: %v", err), http.StatusInternalServerError)
			return
		}

		var sourceIDs []string
		for _, s := range sources {
			if s.Enabled() {
				sourceIDs = append(sourceIDs, s.ID())
			}
		}

		result, err := collectorSvc.Collect(r.Context(), cfg.ActiveProfile, sourceIDs, sources)
		if err != nil {
			http.Error(w, fmt.Sprintf("Collection failed: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":      "completed",
			"items_count": len(result.Trends),
			"errors":      result.Errors,
		})
	})

	mux.HandleFunc("/api/v1/sources", func(w http.ResponseWriter, r *http.Request) {
		sourceLoader := yaml.NewSourceLoader(configPath + "/sources")
		sources, err := sourceLoader.LoadAll(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var dtos []interface{}
		for _, s := range sources {
			dtos = append(dtos, s.ToDTO())
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"sources": dtos,
		})
	})

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		log.Printf("Starting server on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	srv.Shutdown(ctx)
	log.Println("Server stopped")
}
