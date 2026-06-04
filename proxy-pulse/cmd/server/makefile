.PHONY: run build test deploy bench lint

run:
	go run ./cmd/server/

build:
	mkdir -p bin && CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/proxy-pulse ./cmd/server/

test:
	go test ./... -coverprofile=coverage.out

deploy: build
	docker compose -f deploy/docker-compose.yml up -d --build

bench: build
	wrk -t12 -c400 -d30s http://localhost:9090/metrics

lint:
	golangci-lint run ./...

clean:
	rm -rf bin/ cache/ coverage.out
