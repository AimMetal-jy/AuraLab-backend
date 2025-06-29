@echo off
echo Starting AuraLab Backend Services...

echo.
echo Starting WhisperX Flask Service...
start "WhisperX Flask Service" "%~dp0start_flask.bat"

echo Waiting for services to initialize...
timeout /t 5 /nobreak > nul

echo.
echo Starting BlueLM Go Service...
start "BlueLM Go Service" "%~dp0start_go.bat"

echo.
echo Both services are starting in separate windows.
echo WhisperX Flask Service should be available at: http://localhost:5000
echo BlueLM Go Service should be available at: http://localhost:8888
echo.
echo This window will now close.