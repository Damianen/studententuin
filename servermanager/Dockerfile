FROM golang:1.25-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/servermanager ./cmd/servermanager

# The nixpacks release is a static musl build, so it runs on alpine. Keep the
# version in lockstep with the CI install step (servermanager-tests.yml).
# x86_64-only: matches the hosting server and CI runners.
FROM alpine:3.21 AS nixpacks
ARG NIXPACKS_VERSION=1.41.0
RUN apk add --no-cache ca-certificates curl \
 && curl -fsSL -o /tmp/nixpacks.tar.gz "https://github.com/railwayapp/nixpacks/releases/download/v${NIXPACKS_VERSION}/nixpacks-v${NIXPACKS_VERSION}-x86_64-unknown-linux-musl.tar.gz" \
 && tar -xzf /tmp/nixpacks.tar.gz -C /usr/local/bin nixpacks

FROM alpine:3.21 AS runner

# ca-certificates: go-git clones public repos over https.
RUN apk add --no-cache ca-certificates

# No USER drop: the manager needs /var/run/docker.sock, which is
# root-equivalent regardless — that's exactly why its API requires the
# bearer token and user containers never get the socket.
COPY --from=nixpacks /usr/local/bin/nixpacks /usr/local/bin/nixpacks
COPY --from=builder /app/bin/servermanager /usr/local/bin/servermanager

EXPOSE 8080
ENTRYPOINT ["servermanager"]
