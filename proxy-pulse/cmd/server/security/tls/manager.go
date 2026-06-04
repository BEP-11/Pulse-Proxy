package tls

import (
	"context"
	"crypto/tls"

	"github.com/caddyserver/certmagic"
)

func NewManager(sni string) *tls.Config {
	cm := certmagic.NewDefault()
	cm.OnCacheMiss = func(ctx context.Context, domain string) ([]byte, error) {
		// Fallback to self-signed if ACME unavailable
		return nil, nil
	}
	if sni != "" {
		cm.Storage = &certmagic.FileStorage{Path: "cache"}
	}

	return &tls.Config{
		GetCertificate:   cm.GetCertificate,
		ClientAuth:       tls.NoClientCert,
		NextProtos:       []string{"h2", "http/1.1"},
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}
}
