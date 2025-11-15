# Скрипт для тестирования очередей
$baseUrl = "http://localhost:8080/api"
$ErrorActionPreference = "Continue"

Write-Host "=== Тестирование очередей ===" -ForegroundColor Cyan
Write-Host ""

# 1. Вход под админом
Write-Host "1. Вход под админом..." -ForegroundColor Yellow
$loginBody = @{
    email = "admin_repo_test@example.com"
    password = "password123"
} | ConvertTo-Json

try {
    $loginResponse = Invoke-RestMethod -Uri "$baseUrl/auth/login" -Method POST -Body $loginBody -ContentType "application/json"
    $token = $loginResponse.token
    $headers = @{ "Authorization" = "Bearer $token" }
    Write-Host "   OK: Вход выполнен" -ForegroundColor Green
} catch {
    Write-Host "   ERROR: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# 2. Получение проектов
Write-Host ""
Write-Host "2. Получение проектов..." -ForegroundColor Yellow
try {
    $projects = Invoke-RestMethod -Uri "$baseUrl/projects" -Method GET -Headers $headers
    if ($projects.Count -eq 0) {
        Write-Host "   WARNING: Нет проектов в базе" -ForegroundColor Yellow
        exit 1
    }
    Write-Host "   OK: Получено проектов: $($projects.Count)" -ForegroundColor Green
    $projectId = $projects[0].id
    Write-Host "   Используем проект ID: $projectId" -ForegroundColor Gray
} catch {
    Write-Host "   ERROR: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# 3. Постановка задачи синхронизации в очередь
Write-Host ""
Write-Host "3. Постановка задачи синхронизации в очередь (POST /api/sync/$projectId)..." -ForegroundColor Yellow
try {
    $syncResponse = Invoke-RestMethod -Uri "$baseUrl/sync/$projectId" -Method POST -Headers $headers
    Write-Host "   OK: Задача поставлена в очередь!" -ForegroundColor Green
    Write-Host "      Message: $($syncResponse.message)" -ForegroundColor Gray
    Write-Host "      Project ID: $($syncResponse.project_id)" -ForegroundColor Gray
    Write-Host "      Task ID: $($syncResponse.task_id)" -ForegroundColor Gray
    Write-Host "      Queue: $($syncResponse.queue)" -ForegroundColor Gray
} catch {
    Write-Host "   ERROR: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.Exception.Response) {
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        $responseBody = $reader.ReadToEnd()
        Write-Host "   Response: $responseBody" -ForegroundColor Red
    }
    exit 1
}

# 4. Проверка логов воркера
Write-Host ""
Write-Host "4. Проверка логов воркера (подождите 5 секунд)..." -ForegroundColor Yellow
Start-Sleep -Seconds 5

Write-Host ""
Write-Host "=== Тестирование завершено ===" -ForegroundColor Cyan
Write-Host ""
Write-Host "Проверьте логи backend для подтверждения обработки задачи:" -ForegroundColor Yellow
Write-Host "  docker-compose logs backend --tail 20" -ForegroundColor Gray

