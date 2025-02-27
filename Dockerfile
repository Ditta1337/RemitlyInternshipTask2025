FROM golang:1.24.0 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/api

FROM alpine:latest
WORKDIR /root/

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/server .
COPY --from=builder /app/cmd/migrations /root/cmd/migrations
COPY --from=builder /app/internal/db/seed /root/internal/db/seed
COPY --from=builder /app/docs /root/docs


RUN chmod +x /root/server

EXPOSE 8080

CMD ["/root/server"]
