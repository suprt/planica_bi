# Planica BI - PowerShell helper script
# Usage: .\run.ps1 <command>

param(
    [Parameter(Position=0)]
    [string]$Command = "help"
)

function Show-Help {
    Write-Host "Planica BI - PowerShell Commands" -ForegroundColor Cyan
    Write-Host "================================" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Basic commands:"
    Write-Host "  .\run.ps1 build         - Build all Docker images"
    Write-Host "  .\run.ps1 build-backend - Build only backend image"
    Write-Host "  .\run.ps1 up            - Start all services"
    Write-Host "  .\run.ps1 down          - Stop all services"
    Write-Host "  .\run.ps1 restart       - Restart all services"
    Write-Host "  .\run.ps1 ps            - Show running containers"
    Write-Host "  .\run.ps1 status        - Show detailed status"
    Write-Host ""
    Write-Host "Logs:"
    Write-Host "  .\run.ps1 logs          - Show logs from all services"
    Write-Host "  .\run.ps1 logs-backend   - Show logs from backend"
    Write-Host "  .\run.ps1 logs-mysql     - Show logs from MySQL"
    Write-Host ""
    Write-Host "Database:"
    Write-Host "  .\run.ps1 db-shell      - Open MySQL shell"
    Write-Host "  .\run.ps1 db-reset      - Reset database (drop volume and recreate)"
    Write-Host ""
    Write-Host "Development:"
    Write-Host "  .\run.ps1 backend-shell - Open backend container shell"
    Write-Host "  .\run.ps1 adminer       - Open Adminer in browser"
    Write-Host ""
    Write-Host "Cleanup:"
    Write-Host "  .\run.ps1 clean         - Stop and remove containers, volumes, networks"
    Write-Host "  .\run.ps1 rebuild       - Rebuild everything from scratch"
    Write-Host ""
}

switch ($Command) {
    "build" {
        docker-compose build
    }
    "build-backend" {
        docker-compose build backend
    }
    "up" {
        docker-compose up -d
        Write-Host "Services started. Waiting for health checks..." -ForegroundColor Yellow
        Start-Sleep -Seconds 3
        docker-compose ps
    }
    "down" {
        docker-compose down
    }
    "restart" {
        docker-compose down
        docker-compose up -d
        Start-Sleep -Seconds 3
        docker-compose ps
    }
    "logs" {
        docker-compose logs -f
    }
    "logs-backend" {
        docker-compose logs -f backend
    }
    "logs-mysql" {
        docker-compose logs -f mysql
    }
    "ps" {
        docker-compose ps
    }
    "status" {
        docker-compose ps
        Write-Host ""
        Write-Host "Service URLs:" -ForegroundColor Cyan
        Write-Host "  Backend API: http://localhost:8080"
        Write-Host "  Adminer:    http://localhost:8081"
        Write-Host "  MySQL:      localhost:3306"
    }
    "clean" {
        docker-compose down -v
        Write-Host "All containers, volumes and networks removed" -ForegroundColor Green
    }
    "rebuild" {
        docker-compose down -v
        docker-compose build
        docker-compose up -d
        Write-Host "Waiting for services to be ready..." -ForegroundColor Yellow
        Start-Sleep -Seconds 5
        docker-compose ps
    }
    "db-reset" {
        docker-compose down -v
        docker volume rm planica_bi_mysql_data -ErrorAction SilentlyContinue
        docker-compose up -d mysql
        Write-Host "Waiting for MySQL to initialize..." -ForegroundColor Yellow
        Start-Sleep -Seconds 10
        docker-compose up -d
    }
    "db-shell" {
        docker-compose exec mysql mysql -uroot -p1234 --default-character-set=utf8mb4 reports
    }
    "backend-shell" {
        docker-compose exec backend sh
    }
    "adminer" {
        Write-Host "Opening Adminer at http://localhost:8081" -ForegroundColor Cyan
        Start-Process "http://localhost:8081"
    }
    default {
        Show-Help
    }
}
