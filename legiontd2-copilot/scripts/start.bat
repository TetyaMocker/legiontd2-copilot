@echo off
cd /d "%~dp0.."
echo Starting LT2 Copilot...
start "LT2 Copilot" bin\lt2-copilot.exe
echo Web UI: http://localhost:8080
echo Press any key to stop...
pause
taskkill /f /im lt2-copilot.exe >nul 2>&1
