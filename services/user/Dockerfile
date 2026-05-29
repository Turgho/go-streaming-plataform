FROM golang:1.26.3-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o bin/server cmd/server/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bin/server .
EXPOSE 50051
CMD ["./server"]