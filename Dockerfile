# Stage 1: build frontend
FROM oven/bun:1.3 AS frontend

WORKDIR /app
COPY package.json bun.lock ./
COPY packages/observer-web/package.json packages/observer-web/
COPY packages/observer-web/bunfig.toml packages/observer-web/
RUN bun install

COPY packages/observer-web/ packages/observer-web/
RUN cd packages/observer-web && bun run build

# Stage 2: build backend with embedded frontend
FROM golang:1.25 AS backend

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=frontend /app/packages/observer-web/dist internal/spa/dist
RUN CGO_ENABLED=0 GOOS=linux go build -tags production -ldflags="-s -w" -o /observer ./cmd/observer

# Stage 3: final image
FROM scratch

COPY --from=backend /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=backend /observer /observer

EXPOSE 9000

ENTRYPOINT ["/observer"]
CMD ["serve", "--host", "0.0.0.0"]
