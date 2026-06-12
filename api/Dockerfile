FROM golang:1.26-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/api ./cmd/api

FROM alpine:3.21 AS runner

RUN addgroup -S app && adduser -S app -G app
COPY --from=builder /app/bin/api /usr/local/bin/api

USER app
EXPOSE 8080
ENTRYPOINT ["api"]
