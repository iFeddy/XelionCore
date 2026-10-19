# Script PowerShell para iniciar los servicios en segundo plano con logs
$rootDir = (Get-Item $PSScriptRoot).Parent.FullName
Set-Location $rootDir

if (!(Test-Path "logs")) {
    New-Item -ItemType Directory -Path "logs" | Out-Null
}

# Detener instancias previas si existen
Get-Process -Name LoginService, WorldService, ZoneService -ErrorAction SilentlyContinue | Stop-Process -Force

Write-Host "Iniciando XelionCore..." -ForegroundColor Cyan

# 1. LoginService
Start-Process -FilePath "cmd.exe" -ArgumentList "/c bin\LoginService.exe > logs\login.log 2>&1" -WorkingDirectory $rootDir -WindowStyle Hidden
Write-Host "[1/3] Login Server iniciado en 0.0.0.0:9010 (Logs: logs\login.log)" -ForegroundColor Green

# 2. WorldService
Start-Process -FilePath "cmd.exe" -ArgumentList "/c bin\WorldService.exe > logs\world.log 2>&1" -WorkingDirectory $rootDir -WindowStyle Hidden
Write-Host "[2/3] World Server iniciado en 0.0.0.0:9014 (Logs: logs\world.log)" -ForegroundColor Green

# 3. ZoneService
Start-Process -FilePath "cmd.exe" -ArgumentList "/c bin\ZoneService.exe > logs\zone.log 2>&1" -WorkingDirectory $rootDir -WindowStyle Hidden
Write-Host "[3/3] Zone Server iniciado en 0.0.0.0:9210 (Logs: logs\zone.log)" -ForegroundColor Green

Write-Host "`nTodas las instancias estn corriendo en segundo plano." -ForegroundColor Yellow
