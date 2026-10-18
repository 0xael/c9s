package web

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func registerHandlers(mux *http.ServeMux, cache *dataCache) {
	mux.HandleFunc("/api/overview", jsonHandler(func() any { return handleOverview(cache) }))
	mux.HandleFunc("/api/activity", jsonHandler(func() any { return handleActivity(cache) }))
	mux.HandleFunc("/api/tokens", jsonHandler(func() any { return handleTokens(cache) }))
	mux.HandleFunc("/api/cost", jsonHandler(func() any { return handleCost(cache) }))
	mux.HandleFunc("/api/models", jsonHandler(func() any { return handleModels(cache) }))
	mux.HandleFunc("/api/hours", jsonHandler(func() any { return handleHours(cache) }))
	mux.HandleFunc("/api/sessions", jsonHandler(func() any { return handleSessions(cache) }))
	mux.HandleFunc("/api/history", jsonHandler(func() any { return handleHistory(cache) }))
	mux.HandleFunc("/api/session/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/api/session/")
		if id == "" {
			http.Error(w, "missing session id", http.StatusBadRequest)
			return
		}
		writeJSON(w, handleSessionDetail(cache, id))
	})

	sub, _ := fs.Sub(staticFS, "static")
	mux.Handle("/", http.FileServer(http.FS(sub)))
}

func jsonHandler(fn func() any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, fn())
	}
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func handleOverview(c *dataCache) OverviewResponse {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var totalCost float64
	var totalTokens int64
	var linesAdded, linesRemoved int
	for _, sc := range c.costs {
		totalCost += sc.TotalCostUSD
		linesAdded += sc.TotalLinesAdded
		linesRemoved += sc.TotalLinesRemoved
		for _, mu := range sc.ModelUsage {
			totalTokens += mu.InputTokens + mu.OutputTokens + mu.CacheReadInputTokens + mu.CacheCreationInputTokens
		}
	}

	return OverviewResponse{
		TotalSessions:    c.stats.TotalSessions,
		TotalMessages:    c.stats.TotalMessages,
		TotalCostUSD:     totalCost,
		TotalTokens:      totalTokens,
		FirstSessionDate: c.stats.FirstSessionDate,
		LiveSessions:     len(c.sessions),
		TotalLinesAdded:  linesAdded,
		TotalLinesRemoved: linesRemoved,
	}
}

func handleActivity(c *dataCache) any {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.stats.DailyActivity
}

func handleTokens(c *dataCache) any {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.stats.DailyModelTokens
}

func handleCost(c *dataCache) CostResponse {
	c.mu.RLock()
	defer c.mu.RUnlock()

	dailyMap := make(map[string]map[string]float64)
	totalByModel := make(map[string]float64)

	for _, sc := range c.costs {
		date := time.UnixMilli(sc.StartTime).Format("2006-01-02")
		if dailyMap[date] == nil {
			dailyMap[date] = make(map[string]float64)
		}
		for model, mu := range sc.ModelUsage {
			dailyMap[date][model] += mu.CostUSD
			totalByModel[model] += mu.CostUSD
		}
	}

	var dates []string
	for d := range dailyMap {
		dates = append(dates, d)
	}
	sort.Strings(dates)

	daily := make([]DailyCost, 0, len(dates))
	for _, d := range dates {
		daily = append(daily, DailyCost{Date: d, CostByModel: dailyMap[d]})
	}

	return CostResponse{Daily: daily, TotalByModel: totalByModel}
}

func handleModels(c *dataCache) []ModelInfo {
	c.mu.RLock()
	defer c.mu.RUnlock()

	costByModel := make(map[string]float64)
	for _, sc := range c.costs {
		for model, mu := range sc.ModelUsage {
			costByModel[model] += mu.CostUSD
		}
	}

	var models []ModelInfo
	for name, mt := range c.stats.ModelUsage {
		total := mt.Total()
		models = append(models, ModelInfo{
			Name:         name,
			InputTokens:  mt.InputTokens,
			OutputTokens: mt.OutputTokens,
			CacheRead:    mt.CacheReadInputTokens,
			CacheWrite:   mt.CacheCreationInputTokens,
			TotalTokens:  total,
			CostUSD:      costByModel[name],
		})
	}
	sort.Slice(models, func(i, j int) bool { return models[i].TotalTokens > models[j].TotalTokens })
	return models
}

func handleHours(c *dataCache) any {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.stats.HourCounts
}

func handleSessions(c *dataCache) []WebSession {
	c.mu.RLock()
	defer c.mu.RUnlock()

	out := make([]WebSession, len(c.sessions))
	for i, s := range c.sessions {
		out[i] = toWebSession(s)
	}
	return out
}

func handleHistory(c *dataCache) []HistoryRow {
	c.mu.RLock()
	defer c.mu.RUnlock()

	costIdx := make(map[string]*SessionCostSummary, len(c.costs))
	for _, sc := range c.costs {
		var totalTok int64
		var models []string
		for model, mu := range sc.ModelUsage {
			totalTok += mu.InputTokens + mu.OutputTokens + mu.CacheReadInputTokens + mu.CacheCreationInputTokens
			models = append(models, model)
		}
		sort.Strings(models)
		costIdx[sc.SessionID] = &SessionCostSummary{
			CostUSD:      sc.TotalCostUSD,
			TotalTokens:  totalTok,
			DurationMs:   sc.TotalDuration,
			LinesAdded:   sc.TotalLinesAdded,
			LinesRemoved: sc.TotalLinesRemoved,
			Models:       models,
		}
	}

	rows := make([]HistoryRow, 0, len(c.history))
	for _, h := range c.history {
		row := HistoryRow{
			SessionID:    h.SessionID,
			Project:      h.Project,
			ProjectShort: filepath.Base(h.Project),
			FirstPrompt:  h.FirstPrompt,
			PromptCount:  h.PromptCount,
			FirstAt:      h.FirstAt.UnixMilli(),
			LastAt:       h.LastAt.UnixMilli(),
		}
		if cs, ok := costIdx[h.SessionID]; ok {
			row.CostUSD = cs.CostUSD
			row.TotalTokens = cs.TotalTokens
			row.DurationMs = cs.DurationMs
			row.LinesAdded = cs.LinesAdded
			row.LinesRemoved = cs.LinesRemoved
			row.Models = cs.Models
		}
		rows = append(rows, row)
	}
	return rows
}

type SessionCostSummary struct {
	CostUSD      float64
	TotalTokens  int64
	DurationMs   int64
	LinesAdded   int
	LinesRemoved int
	Models       []string
}

func handleSessionDetail(c *dataCache, sessionID string) SessionDetail {
	c.mu.RLock()
	defer c.mu.RUnlock()

	detail := SessionDetail{SessionID: sessionID}

	for _, sc := range c.costs {
		if sc.SessionID == sessionID {
			cost := sc
			detail.Cost = &cost
			detail.Project = sc.Project
			break
		}
	}

	for _, h := range c.history {
		if h.SessionID == sessionID {
			detail.Project = h.Project
			detail.Prompts = h.Prompts
			break
		}
	}

	return detail
}
