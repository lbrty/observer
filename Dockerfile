# Stage 1: build frontend
FROM oven/bun:1 AS frontend
WORKDIR /app
COPY package.json bun.lock ./
COPY packages/observer-web/package.json packages/observer-web/
RUN bun install --frozen-lockfile
COPY packages/observer-web/ packages/observer-web/
RUN cd packages/observer-web && bun run build

# Stage 2: build backend
FROM golang:1.25 AS backend
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /observer ./cmd/observer

# Stage 3: final image
FROM scratch
COPY --from=backend /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=backend /observer /observer
COPY --from=frontend /app/packages/observer-web/dist /spa
COPY migrations /migrations
ENTRYPOINT ["/observer"]
CMD ["serve", "--host", "0.0.0.0"]
