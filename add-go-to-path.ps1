# PowerShell script to add Go to the user's PATH
Write-Host "===============================" -ForegroundColor Cyan
Write-Host "  Script: Add Go to User PATH" -ForegroundColor Cyan
Write-Host "===============================" -ForegroundColor Cyan

$possiblePaths = @("C:\Go\bin", "$env:USERPROFILE\go\bin", "C:\Program Files\Go\bin")
$goPath = $null
foreach ($path in $possiblePaths) {
    if (Test-Path (Join-Path $path 'go.exe')) {
        $goPath = $path
        break
    }
}
Write-Host "Detected Go path: $goPath" -ForegroundColor Yellow

if (-not $goPath) {
    Write-Host "[ERROR] go.exe not found in standard paths." -ForegroundColor Red
    Write-Host "Check your Go installation and try again." -ForegroundColor Yellow
    exit 1
}

$currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($currentPath -notlike "*$goPath*") {
    [Environment]::SetEnvironmentVariable("Path", "$currentPath;$goPath", "User")
    Write-Host "[SUCCESS] Go added to user PATH:" -ForegroundColor Green
    Write-Host "          $goPath" -ForegroundColor White
    Write-Host "`nClose and reopen your terminal to apply the changes." -ForegroundColor Yellow
} else {
    Write-Host "[INFO] Go is already in the user PATH:" -ForegroundColor Green
    Write-Host "       $goPath" -ForegroundColor White
}
Write-Host "===============================" -ForegroundColor Cyan