@echo off
echo Generating Swagger documentation...
swag init
if %ERRORLEVEL% neq 0 (
    echo Failed to generate Swagger documentation.
    pause
    exit /b %ERRORLEVEL%
)

echo Documentation generated successfully!
echo Starting the server...
go run main.go