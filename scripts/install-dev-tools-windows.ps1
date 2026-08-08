# Install Go and Node.js (LTS) on Windows using winget or Chocolatey.
# Run in an elevated PowerShell session (Run as Administrator) for installs.

function Write-Info($m){ Write-Host "[INFO] $m" }
function Write-Warn($m){ Write-Host "[WARN] $m" -ForegroundColor Yellow }
function Write-Err($m){ Write-Host "[ERROR] $m" -ForegroundColor Red }

Write-Info "Checking for package managers (winget, choco)..."
$hasWinget = (Get-Command winget -ErrorAction SilentlyContinue) -ne $null
$hasChoco  = (Get-Command choco  -ErrorAction SilentlyContinue) -ne $null

if ($hasWinget) {
    Write-Info "Using winget to install Node.js (LTS) and Go."
    winget install --id OpenJS.NodeJS.LTS -e --accept-package-agreements --accept-source-agreements
    winget install --id Go.Golang -e --accept-package-agreements --accept-source-agreements
}
elseif ($hasChoco) {
    Write-Info "Using Chocolatey to install Node.js (LTS) and Go."
    choco install nodejs-lts -y
    choco install golang -y
}
else {
    Write-Warn "Neither winget nor choco found."
    Write-Host "Options:"
    Write-Host "  1) Install winget (Microsoft App Installer) or Chocolatey first, then re-run this script."
    Write-Host "  2) Manually download installers: https://golang.org/dl and https://nodejs.org/"
    exit 2
}

Write-Info "Waiting briefly to allow environment to settle..."
Start-Sleep -Seconds 3

Write-Info "Verifying installations (may require a new shell if PATH changed)..."

try { & go version } catch { Write-Warn "go not found in PATH. You may need to open a new terminal." }
try { & node --version } catch { Write-Warn "node not found in PATH. You may need to open a new terminal." }
try { & npm --version } catch { Write-Warn "npm not found in PATH. You may need to open a new terminal." }

Write-Info "Showing git remotes for current repository (if any):"
try { git remote -v } catch { Write-Warn "git not available or not a repo here." }

Write-Info "If versions show, re-run the checks from the PMX README or ask me to continue." 
Write-Host "Suggested next steps:"
Write-Host "  Open a NEW PowerShell window and run:"
Write-Host "    go version"
Write-Host "    node --version"
Write-Host "    npm --version"
Write-Host "    git remote -v"
Write-Host "  Then run the validation steps I will perform for PMX."
