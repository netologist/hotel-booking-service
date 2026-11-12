# Hotel Booking

A Go service for registering and finding hotels, searching room availability, and booking rooms. It uses a spec-first REST API, PostgreSQL, and hexagonal architecture.

## Prerequisites

- [Mise](https://mise.jdx.dev/) for the pinned Go toolchain and development tasks
- [Docker and Docker Compose](https://docs.docker.com/get-docker/) for PostgreSQL

## Quickstart

```bash
git clone https://github.com/netologist/hotel-booking-service.git
cd hotel-booking-service
mise install
cp .env.example .env
mise run docker:compose-up
mise run migrate:up
mise exec -- go run ./cmd/server
```

The server listens on `http://localhost:8080` by default. Check its health with:

```bash
curl http://localhost:8080/healthz
```

Configuration is loaded from environment variables; `.env.example` lists the local defaults. Mise loads `.env` when running tasks and `mise exec` commands.

## API

| Method | Endpoint | Purpose |
| --- | --- | --- |
| GET | `/healthz` | Health check |
| GET | `/api/v1/hotels?name=grand` | Find hotels by name |
| POST | `/api/v1/hotels` | Register a hotel |
| GET | `/api/v1/rooms/available` | Find available rooms by dates and guest count |
| POST | `/api/v1/bookings` | Book a room |
| GET | `/api/v1/bookings/{bookingReference}` | Retrieve a booking |
| POST | `/api/v1/seed` | Seed test data |
| POST | `/api/v1/reset` | Remove test data |

The [OpenAPI specification](api/openapi.yaml) defines request parameters and response schemas. When API documentation is enabled, Swagger UI is available at `http://localhost:8080/docs/` and the specification at `http://localhost:8080/openapi.json`. Both are disabled in production.

## Development

```bash
mise run dev             # Start PostgreSQL, apply migrations, and run with live reload
mise run generate        # Regenerate API, client, SQL, converter, and mock code
mise run test            # Unit, client, and architecture tests with race detection
mise run test:arch       # Architecture boundary tests
mise run test:e2e        # PostgreSQL-backed end-to-end tests
mise run test:coverage   # Unit test coverage
mise run lint            # Lint Go packages
mise run vuln            # Check dependencies for known vulnerabilities
mise run docs:serve      # Serve the MkDocs documentation locally
mise run docs:build      # Build the documentation site
mise run docker:build    # Build the production container
```

## Architecture and documentation

The domain and use cases are separated from HTTP and PostgreSQL adapters. OpenAPI drives the HTTP contract, `sqlc` generates database access code, and Goose manages migrations.

- [Project structure and architecture](docs/folder-structure.md)
- [Technology stack](docs/tech-stack.md)
- [CI/CD and security](docs/pipeline.md)
- [Development guide](docs/user-guide.md)
- [Architecture decisions](docs/decisions/0001-development-environment.md)
- [API reference](docs/api.md)
