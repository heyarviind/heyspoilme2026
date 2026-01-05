#!/bin/bash

# HeySpoilMe - Development Start Script
# This script starts all services for local development

set -e

echo "🚀 Starting HeySpoilMe Development Environment"
echo "================================================"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if docker is running
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}Error: Docker is not running. Please start Docker first.${NC}"
    exit 1
fi

# Start PostgreSQL in Docker
echo -e "${YELLOW}Starting PostgreSQL...${NC}"
docker-compose up -d postgres

# Wait for PostgreSQL to be ready
echo -e "${YELLOW}Waiting for PostgreSQL to be ready...${NC}"
until docker-compose exec -T postgres pg_isready -U postgres > /dev/null 2>&1; do
    sleep 1
done
echo -e "${GREEN}PostgreSQL is ready!${NC}"

# Function to cleanup on exit
cleanup() {
    echo -e "\n${YELLOW}Shutting down services...${NC}"
    kill $BACKEND_PID 2>/dev/null || true
    kill $FRONTEND_PID 2>/dev/null || true
    docker-compose stop postgres
    echo -e "${GREEN}All services stopped.${NC}"
}

trap cleanup EXIT

# Start Backend
echo -e "${YELLOW}Starting Backend (Go)...${NC}"
cd backend

# Load .env file if it exists
if [ -f ".env" ]; then
    echo -e "${YELLOW}Loading environment from .env file...${NC}"
    set -a
    source .env
    set +a
fi

export DATABASE_URL="postgres://postgres:postgres@localhost:5433/heyspoilme?sslmode=disable"
export PORT=8080
export FRONTEND_URL="http://localhost:3003"
export JWT_SECRET="dev-secret-change-in-production"

# Cloudflare R2 / S3 Configuration (loaded from .env or environment)
export AWS_ACCESS_KEY_ID="${AWS_ACCESS_KEY_ID:-}"
export AWS_SECRET_ACCESS_KEY="${AWS_SECRET_ACCESS_KEY:-}"
export AWS_REGION="${AWS_REGION:-auto}"
export S3_BUCKET="${S3_BUCKET:-}"
export S3_ENDPOINT="${S3_ENDPOINT:-}"
export S3_BASE_URL="${S3_BASE_URL:-}"

# ZeptoMail Configuration (loaded from .env or environment)
export ZEPTOMAIL_API_KEY="${ZEPTOMAIL_API_KEY:-}"
export ZEPTOMAIL_FROM_EMAIL="${ZEPTOMAIL_FROM_EMAIL:-}"
export ZEPTOMAIL_FROM_NAME="${ZEPTOMAIL_FROM_NAME:-}"

# Google OAuth (loaded from .env or environment)
export GOOGLE_CLIENT_ID="${GOOGLE_CLIENT_ID:-}"
export GOOGLE_CLIENT_SECRET="${GOOGLE_CLIENT_SECRET:-}"

export ADMIN_CODE_1="${ADMIN_CODE_1:-}"
export ADMIN_CODE_2="${ADMIN_CODE_2:-}"

# Download dependencies if go.sum doesn't exist
if [ ! -f "go.sum" ]; then
    echo -e "${YELLOW}Downloading Go dependencies...${NC}"
    go mod tidy
fi

# Run migrations if migrate is installed
if command -v migrate &> /dev/null; then
    echo -e "${YELLOW}Running database migrations...${NC}"
    migrate -path ./migrations -database "$DATABASE_URL" up 2>/dev/null || true
fi

go run ./cmd/server &
BACKEND_PID=$!
cd ..

# Wait for backend to start
sleep 3
if kill -0 $BACKEND_PID 2>/dev/null; then
    echo -e "${GREEN}Backend running on http://localhost:8080${NC}"
else
    echo -e "${RED}Backend failed to start${NC}"
fi

# Start Frontend
echo -e "${YELLOW}Starting Frontend (SvelteKit)...${NC}"
cd frontend
pnpm install --silent
pnpm dev &
FRONTEND_PID=$!
cd ..

echo ""
echo -e "${GREEN}================================================${NC}"
echo -e "${GREEN}HeySpoilMe Development Environment Started!${NC}"
echo -e "${GREEN}================================================${NC}"
echo ""
echo -e "  Frontend:  ${GREEN}http://localhost:3003${NC}"
echo -e "  Backend:   ${GREEN}http://localhost:8080${NC}"
echo -e "  Database:  ${GREEN}localhost:5433${NC}"
echo ""
echo -e "${YELLOW}Note: Set GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET for OAuth${NC}"
echo -e "${YELLOW}Press Ctrl+C to stop all services${NC}"
echo ""

# Wait for any process to exit
wait

