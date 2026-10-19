# Script PowerShell para detener todos los servicios de XelionCore
Write-Host "Deteniendo servicios de XelionCore..." -ForegroundColor Yellow

$procs = Get-Process -Name LoginService, WorldService, ZoneService -ErrorAction SilentlyContinue
if ($procs) {
    $procs | ForEach-Object {
        Write-Host "Deteniendo $($_.ProcessName) (PID $($_.Id))..."
        $_ | Stop-Process -Force
    }
    Write-Host "Todos los servicios han sido detenidos." -ForegroundColor Green
} else {
    Write-Host "No se encontraron procesos activos de XelionCore." -ForegroundColor Gray
}
