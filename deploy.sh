#!/bin/bash

# HeySpoilMe - Production Deployment Script
# This script builds and deploys all services using Docker Compose

set -e

echo "HeySpoilMe Production Deployment"
echo "===================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default values
ACTION=${1:-"up"}
MODE=${2:-"local"}

# Check if docker is running
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}Error: Docker is not running. Please start Docker first.${NC}"
    exit 1
fi

# Check if docker-compose is available
if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}Error: docker-compose is not installed.${NC}"
    exit 1
fi

# Select compose file based on mode
if [ "$MODE" = "prod" ]; then
    COMPOSE_FILE="docker-compose.prod.yml"
    echo -e "${YELLOW}Using production config with Caddy...${NC}"
else
    COMPOSE_FILE="docker-compose.yml"
fi

case $ACTION in
    "up")
        echo -e "${YELLOW}Building and starting all services...${NC}"
        docker-compose -f $COMPOSE_FILE up -d --build
        
        echo ""
        echo -e "${GREEN}===================================="
        echo -e "Deployment Complete!"
        echo -e "====================================${NC}"
        echo ""
        if [ "$MODE" = "prod" ]; then
            echo -e "  Frontend:  ${GREEN}https://heyspoil.me${NC}"
            echo -e "  Backend:   ${GREEN}https://api.heyspoil.me${NC}"
        else
            echo -e "  Frontend:  ${GREEN}http://localhost:3001${NC}"
            echo -e "  Backend:   ${GREEN}http://localhost:8081${NC}"
            echo -e "  Database:  ${GREEN}localhost:5433${NC}"
        fi
        echo ""
        echo -e "${YELLOW}View logs: docker-compose -f $COMPOSE_FILE logs -f${NC}"
        ;;
    
    "down")
        echo -e "${YELLOW}Stopping all services...${NC}"
        docker-compose -f $COMPOSE_FILE down
        echo -e "${GREEN}All services stopped.${NC}"
        ;;
    
    "restart")
        echo -e "${YELLOW}Restarting all services...${NC}"
        docker-compose -f $COMPOSE_FILE down
        docker-compose -f $COMPOSE_FILE up -d --build
        echo -e "${GREEN}All services restarted.${NC}"
        ;;
    
    "logs")
        docker-compose -f $COMPOSE_FILE logs -f
        ;;
    
    "status")
        echo -e "${YELLOW}Service Status:${NC}"
        docker-compose -f $COMPOSE_FILE ps
        ;;
    
    "rebuild")
        echo -e "${YELLOW}Rebuilding all services (no cache)...${NC}"
        docker-compose -f $COMPOSE_FILE build --no-cache
        docker-compose -f $COMPOSE_FILE up -d
        echo -e "${GREEN}Rebuild complete.${NC}"
        ;;
    
    "clean")
        echo -e "${YELLOW}Stopping services and removing volumes...${NC}"
        docker-compose -f $COMPOSE_FILE down -v
        echo -e "${GREEN}Cleanup complete.${NC}"
        ;;
    
    *)
        echo "Usage: ./deploy.sh [command] [mode]"
        echo ""
        echo "Commands:"
        echo "  up        Build and start all services (default)"
        echo "  down      Stop all services"
        echo "  restart   Restart all services"
        echo "  logs      View logs from all services"
        echo "  status    Show status of all services"
        echo "  rebuild   Rebuild all services without cache"
        echo "  clean     Stop services and remove volumes"
        echo ""
        echo "Modes:"
        echo "  local     Use docker-compose.yml (default)"
        echo "  prod      Use docker-compose.prod.yml with Caddy"
        echo ""
        echo "Examples:"
        echo "  ./deploy.sh up local    # Local deployment"
        echo "  ./deploy.sh up prod     # Production with Caddy"
        echo ""
        ;;
esac

