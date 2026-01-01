# mcinit installation script for Windows

$ErrorActionPreference = "Stop"

Write-Host "mcinit Installation Script" -ForegroundColor Green
Write-Host "===========================" -ForegroundColor Green

# Detect architecture
$Arch = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }

# Get latest release version
Write-Host "Fetching latest release..." -ForegroundColor Yellow
$LatestRelease = Invoke-RestMethod -Uri "https://api.github.com/repos/jackh54/mcinit/releases/latest"
$LatestVersion = $LatestRelease.tag_name

if (-not $LatestVersion) {
    Write-Host "Failed to fetch latest version" -ForegroundColor Red
    exit 1
}

Write-Host "Latest version: $LatestVersion" -ForegroundColor Green

# Download URL
$DownloadUrl = "https://github.com/jackh54/mcinit/releases/download/$LatestVersion/mcinit_$($LatestVersion.TrimStart('v'))_windows_$Arch.zip"

# Download
$TempZip = "$env:TEMP\mcinit.zip"
$TempDir = "$env:TEMP\mcinit"

Write-Host "Downloading mcinit..." -ForegroundColor Yellow
Invoke-WebRequest -Uri $DownloadUrl -OutFile $TempZip

# Extract
Write-Host "Extracting..." -ForegroundColor Yellow
if (Test-Path $TempDir) {
    Remove-Item -Path $TempDir -Recurse -Force
}
Expand-Archive -Path $TempZip -DestinationPath $TempDir -Force

# Install
$InstallDir = "$env:LOCALAPPDATA\mcinit"
if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir | Out-Null
}

Write-Host "Installing to $InstallDir..." -ForegroundColor Yellow
Move-Item -Path "$TempDir\mcinit.exe" -Destination "$InstallDir\mcinit.exe" -Force

# Add to PATH if not already there
$CurrentPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($CurrentPath -notlike "*$InstallDir*") {
    Write-Host "Adding to PATH..." -ForegroundColor Yellow
    [Environment]::SetEnvironmentVariable(
        "Path",
        "$CurrentPath;$InstallDir",
        "User"
    )
    Write-Host "PATH updated. Please restart your terminal." -ForegroundColor Yellow
}

# Cleanup
Remove-Item -Path $TempZip -Force
Remove-Item -Path $TempDir -Recurse -Force

Write-Host ""
Write-Host "mcinit installed successfully!" -ForegroundColor Green
Write-Host "Run 'mcinit --help' to get started" -ForegroundColor Green
Write-Host ""
Write-Host "Note: You may need to restart your terminal for PATH changes to take effect." -ForegroundColor Yellow

