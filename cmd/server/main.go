package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	glm5agent "r3f-trends/internal/adapter/driven/agent/glm5"
	chromecollector "r3f-trends/internal/adapter/driven/collector/chrome"
	httpcollector "r3f-trends/internal/adapter/driven/collector/http"
	"r3f-trends/internal/adapter/driven/config/yaml"
	"r3f-trends/internal/adapter/driven/storage/markdown"
	"r3f-trends/internal/app/service"
	"r3f-trends/internal/domain/entity"
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
	chromeCollector := chromecollector.New()

	collectorSvc := service.NewCollectorService(
		trendRepo,
		map[string]interface{}{
			"http":   httpCollector,
			"chrome": chromeCollector,
		},
	)

	trendSvc := service.NewTrendService(trendRepo)

	var agentSvc *service.AgentService
	if cfg.LLM.APIKey != "" {
		llmAgent := glm5agent.NewAgent(cfg.LLM.APIKey, cfg.LLM.BaseURL)
		agentSvc = service.NewAgentService(llmAgent, trendSvc)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	mux.HandleFunc("/api/v1/trends", trendsHandler(trendSvc))
	mux.HandleFunc("/api/v1/trends/", trendDetailHandler(trendSvc))
	mux.HandleFunc("/api/v1/collect", collectHandler(collectorSvc, configPath, cfg.ActiveProfile))
	mux.HandleFunc("/api/v1/sources", sourcesHandler(configPath))
	mux.HandleFunc("/api/v1/profiles", profilesHandler(configPath))

	if agentSvc != nil {
		mux.HandleFunc("/api/v1/agent/summarize", agentSummarizeHandler(agentSvc))
		mux.HandleFunc("/api/v1/agent/suggest", agentSuggestHandler(agentSvc))
	}

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	go func() {
		log.Printf("Starting server on %s", addr)
		log.Printf("Chrome collector: enabled")
		if agentSvc != nil {
			log.Printf("GLM-5 agent: enabled")
		}
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

func trendsHandler(trendSvc *service.TrendService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		date := r.URL.Query().Get("date")

		var trends []*entity.Trend
		var total int
		var err error

		if date != "" {
			trends, total, err = trendSvc.GetByDate(r.Context(), date, service.ListOptions{
				Limit:  100,
				Offset: 0,
			})
		} else {
			trends, total, err = trendSvc.List(r.Context(), service.ListOptions{
				Limit:  100,
				Offset: 0,
			})
		}

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
	}
}

func trendDetailHandler(trendSvc *service.TrendService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/api/v1/trends/"):]

		if strings.HasSuffix(id, "/star") {
			trendID := id[:len(id)-5]
			if r.Method != http.MethodPost {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}
			if err := trendSvc.Star(r.Context(), trendID); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"status": "starred"})
			return
		}

		if strings.HasSuffix(id, "/unstar") {
			trendID := id[:len(id)-7]
			if r.Method != http.MethodPost {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}
			if err := trendSvc.Unstar(r.Context(), trendID); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"status": "unstarred"})
			return
		}

		trend, err := trendSvc.Get(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(trend.ToDTO())
	}
}

func collectHandler(collectorSvc *service.CollectorService, configPath, activeProfile string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		result, err := collectorSvc.Collect(r.Context(), activeProfile, sourceIDs, sources)
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
	}
}

func sourcesHandler(configPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}

func profilesHandler(configPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		profileLoader := yaml.NewProfileLoader(configPath + "/profiles")
		profiles, err := profileLoader.LoadAll(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var dtos []interface{}
		for _, p := range profiles {
			dtos = append(dtos, p.ToDTO())
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"profiles": dtos,
		})
	}
}

func agentSummarizeHandler(agentSvc *service.AgentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			TrendID string `json:"trend_id"`
			Content string `json:"content"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		var summary string
		var err error

		if req.TrendID != "" {
			summary, err = agentSvc.Summarize(r.Context(), req.TrendID)
		} else if req.Content != "" {
			summary, err = agentSvc.SummarizeContent(r.Context(), req.Content)
		} else {
			http.Error(w, "trend_id or content required", http.StatusBadRequest)
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"summary": summary})
	}
}

func agentSuggestHandler(agentSvc *service.AgentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		profile := r.URL.Query().Get("profile")
		if profile == "" {
			profile = "tech"
		}

		suggestions, err := agentSvc.SuggestTopics(r.Context(), profile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"suggestions": suggestions,
		})
	}
}
