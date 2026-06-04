package proxy

import (
	"context"
	"crypto/tls"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/user/proxy-pulse/internal/security"
)

type Config struct {
	ListenAddr string
	TLSConfig  *tls.Config
	RateLimit  float64
}

type Server struct {
	ln        net.Listener
	limiter   *security.RateLimiter
	rateLimit float64
	wg        sync.WaitGroup
}

func NewServer(c *Config) *Server {
	return &Server{
		ln:        nil, // set in Start()
		limiter:   security.NewRateLimiter(c.RateLimit),
		rateLimit: c.RateLimit,
	}
}

func (s *Server) Start(ctx context.Context) error {
	addr := s.ln.Addr().String()
	if s.rateLimit > 0 {
		go s.rateLimiterCleanup(30)
	}
	return http.ServeTLS(s.ln, s.handleRequest, s.TLSConfig)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.ln.Close()
}

func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	remote := r.RemoteAddr
	if !s.limiter.Allow(remote) {
		http.Error(w, "Too many requests", http.StatusTooManyRequests)
		return
	}

	conn, err := w.(http.Hijacker).Hijack()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	go s.tunnel(r.Context(), conn, remote)
}

func (s *Server) tunnel(ctx context.Context, client net.Conn, remote string) {
	defer client.Close()

	// Simple TCP tunnel to upstream (replace with smart router in prod)
	upstream := "1.1.1.1:443"
	server, err := net.DialTimeout("tcp", upstream, 5*time.Second)
	if err != nil {
		log.Printf("❌ Dial failed %s: %v", remote, err)
		return
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		io.Copy(server, client)
		server.CloseWrite()
	}()

	io.Copy(client, server)
	client.CloseWrite()
}
