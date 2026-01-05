# HeySpoilMe

A premium exclusive dating platform built on trust and discretion.

## Tech Stack

- **Frontend**: SvelteKit 5 with TypeScript
- **Backend**: Go with Gin framework
- **Database**: PostgreSQL

## Quick Start with Docker

```bash
# Build and start all services
docker-compose up --build

# Or run in detached mode
docker-compose up -d --build
```

The services will be available at:
- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **Database**: localhost:5432

## Development

### Frontend

```bash
cd frontend
pnpm install
pnpm dev
```

### Backend

```bash
cd backend
go mod download
go run main.go
```

Make sure PostgreSQL is running and accessible at the DATABASE_URL specified in your environment.

## API Endpoints

### Health Check
```
GET /health
```

### Subscribe (Email Collection)
```
POST /api/subscribe
Content-Type: application/json

{
  "email": "user@example.com"
}
```

## Environment Variables

See `.env.example` for all available configuration options.

## License

Private - All rights reserved

