# Kenjix Persist

Small Go service for Kenjix persistence layer.

## Requirements

- Go 1.24+
- Docker (optional)

## Build & Run (local)

Run with `go run` from repository root:

```bash
go run ./cmd/api
```

Build binary:

```bash
go build -o bin/kenjix ./cmd/api
./bin/kenjix    # or .\bin\kenjix.exe on Windows
```

The server listens on `PORT` (default `8080`). Health endpoint: `/health`.

## Docker

Build image (run from repo root where `go.mod` lives):

```bash
# run from repository root
docker build -t kenjix:latest -f docker/Dockerfile .
docker run --rm -p 8080:8080 kenjix:latest
```

If you run `docker build` from inside the `docker/` folder, pass `..` as context:

```bash
cd docker
docker build -t kenjix:latest -f Dockerfile ..
```

## Docker Compose

Start (from repo root):

```bash
docker compose -f docker/docker-compose.yml up -d --build
```

Stop and remove:

```bash
docker compose -f docker/docker-compose.yml down -v
```

## VS Code Launch

The project includes a `.vscode/launch.json` configured to launch `cmd/api`.

## Generating `go.sum`

If `go.sum` is missing, run from repo root:

```bash
go mod tidy
```

## Notes

- Ensure you run Docker builds with the repository root as the build context so `go.mod` is included.
- Avoid `container_name` in `docker-compose.yml` if you plan to scale services.