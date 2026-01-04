.PHONY: build run dev migrate-up migrate-down migrate-create test clean

# Database URL
DATABASE_URL ?= postgres://postgres:postgres@localhost:5433/heyspoilme?sslmode=disable

# Build the application
build:
	go build -o bin/server ./cmd/server

# Run the application
run: build
	./bin/server

# Run in development mode with hot reload
dev:
	go run ./cmd/server

# Run database migrations up
migrate-up:
	migrate -path ./migrations -database "$(DATABASE_URL)" up

# Run database migrations down
migrate-down:
	migrate -path ./migrations -database "$(DATABASE_URL)" down

# Run specific number of migrations down
migrate-down-one:
	migrate -path ./migrations -database "$(DATABASE_URL)" down 1

# Check migration status
migrate-status:
	migrate -path ./migrations -database "$(DATABASE_URL)" version

# Create a new migration
migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir ./migrations -seq $$name

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Install dependencies
deps:
	go mod download
	go mod tidy

# Install migrate CLI tool
install-migrate:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

