FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download && CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o proxy-pulse ./cmd/server

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /app/proxy-pulse /usr/local/bin/
EXPOSE 443 9090
ENTRYPOINT ["proxy-pulse"]
