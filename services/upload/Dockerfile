FROM golang:1.26.3-alpine AS builder

COPY pkg/events /src/pkg/events

WORKDIR /src/services/upload
COPY services/upload/go.mod services/upload/go.sum ./
RUN go mod download

COPY services/upload .
RUN go build -o /bin/server cmd/server/main.go

FROM alpine:latest
RUN apk add --no-cache ffmpeg
WORKDIR /app
COPY --from=builder /bin/server .
EXPOSE 50052
CMD ["./server"]
