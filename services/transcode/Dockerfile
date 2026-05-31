FROM golang:1.26.3-alpine AS builder

COPY pkg/events /src/pkg/events

WORKDIR /src/services/transcode
COPY services/transcode/go.mod services/transcode/go.sum ./
RUN go mod download

COPY services/transcode .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/worker cmd/worker/main.go

FROM alpine:latest
RUN apk add --no-cache ffmpeg
WORKDIR /app
COPY --from=builder /bin/worker .
CMD ["./worker"]
