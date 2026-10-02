# AccountTutor 9000 - Multi-Platform Release Packaging (PowerShell)
# Builds standalone native binaries for Windows, macOS (Apple Silicon and Intel), and Linux.

$ErrorActionPreference = "Stop"

$RootDir = Split-Path -Parent $PSScriptRoot
Set-Location $RootDir

Write-Host "=================================================================="
Write-Host " AccountTutor 9000 - Multi-Platform Build and Packaging"
Write-Host "=================================================================="

# 1. Detect Version from cmd/acctg/main.go
$mainContent = Get-Content "$RootDir\cmd\acctg\main.go" -Raw
$versionMatch = [regex]::Match($mainContent, 'const AppVersion = "([^"]+)"')
if ($versionMatch.Success) {
    $Version = $versionMatch.Groups[1].Value
} else {
    throw "Could not determine release version"
}
Write-Host "Target Version: v$Version"

# 2. Automated Quality Gates
Write-Host "`n[1/5] Running automated quality gates..."

Write-Host "  -> gofmt formatting check..."
$gofmtOutput = & gofmt -s -d .
if ($gofmtOutput) {
    Write-Error "gofmt check failed! Please run 'gofmt -s -w .' before releasing."
}

Write-Host "  -> go vet static analysis..."
& go vet ./...
if ($LASTEXITCODE -ne 0) {
    Write-Error "go vet reported errors!"
}

Write-Host "  -> go test full test suite..."
& go test ./...
if ($LASTEXITCODE -ne 0) {
    Write-Error "Unit tests failed!"
}
Write-Host "  All quality gates PASSED."

# 3. Prepare dist/ without deleting existing release artifacts
Write-Host "`n[2/5] Preparing output directory (dist/)..."
$DistDir = "$RootDir\dist"
New-Item -ItemType Directory -Path $DistDir -Force | Out-Null

$SavedGoEnv = @{}
foreach ($name in @("GOOS", "GOARCH", "CGO_ENABLED")) {
    $SavedGoEnv[$name] = [Environment]::GetEnvironmentVariable($name, "Process")
}
try {

$Env:CGO_ENABLED = "0"
$LdFlags = "-s -w"

# 4. Cross-compiling release binaries
Write-Host "`n[3/5] Cross-compiling standalone release binaries..."

# Windows AMD64
Write-Host "  -> Compiling Windows AMD64 (dist/acctg-windows-amd64.exe)..."
$Env:GOOS = "windows"
$Env:GOARCH = "amd64"
& go build "-ldflags=$LdFlags" -o "$DistDir\acctg-windows-amd64.exe" ./cmd/acctg
if ($LASTEXITCODE -ne 0) { throw "Windows AMD64 build failed" }
Copy-Item "$DistDir\acctg-windows-amd64.exe" "$DistDir\acctg.exe" -Force

# macOS ARM64 (Apple Silicon M1-M4)
Write-Host "  -> Compiling macOS Apple Silicon ARM64 (dist/acctg-mac-arm64)..."
$Env:GOOS = "darwin"
$Env:GOARCH = "arm64"
& go build "-ldflags=$LdFlags" -o "$DistDir\acctg-mac-arm64" ./cmd/acctg
if ($LASTEXITCODE -ne 0) { throw "macOS ARM64 build failed" }

# macOS AMD64 (Intel Macs)
Write-Host "  -> Compiling macOS Intel AMD64 (dist/acctg-mac-amd64)..."
$Env:GOOS = "darwin"
$Env:GOARCH = "amd64"
& go build "-ldflags=$LdFlags" -o "$DistDir\acctg-mac-amd64" ./cmd/acctg
if ($LASTEXITCODE -ne 0) { throw "macOS AMD64 build failed" }

# Linux AMD64
Write-Host "  -> Compiling Linux AMD64 (dist/acctg-linux-amd64)..."
$Env:GOOS = "linux"
$Env:GOARCH = "amd64"
& go build "-ldflags=$LdFlags" -o "$DistDir\acctg-linux-amd64" ./cmd/acctg
if ($LASTEXITCODE -ne 0) { throw "Linux AMD64 build failed" }

Write-Host "  -> Compiling Linux ARM64 (dist/acctg-linux-arm64)..."
$Env:GOOS = "linux"
$Env:GOARCH = "arm64"
& go build "-ldflags=$LdFlags" -o "$DistDir\acctg-linux-arm64" ./cmd/acctg
if ($LASTEXITCODE -ne 0) { throw "Linux ARM64 build failed" }

# Reset environment variables
} finally {
    foreach ($name in $SavedGoEnv.Keys) {
        [Environment]::SetEnvironmentVariable($name, $SavedGoEnv[$name], "Process")
    }
}

# Copy macOS & Linux Launchers and desktop integration
Copy-Item "$RootDir\Launch-Tutor.command" "$DistDir\Launch-Tutor.command" -Force
Copy-Item "$RootDir\launch-tutor.sh" "$DistDir\launch-tutor.sh" -Force
Copy-Item "$RootDir\accounttutor.desktop" "$DistDir\accounttutor.desktop" -Force

# 5. Package release bundles
Copy-Item "$RootDir\README.md" "$DistDir\README.md" -Force
Copy-Item "$RootDir\docs\RELEASE-NOTES.md" "$DistDir\RELEASE-NOTES.md" -Force
Write-Host "`n[4/5] Creating distribution archive packages..."

# Windows bundle
$WinZip = "$DistDir\acctg-v$Version-windows-amd64.zip"
Compress-Archive -Path "$DistDir\acctg.exe", "$DistDir\README.md", "$DistDir\RELEASE-NOTES.md", "$RootDir\OVERVIEW.md" -DestinationPath $WinZip -Force
Write-Host "  -> Created $WinZip"

# macOS Apple Silicon bundle
$MacArmZip = "$DistDir\acctg-v$Version-macos-arm64.zip"
Compress-Archive -Path "$DistDir\acctg-mac-arm64", "$DistDir\Launch-Tutor.command", "$DistDir\README.md", "$DistDir\RELEASE-NOTES.md" -DestinationPath $MacArmZip -Force
Write-Host "  -> Created $MacArmZip"

# macOS Intel bundle
$MacIntelZip = "$DistDir\acctg-v$Version-macos-amd64.zip"
Compress-Archive -Path "$DistDir\acctg-mac-amd64", "$DistDir\Launch-Tutor.command", "$DistDir\README.md", "$DistDir\RELEASE-NOTES.md" -DestinationPath $MacIntelZip -Force
Write-Host "  -> Created $MacIntelZip"

# macOS Universal Classmate bundle (contains both architectures + Launch-Tutor.command)
$MacUniversalZip = "$DistDir\acctg-v$Version-macos-classmate-bundle.zip"
Compress-Archive -Path "$DistDir\acctg-mac-arm64", "$DistDir\acctg-mac-amd64", "$DistDir\Launch-Tutor.command", "$DistDir\README.md", "$DistDir\RELEASE-NOTES.md" -DestinationPath $MacUniversalZip -Force
Write-Host "  -> Created $MacUniversalZip"

# Linux AMD64 bundle
$LinuxZip = "$DistDir\acctg-v$Version-linux-amd64.zip"
Compress-Archive -Path "$DistDir\acctg-linux-amd64", "$DistDir\launch-tutor.sh", "$DistDir\accounttutor.desktop", "$DistDir\README.md", "$DistDir\RELEASE-NOTES.md", "$RootDir\OVERVIEW.md" -DestinationPath $LinuxZip -Force
Write-Host "  -> Created $LinuxZip"

# Linux ARM64 bundle
$LinuxArmZip = "$DistDir\acctg-v$Version-linux-arm64.zip"
Compress-Archive -Path "$DistDir\acctg-linux-arm64", "$DistDir\launch-tutor.sh", "$DistDir\accounttutor.desktop", "$DistDir\README.md", "$DistDir\RELEASE-NOTES.md", "$RootDir\OVERVIEW.md" -DestinationPath $LinuxArmZip -Force
Write-Host "  -> Created $LinuxArmZip"

# 6. Generate SHA256 Checksums
Write-Host "`n[5/5] Generating SHA256 checksums..."
$ChecksumFile = "$DistDir\checksums.sha256"
if (Test-Path $ChecksumFile) { Remove-Item $ChecksumFile -Force }

$Artifacts = Get-ChildItem -LiteralPath $DistDir -File | Where-Object { $_.Name -like "acctg-v$Version-*.zip" -or $_.Name -in @("acctg.exe", "acctg-windows-amd64.exe", "acctg-mac-arm64", "acctg-mac-amd64", "acctg-linux-amd64", "acctg-linux-arm64", "Launch-Tutor.command", "launch-tutor.sh", "accounttutor.desktop") } | Sort-Object Name
$Summary = @()

foreach ($item in $Artifacts) {
    $fileHash = (Get-FileHash -Path $item.FullName -Algorithm SHA256).Hash.ToLower()
    "$fileHash  $($item.Name)" | Out-File -FilePath $ChecksumFile -Append -Encoding ascii
    $sizeMB = [math]::Round($item.Length / 1MB, 2)
    $Summary += [PSCustomObject]@{
        Artifact = $item.Name
        "Size (MB)" = $sizeMB
        "SHA256" = $fileHash.Substring(0, 16) + "..."
    }
}

Write-Host "`n=================================================================="
Write-Host " Build & Packaging Complete! Release Artifacts in dist/:"
Write-Host "=================================================================="
$Summary | Format-Table -AutoSize

Write-Host "Checksums written to: dist\checksums.sha256`n"


