package tls

import "crypto/tls"

// В production замените на certmagic или загрузку PEM-файлов
func NewManager() *tls.Config {
	return &tls.Config{MinVersion: tls.VersionTLS12}
}
