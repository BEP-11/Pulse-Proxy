package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Security SecurityConfig `mapstructure:"security"`
	Metrics  MetricsConfig  `mapstructure:"metrics"`
}

type ServerConfig struct {
	Listen string `mapstructure:"listen"`
	Port   int    `mapstructure:"port"`
}

type SecurityConfig struct {
	RateLimitQPS float64 `mapstructure:"rate_limit_qps"`
}

type MetricsConfig struct {
	Addr string `mapstructure:"addr"`
}

func Load(path string) (*Config, error) {
	v := viper.New()
	if path != "" {
		v.SetConfigFile(path)
	} else {
		v.SetDefault("server.listen", "0.0.0.0")
		v.SetDefault("server.port", 8443)
		v.SetDefault("security.rate_limit_qps", 50)
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
