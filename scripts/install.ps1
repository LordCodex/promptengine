# PromptEngine Windows PowerShell Installer Script
# Downloads and installs the correct binary for your architecture.

$Owner = "LordCodex"
$Repo = "promptengine"
$BinaryName = "promptengine"

# Detect architecture
$Arch = "amd64"
if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64" -or $env:PROCESSOR_ARCHITEW6432 -eq "ARM64") {
    $Arch = "arm64"
}

# Fetch latest release
Write-Host "Fetching latest version tag..."
$LatestRelease = Invoke-RestMethod -Uri "https://api.github.com/repos/$Owner/$Repo/releases/latest"
$LatestTag = $LatestRelease.tag_name

if (-not $LatestTag) {
    $LatestTag = "v1.0.0"
}

$Version = $LatestTag.TrimStart('v')

Write-Host "Installing $BinaryName $LatestTag for windows-$Arch..."

# Setup target paths
$InstallDir = Join-Path $env:USERPROFILE ".promptengine\bin"
if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
}

# Build download url
$DownloadUrl = "https://github.com/$Owner/$Repo/releases/download/$LatestTag/${BinaryName}_windows_${Arch}.zip"
$TempZip = Join-Path $env:TEMP "${BinaryName}_windows.zip"
$TempDir = Join-Path $env:TEMP "${BinaryName}_extracted"

# Clean previous temp directories
if (Test-Path $TempDir) { Remove-Item -Recurse -Force $TempDir | Out-Null }

Write-Host "Downloading from: $DownloadUrl"
try {
    Invoke-WebRequest -Uri $DownloadUrl -OutFile $TempZip -UseBasicParsing
} catch {
    Write-Error "Failed to download PromptEngine archive."
    exit 1
}

# Expand and install
Expand-Archive -Path $TempZip -DestinationPath $TempDir -Force

$ExePath = Join-Path $TempDir "${BinaryName}.exe"
if (Test-Path $ExePath) {
    Copy-Item -Path $ExePath -Destination (Join-Path $InstallDir "${BinaryName}.exe") -Force
    Write-Host "✓ Successfully installed ${BinaryName}.exe to $InstallDir!"
    
    # Check if PATH update is needed
    $UserPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    if ($UserPath -notlike "*$InstallDir*") {
        Write-Host "Adding $InstallDir to user PATH environment variable..."
        $NewUserPath = "$UserPath;$InstallDir"
        [Environment]::SetEnvironmentVariable("PATH", $NewUserPath, "User")
        Write-Host "✓ PATH updated! Please restart your terminal/PowerShell session."
    }
} else {
    Write-Error "Binary '${BinaryName}.exe' was not found in extracted archive."
    exit 1
}
