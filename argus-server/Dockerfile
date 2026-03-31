FROM golang:1.25-alpine AS builder
LABEL authors="noboaki"

WORKDIR /app

COPY go.work go.work.sum ./
COPY /proto ./proto
COPY /argus-server ./argus-server

RUN go work sync
RUN go build -o argus-server ./argus-server/cmd/server

FROM alpine:3.23

WORKDIR /app

COPY --from=builder /app/argus-server .

EXPOSE 50051

ENTRYPOINT ["./argus-server"]
