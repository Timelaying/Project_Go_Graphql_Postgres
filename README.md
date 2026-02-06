# Job Tracker (Go + GraphQL + Postgres)

A lightweight job-application tracker built with Go, gqlgen (GraphQL), and PostgreSQL. The repository contains the GraphQL schema, database models, and helper utilities for a jobs table, plus a simple HTTP server that hosts the GraphQL playground and API endpoint.

## Features

- GraphQL API scaffolded with gqlgen.
- PostgreSQL-backed `jobs` table with status and timestamps.
- Helper package for database connections and SQL migrations.
- Docker Compose setup for local Postgres + Adminer.

> **Project status**: the GraphQL server is fully wired to Postgres with validation, demo seed data, and consistent error codes.

## Tech Stack

- **Go** (module: `job-tracker`)
- **GraphQL** via [gqlgen](https://github.com/99designs/gqlgen)
- **PostgreSQL**
- **Docker Compose** (for local Postgres and Adminer)

## Project Structure

```
.
├── server.go                # HTTP server + GraphQL playground
├── internal/
│   ├── db/                  # db pool + migration helpers
│   ├── graph/               # GraphQL schema + generated code
│   └── jobs/                # Job domain model + repo
├── migrations/              # SQL migrations (schema)
└── docker-compose.yml       # Local Postgres + Adminer
```

## Getting Started

### Prerequisites

- Go installed (see `go.mod` for the version reference).
- Docker (optional, for local Postgres).

### 1) Start Postgres (Docker)

```bash
make up
```

Adminer is available at `http://localhost:8080` (default). Note that the Go server also defaults to port 8080, so you may want to set `PORT=8081` when running the API locally.

### 2) Configure Environment

```bash
cp .env.example .env
```

### 3) Run the API Server

```bash
make run
```

- GraphQL Playground: `http://localhost:4000/`
- GraphQL endpoint: `http://localhost:4000/query`

## Environment Variables

- `PORT`: HTTP server port (default: `8080`).
- `DATABASE_URL`: Postgres connection string.

For database connectivity, the helpers in `internal/db` expect a Postgres URL. A common local URL is:

```
postgres://postgres:postgres@localhost:5432/jobtracker?sslmode=disable
```

## Database Schema

The `jobs` table is created by `migrations/001_init.sql` with these fields:

- `id` (UUID, primary key)
- `company` (text)
- `role` (text)
- `link` (text, optional)
- `status` (enum-like text: `APPLIED`, `INTERVIEW`, `OFFER`, `REJECTED`)
- `created_at` (timestamp)

Stats query example:

```sql
SELECT status, COUNT(*) FROM jobs GROUP BY status ORDER BY status;
```

## GraphQL Schema

The GraphQL schema is defined in `internal/graph/schema.graphqls` and includes:

- **Queries**: `jobs`, `job`, `statsByStatus`
- **Mutations**: `createJob`, `updateJobStatus`, `updateJobLink`, `deleteJob`, `seedDemoJobs`
- **Types**: `Job`, `Status`, `StatusCount`

You can use the playground to execute queries like:

```graphql
query GetJobs {
  jobs {
    id
    company
    role
    status
    createdAt
  }
}
```

Sample GraphQL operations:

```graphql
mutation {
  createJob(company: "Monzo", role: "Backend Engineer", link: "https://example.com") {
    id
    company
    role
    status
    createdAt
  }
}
```

```graphql
query {
  jobs(status: APPLIED, q: "Engineer") {
    id
    company
    role
    link
    status
    createdAt
  }
}
```

```graphql
mutation {
  updateJobStatus(id: "<uuid>", status: INTERVIEW) {
    id
    status
  }
}
```

```graphql
mutation {
  updateJobLink(id: "<uuid>", link: "https://example.com/job") {
    id
    link
  }
}
```

```graphql
mutation {
  deleteJob(id: "<uuid>")
}
```

```graphql
mutation {
  seedDemoJobs
}
```

```graphql
query {
  statsByStatus {
    status
    count
  }
}
```

## Development Notes

- gqlgen configuration lives in `gqlgen.yml`.
- Generated code is under `internal/graph/`.
- The job repository (`internal/jobs/repository.go`) provides CRUD operations against Postgres.

## Definition of Done

- ✅ docker compose up spins Postgres.
- ✅ App auto-applies `migrations/001_init.sql` on startup.
- ✅ Playground opens and all operations work (`jobs`, `job`, `createJob`, `updateJobStatus`, `updateJobLink`, `deleteJob`, `seedDemoJobs`, `statsByStatus`).
- ✅ Filters + search work (`status`, `q`).
- ✅ Consistent GraphQL error codes for bad input + not found.
- ✅ README has setup + sample operations.

## License

This project is currently unlicensed. Add a LICENSE file if you plan to distribute it.
