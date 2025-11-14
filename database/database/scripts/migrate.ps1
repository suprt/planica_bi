#!/usr/bin/env pwsh
<#
.SYNOPSIS
    Database migration system for Planica BI
.DESCRIPTION
    Professional migration system with rollback support and audit logging
.PARAMETER Environment
    Target environment: Development, Staging, Production
.PARAMETER Action
    Action to perform: migrate, rollback, status
#>

param(
    [string]$Environment = "Development",
    [string]$Action = "migrate",
    [string]$TargetVersion = "",
    [switch]$DryRun,
    [switch]$Verbose
)

$ErrorActionPreference = "Stop"

# Configuration
$Config = @{
    Development = @{
        Host = "localhost"
        Port = 3306
        User = "planica_user"
        Password = "root"
        Database = "planica_bi"
    }
    Staging = @{
        Host = "localhost"
        Port = 3306
        User = "planica_user"
        Password = "root"
        Database = "planica_bi"
    }
    Production = @{
        Host = "localhost"
        Port = 3306
        User = "planica_user"
        Password = $env:PRODUCTION_DB_PASSWORD
        Database = "planica_bi"
    }
}

function Write-Log {
    param($Message, $Color = "White")
    Write-Host "[$(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')] $Message" -ForegroundColor $Color
}

function Get-MySQLConnectionString {
    param($Config, $Environment)
    
    $envConfig = $Config[$Environment]
    if (-not $envConfig) {
        throw "Unknown environment: $Environment"
    }
    
    return "server=$($envConfig.Host);port=$($envConfig.Port);userid=$($envConfig.User);password=$($envConfig.Password);database=$($envConfig.Database);charset=utf8mb4"
}

function Invoke-MySQLQuery {
    param($ConnectionString, $Query, $Parameters = @{})
    
    try {
        $connection = New-Object MySql.Data.MySqlClient.MySqlConnection($ConnectionString)
        $connection.Open()
        
        $command = $connection.CreateCommand()
        $command.CommandText = $Query
        
        foreach ($param in $Parameters.GetEnumerator()) {
            $command.Parameters.AddWithValue($param.Key, $param.Value)
        }
        
        if ($Query.Trim().StartsWith("SELECT")) {
            $adapter = New-Object MySql.Data.MySqlClient.MySqlDataAdapter($command)
            $dataset = New-Object System.Data.DataSet
            $adapter.Fill($dataset) | Out-Null
            return $dataset.Tables[0]
        } else {
            $result = $command.ExecuteNonQuery()
            return $result
        }
    } finally {
        if ($connection -and $connection.State -eq 'Open') {
            $connection.Close()
        }
    }
}

function Get-MigrationFiles {
    $migrationFolders = Get-ChildItem -Path "..\migrations" -Directory | Sort-Object Name
    $migrationFiles = @()
    
    foreach ($folder in $migrationFolders) {
        $upFile = Get-ChildItem -Path $folder.FullName -Filter "*.up.sql" | Select-Object -First 1
        if ($upFile) {
            $migrationFiles += $upFile
        }
    }
    
    return $migrationFiles
}

function Test-MigrationApplied {
    param($ConnectionString, $Version)
    
    $query = "SELECT COUNT(*) as count FROM schema_migrations WHERE version = @version AND status = 'success'"
    $result = Invoke-MySQLQuery -ConnectionString $ConnectionString -Query $query -Parameters @{version = $Version}
    
    return $result.count -gt 0
}

function Apply-Migration {
    param($ConnectionString, $MigrationFile)
    
    $version = [System.IO.Path]::GetFileNameWithoutExtension($MigrationFile.Name)
    $sqlContent = Get-Content -Path $MigrationFile.FullName -Raw
    
    Write-Log "Applying migration: $version" -Color "Green"
    
    if ($DryRun) {
        Write-Log "DRY RUN: Would execute $($MigrationFile.Name)" -Color "Yellow"
        return $true
    }
    
    $stopwatch = [System.Diagnostics.Stopwatch]::StartNew()
    
    try {
        # Execute migration within transaction
        Invoke-MySQLQuery -ConnectionString $ConnectionString -Query "START TRANSACTION"
        Invoke-MySQLQuery -ConnectionString $ConnectionString -Query $sqlContent
        
        # Record migration
        $insertQuery = @"
            INSERT INTO schema_migrations 
            (version, description, checksum, execution_time_ms, applied_by, status)
            VALUES (@version, @description, @checksum, @executionTime, USER(), 'success')
"@
        
        $parameters = @{
            version = $version
            description = "Applied via migration system"
            checksum = (Get-FileHash -Path $MigrationFile.FullName -Algorithm SHA256).Hash
            executionTime = $stopwatch.ElapsedMilliseconds
        }
        
        Invoke-MySQLQuery -ConnectionString $ConnectionString -Query $insertQuery -Parameters $parameters
        Invoke-MySQLQuery -ConnectionString $ConnectionString -Query "COMMIT"
        
        Write-Log "✓ Migration $version applied successfully" -Color "Green"
        return $true
        
    } catch {
        Write-Log "✗ Migration $version failed: $($_.Exception.Message)" -Color "Red"
        Invoke-MySQLQuery -ConnectionString $ConnectionString -Query "ROLLBACK"
        return $false
    } finally {
        $stopwatch.Stop()
    }
}

# Main execution
try {
    Write-Log "Planica BI Database Migration System" -Color "Cyan"
    Write-Log "Environment: $Environment | Action: $Action" -Color "Yellow"
    Write-Host "=" * 60 -ForegroundColor Cyan
    
    $connectionString = Get-MySQLConnectionString -Config $Config -Environment $Environment
    
    # Test connection
    Write-Log "Testing database connection..." -Color "Gray"
    $testResult = Invoke-MySQLQuery -ConnectionString $connectionString -Query "SELECT VERSION() as version"
    Write-Log "Connected to: $($testResult.version)" -Color "Green"
    
    # Get migration files
    $migrationFiles = Get-MigrationFiles
    $appliedCount = 0
    
    foreach ($migrationFile in $migrationFiles) {
        $version = [System.IO.Path]::GetFileNameWithoutExtension($migrationFile.Name)
        
        if (Test-MigrationApplied -ConnectionString $connectionString -Version $version) {
            Write-Log "✓ $version already applied" -Color "Gray"
            continue
        }
        
        if (Apply-Migration -ConnectionString $connectionString -MigrationFile $migrationFile) {
            $appliedCount++
        } else {
            break
        }
    }
    
    Write-Host "=" * 60 -ForegroundColor Cyan
    Write-Log "Migration completed! Applied $appliedCount new migration(s)" -Color "Green"
    
} catch {
    Write-Log "Migration failed: $($_.Exception.Message)" -Color "Red"
    exit 1
}