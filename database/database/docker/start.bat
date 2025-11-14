@echo off
echo Starting Planica BI Database with MySQL 9.5...

docker-compose up -d

echo Waiting for MySQL to start...
timeout /t 30 /nobreak > nul

echo Database started!
echo MySQL 9.5 is available on localhost:3306
echo Adminer is available on http://localhost:8080
pause