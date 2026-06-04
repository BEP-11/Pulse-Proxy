
# 🔒 ProxyPulse
High-performance TLS proxy + load balancer + observability. Zero-config mode for dev, YAML for prod.

## ⚡ Quick Start
```bash
git clone https://github.com/BEP-11/proxy-pulse.git
cd proxy-pulse
make run   # starts on 0.0.0.0:443 + metrics :9090

 Как активировать конфиг?
В вашем cmd/server/main.go уже используется viper. Запустите с нужным файлом:

PROXY_CONFIG=config.prod.yaml go run ./cmd/server/
Или через env-переменные (Viper автоматически мапит):

export SERVER_PORT=443
export SECURITY_RATE_LIMIT_QPS=50
export TLS_ACME=true
go run ./cmd/server/ --config config.dev.yam
