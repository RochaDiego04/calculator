# ---- build ----
FROM golang:1.27.1 AS build
WORKDIR /src

# Manifests first, so the module download layer only busts when deps change.
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# CGO_ENABLED=0 produces a genuinely static binary. Without it the binary can
# link against the build image's libc and die on distroless with a loader error.
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app ./

# ---- run ----
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /app /app
EXPOSE 8080
USER nonroot:nonroot

# Distroless has no shell and no curl, so `CMD curl ...` would fail forever
# while the service works fine. The binary probes itself instead.
HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=3 \
  CMD ["/app", "-health"]

# Exec form is required. The shell form makes /bin/sh PID 1, SIGTERM never
# reaches the process, and graceful shutdown silently does nothing.
ENTRYPOINT ["/app"]
