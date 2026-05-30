FROM golang:1.26.3-alpine AS builder

WORKDIR /src/services/user
COPY services/user/go.mod services/user/go.sum ./
RUN go mod download

COPY services/user .
RUN go build -o /bin/server cmd/server/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /bin/server .
EXPOSE 50051
CMD ["./server"]
