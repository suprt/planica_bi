@echo off
echo Running Planica BI Database Migrations...

cd /d %~dp0
powershell -ExecutionPolicy Bypass -File migrate.ps1 -Environment Development

if %errorlevel% equ 0 (
    echo.
    echo Migration completed successfully!
) else (
    echo.
    echo Migration failed!
    exit /b 1
)

echo.
echo Press any key to exit...
pause > nul