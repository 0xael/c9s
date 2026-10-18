package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/0xael/c9s/internal/config"
	"github.com/0xael/c9s/internal/data"
)

type dataCache struct {
	mu         sync.RWMutex
	claudeHome string
	stats      data.StatsCache
	history    []data.HistorySession
	costs      []data.SessionCost
	sessions   []data.Session
	lastRefresh time.Time
}

func newDataCache(claudeHome string) *dataCache {
	return &dataCache{claudeHome: claudeHome}
}

func (c *dataCache) refresh() {
	stats, _ := data.LoadStatsCache(c.claudeHome)
	history, _ := data.LoadHistory(c.claudeHome)
	costs, _ := data.LoadAllSessionCosts(c.claudeHome)
	sessions, _ := data.FetchAgents(c.claudeHome)
	for i := range sessions {
		sessions[i].Tokens = data.LoadSessionTokens(c.claudeHome, sessions[i].SessionID)
	}

	c.mu.Lock()
	c.stats = stats
	c.history = history
	c.costs = costs
	c.sessions = sessions
	c.lastRefresh = time.Now()
	c.mu.Unlock()
}

func (c *dataCache) refreshLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		c.refresh()
	}
}

func Run(cfg config.Config, port int, noOpen bool) {
	cache := newDataCache(cfg.Claude.Home)
	cache.refresh()
	go cache.refreshLoop(10 * time.Second)

	mux := http.NewServeMux()
	registerHandlers(mux, cache)

	addr := fmt.Sprintf("127.0.0.1:%d", port)
	fmt.Printf("c9s web dashboard: http://%s\n", addr)

	if !noOpen {
		openBrowser(fmt.Sprintf("http://%s", addr))
	}

	srv := &http.Server{Addr: addr, Handler: mux}
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		fmt.Println("\nShutting down...")
		srv.Shutdown(context.Background())
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}

func openBrowser(url string) {
	var cmd string
	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
	case "linux":
		cmd = "xdg-open"
	case "windows":
		cmd = "rundll32"
		exec.Command(cmd, "url.dll,FileProtocolHandler", url).Start()
		return
	default:
		return
	}
	exec.Command(cmd, url).Start()
}
