FROM golang:1.25-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/servermanager ./cmd/servermanager

FROM alpine:3.21 AS runner

# No USER drop: the manager needs /var/run/docker.sock, which is
# root-equivalent regardless — that's exactly why its API requires the
# bearer token and user containers never get the socket.
COPY --from=builder /app/bin/servermanager /usr/local/bin/servermanager

EXPOSE 8080
ENTRYPOINT ["servermanager"]
