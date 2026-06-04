package metrics

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	ActiveConns = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "proxy_pulse_active_connections_total", Help: "Current active connections"})
	RateLimited = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "proxy_pulse_rate_limited_total", Help: "Connections rejected by rate limiter"})
)

func init() {
	prometheus.MustRegister(ActiveConns, RateLimited)
}

type Server struct {
	addr string
	srv  *http.Server
}

func NewServer(addr string) *Server { return &Server{addr: addr} }

func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ProxyPulse Metrics"))
	})
	s.srv = &http.Server{Addr: s.addr, Handler: mux}
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error { return s.srv.Shutdown(ctx) }
