FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY main.go ./
COPY banner.txt ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/sgen ./main.go

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app .

ENTRYPOINT ["./sgen"]

LABEL org.opencontainers.image.source=https://github.com/hnnsly/tg-sgen
LABEL org.opencontainers.image.description="Semi-automated Telegram session creation tool built with Gogram."
