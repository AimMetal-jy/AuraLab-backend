@echo off
echo === AuraLab Backend Configuration Checker ===
echo.
echo Running configuration check...
echo.
go run check_config.go
echo.
echo Press any key to exit...
pause >nul