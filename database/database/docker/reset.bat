@echo off
echo Resetting Planica BI Database...

docker-compose down -v
docker-compose up -d

echo Database reset complete!
echo Waiting for MySQL to start...
timeout /t 30 /nobreak > nul

echo Press any key to exit...
pause > nul