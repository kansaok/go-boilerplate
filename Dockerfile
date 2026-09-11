# Build stage
FROM golang:1.22.5-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git make

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Production stage
FROM alpine:3.20

RUN apk --no-cache add ca-certificates

WORKDIR /app

RUN mkdir -p /app/storage/sqlite /app/uploads /app/log

COPY --from=builder /app/main .

RUN adduser -D -u 1000 appuser && \
    chown -R appuser:appuser /app

USER appuser

EXPOSE 8080

CMD ["./main"]
