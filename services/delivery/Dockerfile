FROM golang:1.26.3-alpine AS builder

WORKDIR /src/services/delivery
COPY services/delivery/go.mod ./
RUN go mod download

COPY services/delivery .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/server cmd/server/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /bin/server .
CMD ["./server"]