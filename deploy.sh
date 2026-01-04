#!/bin/bash

# HeySpoilMe - Production Deployment Script
# This script builds and deploys all services using Docker Compose

set -e

echo "🚀 HeySpoilMe Production Deployment"
echo "===================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default values
ACTION=${1:-"up"}
BUILD=${2:-"--build"}

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

case $ACTION in
    "up")
        echo -e "${YELLOW}Building and starting all services...${NC}"
        docker-compose up -d $BUILD
        
        echo ""
        echo -e "${GREEN}===================================="
        echo -e "Deployment Complete!"
        echo -e "====================================${NC}"
        echo ""
        echo -e "  Frontend:  ${GREEN}http://localhost:3001${NC}"
        echo -e "  Backend:   ${GREEN}http://localhost:8081${NC}"
        echo -e "  Database:  ${GREEN}localhost:5433${NC}"
        echo ""
        echo -e "${YELLOW}View logs: docker-compose logs -f${NC}"
        ;;
    
    "down")
        echo -e "${YELLOW}Stopping all services...${NC}"
        docker-compose down
        echo -e "${GREEN}All services stopped.${NC}"
        ;;
    
    "restart")
        echo -e "${YELLOW}Restarting all services...${NC}"
        docker-compose down
        docker-compose up -d $BUILD
        echo -e "${GREEN}All services restarted.${NC}"
        ;;
    
    "logs")
        docker-compose logs -f
        ;;
    
    "status")
        echo -e "${YELLOW}Service Status:${NC}"
        docker-compose ps
        ;;
    
    "rebuild")
        echo -e "${YELLOW}Rebuilding all services (no cache)...${NC}"
        docker-compose build --no-cache
        docker-compose up -d
        echo -e "${GREEN}Rebuild complete.${NC}"
        ;;
    
    "clean")
        echo -e "${YELLOW}Stopping services and removing volumes...${NC}"
        docker-compose down -v
        echo -e "${GREEN}Cleanup complete.${NC}"
        ;;
    
    *)
        echo "Usage: ./deploy.sh [command]"
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
        ;;
esac

