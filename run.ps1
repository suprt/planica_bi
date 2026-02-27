# Planica BI - PowerShell helper script
# Usage: .\run.ps1 <command>

param(
    [Parameter(Position=0)]
    [string]$Command = "help"
)

# Environment file path
$EnvFile = "backend\.env"

# Helper function to run docker compose with env file
function Invoke-DockerCompose {
    param([string]$CommandArgs)
    Write-Host "DEBUG: Running: docker compose --env-file $EnvFile $CommandArgs" -ForegroundColor Yellow
    $command = "docker compose --env-file `"$EnvFile`" $CommandArgs"
    Write-Host "DEBUG: Full command: $command" -ForegroundColor Cyan
    Invoke-Expression $command
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
        Invoke-DockerCompose -CommandArgs "build"
    }
    "build-backend" {
        Invoke-DockerCompose -CommandArgs "build backend"
    }
    "up" {
        Invoke-DockerCompose -CommandArgs "up -d"
        Write-Host "Services started. Waiting for health checks..." -ForegroundColor Yellow
        Start-Sleep -Seconds 3
        Invoke-DockerCompose -CommandArgs "ps"
    }
    "down" {
        Invoke-DockerCompose -CommandArgs "down"
    }
    "restart" {
        Invoke-DockerCompose -CommandArgs "down"
        Invoke-DockerCompose -CommandArgs "up -d"
        Start-Sleep -Seconds 3
        Invoke-DockerCompose -CommandArgs "ps"
    }
    "logs" {
        Invoke-DockerCompose -CommandArgs "logs -f"
    }
    "logs-backend" {
        Invoke-DockerCompose -CommandArgs "logs -f backend"
    }
    "logs-mysql" {
        Invoke-DockerCompose -CommandArgs "logs -f mysql"
    }
    "ps" {
        Invoke-DockerCompose -CommandArgs "ps"
    }
    "status" {
        Invoke-DockerCompose -CommandArgs "ps"
        Write-Host ""
        Write-Host "Service URLs:" -ForegroundColor Cyan
        Write-Host "  Backend API: http://localhost:8080"
        Write-Host "  Adminer:    http://localhost:8081"
        Write-Host "  MySQL:      localhost:3306"
    }
    "clean" {
        Invoke-DockerCompose -CommandArgs "down -v"
        Write-Host "All containers, volumes and networks removed" -ForegroundColor Green
    }
    "rebuild" {
        Invoke-DockerCompose -CommandArgs "down -v"
        Invoke-DockerCompose -CommandArgs "build"
        Invoke-DockerCompose -CommandArgs "up -d"
        Write-Host "Waiting for services to be ready..." -ForegroundColor Yellow
        Start-Sleep -Seconds 5
        Invoke-DockerCompose -CommandArgs "ps"
    }
    "db-reset" {
        Invoke-DockerCompose -CommandArgs "down -v"
        docker volume rm planica_bi_mysql_data -ErrorAction SilentlyContinue
        Invoke-DockerCompose -CommandArgs "up -d mysql"
        Write-Host "Waiting for MySQL to initialize..." -ForegroundColor Yellow
        Start-Sleep -Seconds 10
        Invoke-DockerCompose -CommandArgs "up -d"
    }
    "db-shell" {
        Invoke-DockerCompose -CommandArgs "exec mysql mysql -uroot -p$($env:DB_PASSWORD) --default-character-set=utf8mb4 reports"
    }
    "backend-shell" {
        Invoke-DockerCompose -CommandArgs "exec backend sh"
    }
    "adminer" {
        Write-Host "Opening Adminer at http://localhost:8081" -ForegroundColor Cyan
        Start-Process "http://localhost:8081"
    }
    default {
        Show-Help
    }
}
