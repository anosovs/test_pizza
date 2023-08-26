FROM golang:1.20 as builder
WORKDIR /app
COPY . /app
RUN GO111MODULE=auto CGO_ENABLED=0 GOOS=linux GOPROXY=https://proxy.golang.org go build -o app cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app .
COPY --from=builder /app/config/.env ./config/.env

ENTRYPOINT ["./app"]
