$lines = git show "aca4d4e:dashboard/app/page.js" 2>$null
if ($lines) {
    for ($i = 175; $i -lt [Math]::Min(190, $lines.Count); $i++) {
        Write-Host ("{0,4}: {1}" -f ($i+1), $lines[$i])
    }
}
