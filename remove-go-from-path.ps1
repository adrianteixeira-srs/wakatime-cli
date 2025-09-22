# PowerShell script to remove C:\Users\devad\go\bin from the user's PATH
$removePath = "C:\Users\devad\go\bin"
$currentPath = [Environment]::GetEnvironmentVariable("Path", "User")

# Split the PATH, remove the target path, and rejoin
$newPath = (($currentPath -split ";") | Where-Object { $_ -ne $removePath }) -join ";"

[Environment]::SetEnvironmentVariable("Path", $newPath, "User")
Write-Host "Path removed from user PATH: $removePath"
Write-Host "Close and reopen your terminal to apply the changes."
