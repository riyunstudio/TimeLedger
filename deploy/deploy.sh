#!/bin/bash

# TimeLedger Deployment Script for VPS

set -e

echo "🚀 Starting TimeLedger deployment..."

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Configuration
COMPOSE_FILE="docker-compose.yml"
ENV_FILE=".env.production"

# Check if .env.production exists
if [ ! -f "$ENV_FILE" ]; then
    echo -e "${RED}Error: $ENV_FILE not found!${NC}"
    echo "Please copy .env.production.example to $ENV_FILE and update the values."
    exit 1
fi

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo -e "${RED}Error: Docker is not installed!${NC}"
    echo "Please install Docker first."
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}Error: Docker Compose is not installed!${NC}"
    echo "Please install Docker Compose first."
    exit 1
fi

# Pull latest code (optional)
echo -e "${YELLOW}Pulling latest code...${NC}"
# git pull origin main

# Build and start containers
echo -e "${YELLOW}Building Docker images...${NC}"
docker-compose -f $COMPOSE_FILE --env-file $ENV_FILE build

echo -e "${YELLOW}Starting containers...${NC}"
docker-compose -f $COMPOSE_FILE --env-file $ENV_FILE up -d

# Wait for services to be healthy
echo -e "${YELLOW}Waiting for services to be healthy...${NC}"
sleep 30

# Check container status
echo -e "${GREEN}Container status:${NC}"
docker-compose -f $COMPOSE_FILE --env-file $ENV_FILE ps

# Run database migrations (if needed)
echo -e "${YELLOW}Running database migrations...${NC}"
# docker-compose -f $COMPOSE_FILE --env-file $ENV_FILE exec -T app go run main.go migrate

echo -e "${GREEN}✅ Deployment completed successfully!${NC}"
echo ""
echo "Access the application at:"
echo "  - API: http://your-server-ip:8080"
echo "  - Swagger UI: http://your-server-ip:8080/swagger/index.html"
echo "  - RabbitMQ Management: http://your-server-ip:15672"
echo ""
echo "To view logs:"
echo "  docker-compose -f $COMPOSE_FILE --env-file $ENV_FILE logs -f app"
