FROM golang:1.25-alpine AS builder

WORKDIR /app/argus-agent

COPY proto ../proto
COPY argus-agent/go.mod argus-agent/go.sum ./
RUN go mod download

COPY argus-agent .

WORKDIR /app/argus-agent

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /app/argus-agent ./cmd/agent/

FROM alpine:3.23

WORKDIR /app

COPY --from=builder /app/argus-agent .

RUN adduser -D appuser
USER appuser

ENV ARGUS_SERVER_ADDR=argus-server:50051
ENV ARGUS_AGENT_ID=""

ENTRYPOINT ["./argus-agent"]