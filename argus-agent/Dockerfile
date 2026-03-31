FROM golang:1.25-alpine AS builder

WORKDIR /app/argus-agent

COPY proto ../proto
COPY argus-agent/go.mod argus-agent/go.sum ./
RUN go mod download

COPY argus-agent .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /out/argus-agent ./cmd/agent

FROM alpine:3.23

WORKDIR /app

COPY --from=builder /out/argus-agent /app/argus-agent

RUN adduser -D appuser
USER appuser

ENV ARGUS_SERVER_ADDR=argus-server:50051
ENV ARGUS_AGENT_ID=""

ENTRYPOINT ["/app/argus-agent"]