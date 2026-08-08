<#
.SYNOPSIS
    Full PMX build and validation script for the picomatch-go repository.

.DESCRIPTION
    Runs the complete Go build/test/lint/benchmark/fuzz pipeline, builds the Wasm target,
    and installs/lints/builds the dashboard.

    This script is designed for Windows PowerShell and can also run in PowerShell Core.

.PARAMETER SkipDashboard
    Skip dashboard install/lint/build steps.

.PARAMETER SkipFuzz
    Skip the fuzz target execution.

.PARAMETER SkipBench
    Skip benchmarks.

.PARAMETER NodeVersion
    Expected Node.js version for the dashboard build.

.PARAMETER GoVersion
    Expected Go version for the project.

.EXAMPLE
    .\build.ps1

.EXAMPLE
    .\build.ps1 -SkipFuzz -SkipBench
#>

[CmdletBinding()]
param(
    [switch]$SkipDashboard,
    [switch]$SkipFuzz,
    [switch]$SkipBench,
    [string]$NodeVersion = '22',
    [string]$GoVersion = '1.21'
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

function Write-Heading {
    param([string]$Text)
    Write-Host "`n===== $Text =====" -ForegroundColor Cyan
}

function Run-Command {
    param(
        [string]$Message,
        [string]$Command,
        [string[]]$Args = @(),
        [string]$WorkingDirectory = $PWD
    )

    Write-Heading $Message
    Push-Location $WorkingDirectory
    try {
        & $Command @Args
        if ($LASTEXITCODE -ne 0) {
            throw "Command '$Command' failed with exit code $LASTEXITCODE."
        }
    } finally {
        Pop-Location
    }
}

function Assert-Tool {
    param([string]$Tool)
    if (-not (Get-Command $Tool -ErrorAction SilentlyContinue)) {
        throw "Required tool '$Tool' was not found in PATH."
    }
}

Write-Heading "PMX Build Script"
Write-Host "Repository root: $PWD"
Write-Host "Go version: $GoVersion"
Write-Host "Node version: $NodeVersion"

Assert-Tool 'go'
Assert-Tool 'npm'

# Validate git clean state only for formatting verification.
Write-Heading "Preflight checks"
if (-not (Get-Command git -ErrorAction SilentlyContinue)) {
    throw 'git is required for this build script.'
}

# Core Go workflow
Run-Command -Message 'Formatting Go code' -Command 'go' -Args @('fmt', './...')
Run-Command -Message 'Verifying formatting' -Command 'git' -Args @('diff', '--exit-code', '--', '*.go')
Run-Command -Message 'Running golangci-lint' -Command 'go' -Args @('install', 'github.com/golangci/golangci-lint/cmd/golangci-lint@v1.60.3')
Run-Command -Message 'Executing lint step' -Command (Join-Path (go env GOPATH) 'bin\golangci-lint.exe') -Args @('run', './...')
Run-Command -Message 'Running go vet' -Command 'go' -Args @('vet', './...')
Run-Command -Message 'Running unit tests' -Command 'go' -Args @('test', '-v', './...')
Run-Command -Message 'Building CLI binary' -Command 'go' -Args @('build', '-o', 'pmx', '.\cmd\pmx')

if (-not $SkipFuzz) {
    Run-Command -Message 'Running fuzz targets' -Command 'go' -Args @('test', '-fuzz=Fuzz', '-fuzztime=30s', './...')
} else {
    Write-Host 'Skipping fuzz targets.' -ForegroundColor Yellow
}

if (-not $SkipBench) {
    Run-Command -Message 'Running benchmark suite' -Command 'go' -Args @('test', '-bench=.', '-benchmem', '-run=^$', './...')
} else {
    Write-Host 'Skipping benchmarks.' -ForegroundColor Yellow
}

Run-Command -Message 'Building WebAssembly target' -Command 'go' -Args @('build', '-o', 'dashboard/public/picomatch.wasm', 'cmd/wasm/main.go')

if (-not $SkipDashboard) {
    Write-Heading 'Dashboard build'
    $dashboardDir = Join-Path $PWD 'dashboard'

    Run-Command -Message 'Installing dashboard dependencies' -Command 'npm' -Args @('ci') -WorkingDirectory $dashboardDir
    Run-Command -Message 'Linting dashboard' -Command 'npm' -Args @('run', 'lint') -WorkingDirectory $dashboardDir
    Run-Command -Message 'Building dashboard' -Command 'npm' -Args @('run', 'build') -WorkingDirectory $dashboardDir
} else {
    Write-Host 'Skipping dashboard build.' -ForegroundColor Yellow
}

Write-Heading 'Build complete'
Write-Host 'All requested PMX build and validation tasks finished successfully.' -ForegroundColor Green
