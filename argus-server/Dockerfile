FROM golang:1.25-alpine AS builder

WORKDIR /app/argus-server

COPY proto ../proto
COPY argus-server/go.mod argus-server/go.sum ./
RUN go mod download

COPY argus-server .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /app/argus-server ./cmd/server/


FROM alpine:3.23

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/argus-server .

RUN adduser -D appuser
USER appuser

EXPOSE 50051

ENTRYPOINT ["./argus-server"]