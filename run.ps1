# Planica BI - PowerShell helper script
# Usage: .\run.ps1 <command>

param(
    [Parameter(Position=0)]
    [string]$Command = "help"
)

# Environment file path
$EnvFile = "backend\.env"

# Helper function to run docker-compose with env file
function Invoke-DockerCompose {
    param([string]$Args)
    docker-compose --env-file $EnvFile $Args
}

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
        Invoke-DockerCompose -Args "build"
    }
    "build-backend" {
        Invoke-DockerCompose -Args "build backend"
    }
    "up" {
        Invoke-DockerCompose -Args "up -d"
        Write-Host "Services started. Waiting for health checks..." -ForegroundColor Yellow
        Start-Sleep -Seconds 3
        Invoke-DockerCompose -Args "ps"
    }
    "down" {
        Invoke-DockerCompose -Args "down"
    }
    "restart" {
        Invoke-DockerCompose -Args "down"
        Invoke-DockerCompose -Args "up -d"
        Start-Sleep -Seconds 3
        Invoke-DockerCompose -Args "ps"
    }
    "logs" {
        Invoke-DockerCompose -Args "logs -f"
    }
    "logs-backend" {
        Invoke-DockerCompose -Args "logs -f backend"
    }
    "logs-mysql" {
        Invoke-DockerCompose -Args "logs -f mysql"
    }
    "ps" {
        Invoke-DockerCompose -Args "ps"
    }
    "status" {
        Invoke-DockerCompose -Args "ps"
        Write-Host ""
        Write-Host "Service URLs:" -ForegroundColor Cyan
        Write-Host "  Backend API: http://localhost:8080"
        Write-Host "  Adminer:    http://localhost:8081"
        Write-Host "  MySQL:      localhost:3306"
    }
    "clean" {
        Invoke-DockerCompose -Args "down -v"
        Write-Host "All containers, volumes and networks removed" -ForegroundColor Green
    }
    "rebuild" {
        Invoke-DockerCompose -Args "down -v"
        Invoke-DockerCompose -Args "build"
        Invoke-DockerCompose -Args "up -d"
        Write-Host "Waiting for services to be ready..." -ForegroundColor Yellow
        Start-Sleep -Seconds 5
        Invoke-DockerCompose -Args "ps"
    }
    "db-reset" {
        Invoke-DockerCompose -Args "down -v"
        docker volume rm planica_bi_mysql_data -ErrorAction SilentlyContinue
        Invoke-DockerCompose -Args "up -d mysql"
        Write-Host "Waiting for MySQL to initialize..." -ForegroundColor Yellow
        Start-Sleep -Seconds 10
        Invoke-DockerCompose -Args "up -d"
    }
    "db-shell" {
        Invoke-DockerCompose -Args "exec mysql mysql -uroot -p$($env:DB_PASSWORD) --default-character-set=utf8mb4 reports"
    }
    "backend-shell" {
        Invoke-DockerCompose -Args "exec backend sh"
    }
    "adminer" {
        Write-Host "Opening Adminer at http://localhost:8081" -ForegroundColor Cyan
        Start-Process "http://localhost:8081"
    }
    default {
        Show-Help
    }
}
