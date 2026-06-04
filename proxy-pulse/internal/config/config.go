package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Security SecurityConfig
	Metrics  MetricsConfig
}

type ServerConfig struct {
	Listen string
	Port   int
	SNI    string
}

type SecurityConfig struct {
	RateLimitQPS float64
	BlockSubnets []string
}

type MetricsConfig struct {
	Addr string
}

func Load(path string) (*Config, error) {
	v := viper.New()
	if path != "" {
		v.SetConfigFile(path)
	} else {
		v.SetConfigType("yaml")
		v.SetDefault("server.listen", "0.0.0.0")
		v.SetDefault("server.port", 443)
		v.SetDefault("server.sni", "")
		v.SetDefault("security.rate_limit_qps", 100)
		v.SetDefault("metrics.addr", ":9090")
	}

	if err := v.ReadInConfig(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
