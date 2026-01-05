#!/bin/bash

# HeySpoilMe - Production Deployment Script
# This script handles the complete deployment lifecycle

set -e

echo ""
echo "╔═══════════════════════════════════════════╗"
echo "║     HeySpoilMe Deployment Script          ║"
echo "╚═══════════════════════════════════════════╝"
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default values
ACTION=${1:-"deploy"}
MODE=${2:-"prod"}

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# ===========================================
# Helper Functions
# ===========================================

# Add Go bin paths to PATH
export PATH="$PATH:$HOME/go/bin:/root/go/bin:/usr/local/go/bin"

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_dependencies() {
    log_info "Checking dependencies..."
    
    if ! docker info > /dev/null 2>&1; then
        log_error "Docker is not running. Please start Docker first."
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null; then
        log_error "docker-compose is not installed."
        exit 1
    fi
    
    log_success "All dependencies found"
}

check_env_file() {
    # .env file should be at the ROOT level (same directory as docker-compose files)
    if [ ! -f ".env" ]; then
        log_error ".env file not found!"
        echo ""
        echo -e "${YELLOW}The .env file should be in the project root directory:${NC}"
        echo -e "  ${SCRIPT_DIR}/.env"
        echo ""
        echo -e "Create it from the template:"
        echo -e "  cp .env.example .env"
        echo -e "  nano .env  # Edit with your actual values"
        exit 1
    fi
    
    log_success ".env file found at project root"
}

validate_env_vars() {
    log_info "Validating environment variables..."
    
    source .env
    
    MISSING_VARS=""
    
    [ -z "$POSTGRES_PASSWORD" ] && MISSING_VARS="$MISSING_VARS POSTGRES_PASSWORD"
    [ -z "$GOOGLE_CLIENT_ID" ] && MISSING_VARS="$MISSING_VARS GOOGLE_CLIENT_ID"
    [ -z "$GOOGLE_CLIENT_SECRET" ] && MISSING_VARS="$MISSING_VARS GOOGLE_CLIENT_SECRET"
    [ -z "$JWT_SECRET" ] && MISSING_VARS="$MISSING_VARS JWT_SECRET"
    [ -z "$AWS_ACCESS_KEY_ID" ] && MISSING_VARS="$MISSING_VARS AWS_ACCESS_KEY_ID"
    [ -z "$AWS_SECRET_ACCESS_KEY" ] && MISSING_VARS="$MISSING_VARS AWS_SECRET_ACCESS_KEY"
    [ -z "$S3_BASE_URL" ] && MISSING_VARS="$MISSING_VARS S3_BASE_URL"
    [ -z "$S3_ENDPOINT" ] && MISSING_VARS="$MISSING_VARS S3_ENDPOINT"
    [ -z "$ZEPTOMAIL_API_KEY" ] && MISSING_VARS="$MISSING_VARS ZEPTOMAIL_API_KEY"
    [ -z "$ADMIN_CODE_1" ] && MISSING_VARS="$MISSING_VARS ADMIN_CODE_1"
    [ -z "$ADMIN_CODE_2" ] && MISSING_VARS="$MISSING_VARS ADMIN_CODE_2"
    
    if [ -n "$MISSING_VARS" ]; then
        log_error "Missing required environment variables:"
        echo -e "${YELLOW}$MISSING_VARS${NC}"
        echo ""
        echo "Please update your .env file with the missing values."
        exit 1
    fi
    
    log_success "All required environment variables found"
}

# ===========================================
# Deployment Steps
# ===========================================

pull_latest() {
    log_info "Pulling latest changes from git..."
    
    if [ -d ".git" ]; then
        git fetch origin
        git pull origin main --ff-only || git pull origin master --ff-only || {
            log_warn "Could not fast-forward. Manual merge may be required."
        }
        log_success "Git pull complete"
    else
        log_warn "Not a git repository, skipping git pull"
    fi
}

stop_conflicting_services() {
    log_info "Checking for port conflicts..."
    
    source .env
    FRONTEND_PORT="${FRONTEND_PORT:-3003}"
    BACKEND_PORT="8081"
    
    # Check if ports are in use by non-heyspoilme services
    for port in $FRONTEND_PORT $BACKEND_PORT; do
        if lsof -i :$port > /dev/null 2>&1; then
            # Check if it's our container
            if docker ps --format '{{.Names}}' | grep -q "heyspoilme"; then
                log_info "Port $port in use by heyspoilme container (will be replaced)"
            else
                log_warn "Port $port is in use by another process:"
                lsof -i :$port | head -3 || true
                log_error "Please stop the process using port $port first"
                exit 1
            fi
        fi
    done
    
    log_success "Port checks complete"
}

run_migrations() {
    log_info "Running database migrations..."
    
    source .env
    
    # Wait for database to be ready
    log_info "Waiting for database to be ready..."
    
    # Construct the database URL for external access
    DB_HOST="${DB_HOST:-localhost}"
    DB_PORT="${DB_PORT:-5434}"
    
    # For docker, we connect to the postgres container
    MIGRATION_DB_URL="postgres://${POSTGRES_USER:-postgres}:${POSTGRES_PASSWORD}@${DB_HOST}:${DB_PORT}/${POSTGRES_DB:-heyspoilme}?sslmode=disable"
    
    # Check if migrate is installed
    if ! command -v migrate &> /dev/null; then
        log_warn "migrate CLI not installed. Installing..."
        go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest || {
            log_warn "Could not install migrate CLI. Skipping migrations."
            log_warn "Migrations will need to be run manually."
            return 0
        }
    fi
    
    # Run migrations
    cd backend
    migrate -path ./migrations -database "$MIGRATION_DB_URL" up || {
        log_warn "Migration failed or no new migrations to run"
    }
    cd ..
    
    log_success "Migrations complete"
}

build_images() {
    log_info "Building Docker images..."
    
    # Build images without starting containers (for zero-downtime)
    docker-compose -f docker-compose.prod.yml build --parallel
    
    log_success "Docker images built"
}

deploy_services() {
    log_info "Deploying services with minimal downtime..."
    
    source .env
    
    # Start/update database first (it's stateful, needs to be up)
    log_info "Starting database..."
    docker-compose -f docker-compose.prod.yml up -d postgres
    
    # Wait for database to be healthy
    log_info "Waiting for database to be healthy..."
    for i in {1..30}; do
        if docker exec heyspoilme-db pg_isready -U ${POSTGRES_USER:-postgres} > /dev/null 2>&1; then
            log_success "Database is ready"
            break
        fi
        sleep 1
    done
    
    # Run migrations
    log_info "Running migrations..."
    if command -v migrate &> /dev/null; then
        DOCKER_DB_URL="postgres://${POSTGRES_USER:-postgres}:${POSTGRES_PASSWORD}@localhost:5434/${POSTGRES_DB:-heyspoilme}?sslmode=disable"
        cd backend
        migrate -path ./migrations -database "$DOCKER_DB_URL" up 2>/dev/null || log_warn "No new migrations or migration skipped"
        cd ..
    else
        log_warn "migrate CLI not found. Skipping migrations."
    fi
    
    # ============================================
    # ZERO-DOWNTIME DEPLOYMENT STRATEGY
    # ============================================
    # 1. Start new containers alongside old ones
    # 2. Wait for new containers to be healthy
    # 3. Update Caddy to route to new containers
    # 4. Remove old containers
    
    # Get current container IDs (if any)
    OLD_BACKEND=$(docker ps -q -f name=heyspoilme-backend 2>/dev/null || true)
    OLD_FRONTEND=$(docker ps -q -f name=heyspoilme-frontend 2>/dev/null || true)
    
    if [ -n "$OLD_BACKEND" ] || [ -n "$OLD_FRONTEND" ]; then
        log_info "Existing containers found - performing rolling update..."
        
        # Rename old containers temporarily
        if [ -n "$OLD_BACKEND" ]; then
            docker rename heyspoilme-backend heyspoilme-backend-old 2>/dev/null || true
        fi
        if [ -n "$OLD_FRONTEND" ]; then
            docker rename heyspoilme-frontend heyspoilme-frontend-old 2>/dev/null || true
        fi
        
        # Start new backend
        log_info "Starting new backend..."
        docker-compose -f docker-compose.prod.yml up -d --no-deps --force-recreate backend
        
        # Wait for new backend to be ready
        log_info "Waiting for backend to be ready..."
        for i in {1..30}; do
            if docker exec heyspoilme-backend wget -q --spider http://localhost:8081/health 2>/dev/null || \
               curl -sf http://localhost:8081/health > /dev/null 2>&1; then
                log_success "Backend is healthy"
                break
            fi
            # Also check if container is running
            if docker ps -q -f name=heyspoilme-backend -f status=running | grep -q .; then
                sleep 1
            else
                log_warn "Backend container not running yet..."
                sleep 2
            fi
        done
        sleep 2
        
        # Start new frontend
        log_info "Starting new frontend..."
        docker-compose -f docker-compose.prod.yml up -d --no-deps --force-recreate frontend
        
        # Wait for frontend to be ready
        log_info "Waiting for frontend to be ready..."
        for i in {1..30}; do
            if docker ps -q -f name=heyspoilme-frontend -f status=running | grep -q .; then
                log_success "Frontend container is running"
                break
            fi
            sleep 1
        done
        sleep 2
        
        # Remove old containers
        log_info "Removing old containers..."
        docker stop heyspoilme-backend-old 2>/dev/null || true
        docker rm heyspoilme-backend-old 2>/dev/null || true
        docker stop heyspoilme-frontend-old 2>/dev/null || true
        docker rm heyspoilme-frontend-old 2>/dev/null || true
        
        log_success "Rolling update complete - zero downtime achieved!"
    else
        log_info "No existing containers - performing fresh deployment..."
        
        # Deploy backend
        log_info "Deploying backend..."
        docker-compose -f docker-compose.prod.yml up -d --no-deps backend
        sleep 3
        
        # Deploy frontend
        log_info "Deploying frontend..."
        docker-compose -f docker-compose.prod.yml up -d --no-deps frontend
        
        log_success "Fresh deployment complete"
        log_info "Make sure your global Caddy is configured to proxy to localhost:3001 and localhost:8081"
    fi
}

show_status() {
    source .env 2>/dev/null || true
    FRONTEND_PORT="${FRONTEND_PORT:-3003}"
    
    echo ""
    echo -e "${GREEN}╔═══════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║         Deployment Complete!              ║${NC}"
    echo -e "${GREEN}╚═══════════════════════════════════════════╝${NC}"
    echo ""
    
    docker-compose -f docker-compose.prod.yml ps
    
    echo ""
    echo -e "  Frontend:  ${GREEN}http://localhost:${FRONTEND_PORT}${NC} → ${GREEN}https://heyspoil.me${NC}"
    echo -e "  Backend:   ${GREEN}http://localhost:8081${NC} → ${GREEN}https://api.heyspoil.me${NC}"
    echo ""
    echo -e "${YELLOW}Global Caddy should proxy:${NC}"
    echo -e "  heyspoil.me       → localhost:${FRONTEND_PORT}"
    echo -e "  api.heyspoil.me   → localhost:8081"
    echo ""
    echo -e "${YELLOW}View logs:${NC} docker-compose -f docker-compose.prod.yml logs -f"
    echo -e "${YELLOW}Stop all:${NC}  ./deploy.sh down"
    echo ""
}

cleanup_old_images() {
    log_info "Cleaning up old Docker images..."
    docker image prune -f
    log_success "Cleanup complete"
}

# ===========================================
# Main Commands
# ===========================================

case $ACTION in
    "deploy"|"up")
        echo -e "${BLUE}Starting full deployment...${NC}"
        echo ""
        
        check_dependencies
        check_env_file
        validate_env_vars
        pull_latest
        stop_conflicting_services
        build_images
        deploy_services
        cleanup_old_images
        show_status
        ;;
    
    "quick")
        # Quick deploy without git pull
        echo -e "${BLUE}Quick deployment (no git pull)...${NC}"
        echo ""
        
        check_dependencies
        check_env_file
        validate_env_vars
        stop_conflicting_services
        build_images
        deploy_services
        show_status
        ;;
    
    "down"|"stop")
        log_info "Stopping all services..."
        docker-compose -f docker-compose.prod.yml down
        log_success "All services stopped"
        ;;
    
    "restart")
        log_info "Restarting all services..."
        docker-compose -f docker-compose.prod.yml restart
        log_success "All services restarted"
        ;;
    
    "logs")
        docker-compose -f docker-compose.prod.yml logs -f ${3:-}
        ;;
    
    "status"|"ps")
        docker-compose -f docker-compose.prod.yml ps
        ;;
    
    "migrate")
        check_env_file
        source .env
        
        if ! command -v migrate &> /dev/null; then
            log_error "migrate CLI not installed. Run: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
            exit 1
        fi
        
        DOCKER_DB_URL="postgres://${POSTGRES_USER:-postgres}:${POSTGRES_PASSWORD}@localhost:5434/${POSTGRES_DB:-heyspoilme}?sslmode=disable"
        cd backend
        migrate -path ./migrations -database "$DOCKER_DB_URL" up
        cd ..
        log_success "Migrations complete"
        ;;
    
    "rebuild")
        log_info "Full rebuild (no cache)..."
        check_dependencies
        check_env_file
        validate_env_vars
        
        docker-compose -f docker-compose.prod.yml down
        docker-compose -f docker-compose.prod.yml build --no-cache
        docker-compose -f docker-compose.prod.yml up -d
        
        show_status
        ;;
    
    "clean")
        log_warn "This will remove all containers, volumes, and images!"
        read -p "Are you sure? (y/N): " confirm
        if [ "$confirm" = "y" ] || [ "$confirm" = "Y" ]; then
            docker-compose -f docker-compose.prod.yml down -v --rmi all
            docker system prune -f
            log_success "Full cleanup complete"
        else
            log_info "Cancelled"
        fi
        ;;
    
    "local")
        # Local development deployment
        echo -e "${BLUE}Local development deployment...${NC}"
        
        check_dependencies
        
        if [ ! -f ".env" ]; then
            log_warn "No .env file found, using defaults for local development"
        fi
        
        docker-compose -f docker-compose.yml up -d --build
        
        echo ""
        echo -e "${GREEN}Local deployment complete!${NC}"
        echo -e "  Frontend:  ${GREEN}http://localhost:3001${NC}"
        echo -e "  Backend:   ${GREEN}http://localhost:8081${NC}"
        echo -e "  Database:  ${GREEN}localhost:5433${NC}"
        echo ""
        ;;
    
    *)
        echo "Usage: ./deploy.sh [command]"
        echo ""
        echo "Commands:"
        echo "  deploy    Full deployment: git pull → build → migrate → deploy (default)"
        echo "  quick     Quick deploy without git pull"
        echo "  down      Stop all services"
        echo "  restart   Restart all services"
        echo "  logs      View logs (optionally: ./deploy.sh logs backend)"
        echo "  status    Show status of all services"
        echo "  migrate   Run database migrations only"
        echo "  rebuild   Full rebuild with no cache"
        echo "  clean     Remove everything (containers, volumes, images)"
        echo "  local     Local development with docker-compose.yml"
        echo ""
        echo "Examples:"
        echo "  ./deploy.sh              # Full production deployment"
        echo "  ./deploy.sh quick        # Deploy without git pull"
        echo "  ./deploy.sh logs backend # View backend logs"
        echo "  ./deploy.sh local        # Local development"
        echo ""
        echo "Environment:"
        echo "  The .env file should be placed in the project root:"
        echo "  ${SCRIPT_DIR}/.env"
        echo ""
        echo "  Create from template: cp .env.example .env"
        echo ""
        ;;
esac
