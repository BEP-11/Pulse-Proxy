package proxy
io

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/BEP-11/proxy-pulse/internal/security"
)

type Config struct {
	ListenAddr string
	RateLimit  float64
}

type Server struct {
	addr    string
	ln      net.Listener
	limiter *security.RateLimiter
	rateLimit float64
	wg      sync.WaitGroup
	httpSrv *http.Server
}

func NewServer(c *Config) *Server {
	return &Server{
		addr:      c.ListenAddr,
		limiter:   security.NewRateLimiter(c.RateLimit),
		rateLimit: c.RateLimit,
	}
}

func (s *Server) Addr() string { return s.addr }

func (s *Server) Start(ctx context.Context) error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	s.ln = ln

	s.httpSrv = &http.Server{Handler: http.HandlerFunc(s.handleProxy)}
	go func() {
		<-ctx.Done()
		s.httpSrv.Shutdown(context.Background())
	}()
	return s.httpSrv.Serve(ln)
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.ln != nil { s.ln.Close() }
	if s.httpSrv != nil { return s.httpSrv.Shutdown(ctx) }
	return nil
}

func (s *Server) handleProxy(w http.ResponseWriter, r *http.Request) {
	remote := r.RemoteAddr
	if !s.limiter.Allow(remote) {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	conn, _, err := hijacker.Hijack()
	if err != nil {
		log.Printf("❌ Hijack failed %s: %v", remote, err)
		return
	}
	s.wg.Add(1)
	go s.tunnel(r.Context(), conn, remote)
}

func (s *Server) tunnel(ctx context.Context, client net.Conn, remote string) {
	defer client.Close()
	defer s.wg.Done()

	// Пример маршрутизации на публичный DNS/TCP (замените на свой upstream)
	server, err := net.DialTimeout("tcp", "1.1.1.1:443", 5*time.Second)
	if err != nil {
		log.Printf("❌ Dial fail %s: %v", remote, err)
		return
	}
	defer server.Close()

	go func() { io.Copy(server, client); server.CloseWrite() }()
	io.Copy(client, server)
	client.CloseWrite()
}
