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

	"github.com/user/proxy-pulse/internal/config"
	"github.com/user/proxy-pulse/internal/metrics"
	"github.com/user/proxy-pulse/internal/proxy"
	"github.com/user/proxy-pulse/internal/tls"
)

func main() {
	cfg, err := config.Load("examples/config.yaml")
	if err != nil {
		log.Fatalf("❌ Config load failed: %v", err)
	}

	tlsCfg := tls.NewManager(cfg.Server.SNI)

	srvProxy := proxy.NewServer(&proxy.Config{
		ListenAddr: cfg.Server.Listen,
		TLSConfig:  tlsCfg,
		RateLimit:  cfg.Security.RateLimitQPS,
	})

	metricsSrv := metrics.NewServer(cfg.Metrics.Addr)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		fmt.Printf("🚀 Proxy starting on %s\n", cfg.Server.Listen)
		if err := srvProxy.Start(ctx); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Proxy failed: %v", err)
		}
	}()

	go func() {
		fmt.Printf("📊 Metrics starting on %s\n", cfg.Metrics.Addr)
		if err := metricsSrv.Start(); err != nil {
			log.Fatalf("❌ Metrics failed: %v", err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println("🛑 Graceful shutdown...")
	srvProxy.Shutdown(shutdownCtx)
	metricsSrv.Shutdown(shutdownCtx)
	log.Println("✅ Stopped.")
}
