.PHONY: help build build-backend up down restart logs clean rebuild test db-shell backend-shell ps status adminer

# Default target
help:
	@echo "Planica BI - Makefile Commands"
	@echo "=============================="
	@echo ""
	@echo "Basic commands:"
	@echo "  make build         - Build all Docker images"
	@echo "  make build-backend - Build only backend image"
	@echo "  make up            - Start all services"
	@echo "  make down          - Stop all services"
	@echo "  make restart       - Restart all services"
	@echo "  make ps            - Show running containers"
	@echo "  make status        - Show detailed status"
	@echo ""
	@echo "Logs:"
	@echo "  make logs          - Show logs from all services"
	@echo "  make logs-backend  - Show logs from backend"
	@echo "  make logs-mysql    - Show logs from MySQL"
	@echo ""
	@echo "Database:"
	@echo "  make db-shell      - Open MySQL shell"
	@echo "  make db-reset      - Reset database (drop volume and recreate)"
	@echo ""
	@echo "Development:"
	@echo "  make backend-shell - Open backend container shell"
	@echo "  make adminer       - Open Adminer in browser (http://localhost:8081)"
	@echo ""
	@echo "Cleanup:"
	@echo "  make clean         - Stop and remove containers, volumes, networks"
	@echo "  make rebuild       - Rebuild everything from scratch"
	@echo ""

# Build Docker images
build:
	docker-compose build

# Build only backend
build-backend:
	docker-compose build backend

# Start all services
up:
	docker-compose up -d
	@echo "Services started. Waiting for health checks..."
	@sleep 3
	@docker-compose ps

# Stop all services
down:
	docker-compose down

# Restart all services
restart: down up

# Show logs
logs:
	docker-compose logs -f

# Show backend logs
logs-backend:
	docker-compose logs -f backend

# Show MySQL logs
logs-mysql:
	docker-compose logs -f mysql

# Show running containers
ps:
	docker-compose ps

# Show detailed status
status: ps
	@echo ""
	@echo "Service URLs:"
	@echo "  Backend API: http://localhost:8080"
	@echo "  Adminer:    http://localhost:8081"
	@echo "  MySQL:      localhost:3306"

# Clean everything (containers, volumes, networks)
clean:
	docker-compose down -v
	@echo "All containers, volumes and networks removed"

# Rebuild everything from scratch
rebuild: clean build up
	@echo "Waiting for services to be ready..."
	@sleep 5
	@docker-compose ps

# Reset database (drop volume and recreate)
db-reset:
	docker-compose down -v
	docker volume rm planica_bi_mysql_data 2>/dev/null || true
	docker-compose up -d mysql
	@echo "Waiting for MySQL to initialize..."
	@sleep 10
	@docker-compose up -d

# Open MySQL shell
db-shell:
	docker-compose exec mysql mysql -uroot -p1234 --default-character-set=utf8mb4 reports

# Open backend container shell
backend-shell:
	docker-compose exec backend sh

# Open Adminer in browser (cross-platform)
adminer:
	@echo "Opening Adminer at http://localhost:8081"
	@if command -v xdg-open > /dev/null; then \
		xdg-open http://localhost:8081; \
	elif command -v open > /dev/null; then \
		open http://localhost:8081; \
	elif command -v start > /dev/null; then \
		start http://localhost:8081; \
	else \
		echo "Please open http://localhost:8081 in your browser"; \
	fi
