# Script PowerShell para iniciar los servicios en segundo plano con logs
$rootDir = (Get-Item $PSScriptRoot).Parent.FullName
Set-Location $rootDir

if (!(Test-Path "logs")) {
    New-Item -ItemType Directory -Path "logs" | Out-Null
}

# Detener instancias previas si existen
Get-Process -Name LoginService, WorldService, ZoneService -ErrorAction SilentlyContinue | Stop-Process -Force
Start-Sleep -Milliseconds 500

Write-Host "Iniciando XelionCore..." -ForegroundColor Cyan

# 1. LoginService
Invoke-CimMethod -ClassName Win32_Process -MethodName Create -Arguments @{
    CommandLine = "$rootDir\bin\LoginService.exe"
    CurrentDirectory = $rootDir
} | Out-Null
Write-Host "[1/3] Login Server iniciado en 0.0.0.0:9010 (Logs: logs\LoginService_*.log)" -ForegroundColor Green

# 2. WorldService
Invoke-CimMethod -ClassName Win32_Process -MethodName Create -Arguments @{
    CommandLine = "$rootDir\bin\WorldService.exe"
    CurrentDirectory = $rootDir
} | Out-Null
Write-Host "[2/3] World Server iniciado en 0.0.0.0:9014 (Logs: WorldService_*.log)" -ForegroundColor Green

# 3. ZoneService
Invoke-CimMethod -ClassName Win32_Process -MethodName Create -Arguments @{
    CommandLine = "$rootDir\bin\ZoneService.exe"
    CurrentDirectory = $rootDir
} | Out-Null
Write-Host "[3/3] Zone Server iniciado en 0.0.0.0:9210 (Logs: ZoneService_*.log)" -ForegroundColor Green

Write-Host "`nTodas las instancias estan corriendo en segundo plano." -ForegroundColor Yellow
