FROM golang:1.25-alpine AS builder
LABEL authors="noboaki"

WORKDIR /app

COPY go.work go.work.sum ./
COPY /proto ./proto
COPY /argus-agent ./argus-agent

RUN go work sync
RUN go build -o argus-agent ./argus-agent/cmd/agent/

FROM alpine:3.23

WORKDIR /app

COPY --from=builder /app/argus-agent .

ENV ARGUS_SERVER_ADDR=argus-server:50051
ENV ARGUS_AGENT_ID=""

ENTRYPOINT ["./argus-agent"]
