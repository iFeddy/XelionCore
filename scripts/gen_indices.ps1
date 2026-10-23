# Genera docs/Indices/*.md (firmas func/type con n° de línea) para todo el proyecto.
# Uso (desde la raíz del proyecto): powershell -File scripts/gen_indices.ps1
$ErrorActionPreference = 'Stop'
$root = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
Set-Location $root
New-Item -Force -ItemType Directory docs/Indices | Out-Null

function New-Index($name, $title, $dirs, $note) {
    $o = @("# Índice: $title", "",
        "> Generado automáticamente por scripts/gen_indices.ps1. Solo firmas con n° de línea: abrir el archivo con StartLine/EndLine. No reemplaza a docs/Servicios/. $note", "")
    foreach ($d in $dirs) {
        if (-not (Test-Path $d)) { continue }
        Get-ChildItem $d -Recurse -Filter *.go | Sort-Object FullName | ForEach-Object {
            $rel = $_.FullName.Substring($root.Length + 1).Replace('\', '/')
            $o += "## $rel"
            $i = 0
            foreach ($l in Get-Content $_.FullName) {
                $i++
                if ($l -match '^(func|type)\s') {
                    $s = ($l -replace '\s*\{\s*$', '')
                    if ($s.Length -gt 140) { $s = $s.Substring(0, 140) }
                    $o += "- L${i}: ``$s``"
                }
            }
            $o += ""
        }
    }
    Set-Content "docs/Indices/$name.md" $o -Encoding UTF8
    Write-Host ("{0}.md {1} KB" -f $name, [int]((Get-Item "docs/Indices/$name.md").Length / 1KB))
}

New-Index 'Data'      'pkg/data (cargadores de datos)'  @('pkg/data', 'pkg/script') 'Datos JSON en data/.'
New-Index 'Zone'      'ZoneService'                     @('pkg/server/zone', 'cmd/ZoneService') 'Protocolos en pkg/protocol.'
New-Index 'Character' 'WorldService / Character'        @('pkg/server/world', 'pkg/database', 'cmd/WorldService') 'Equivale a WorldManager+Character.'
New-Index 'Protocol'  'pkg/protocol'                    @('pkg/protocol') ''
New-Index 'Login'     'LoginService'                    @('pkg/server/login', 'cmd/LoginService') ''
New-Index 'Core'      'config/crypto/logger/network'    @('pkg/config', 'pkg/crypto', 'pkg/logger', 'pkg/network') ''
