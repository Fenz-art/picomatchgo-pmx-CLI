$content = Get-Content 'c:\Users\KIIT\OneDrive\Desktop\pmx\dashboard\app\page.js' -Encoding UTF8
$braceDepth = 0
$parenDepth = 0
$bracketDepth = 0
$inString = $false
$stringChar = ''
$inTemplate = $false
$templateDepth = 0

for ($lineNum = 0; $lineNum -lt [Math]::Min(185, $content.Count); $lineNum++) {
    $line = $content[$lineNum]
    $prevBrace = $braceDepth
    $prevParen = $parenDepth
    $prevBracket = $bracketDepth
    
    for ($i = 0; $i -lt $line.Length; $i++) {
        $ch = $line[$i]
        
        # Skip string contents (simplified - doesn't handle escapes perfectly)
        if ($inString) {
            if ($ch -eq $stringChar -and ($i -eq 0 -or $line[$i-1] -ne '\')) {
                $inString = $false
            }
            continue
        }
        
        if ($ch -eq "'" -or $ch -eq '"') {
            $inString = $true
            $stringChar = $ch
            continue
        }
        
        # Skip template literal contents (simplified)
        if ($ch -eq '`') {
            $inTemplate = -not $inTemplate
            continue
        }
        
        if ($inTemplate) {
            continue
        }
        
        # Skip single-line comments
        if ($ch -eq '/' -and $i + 1 -lt $line.Length -and $line[$i+1] -eq '/') {
            break
        }
        
        switch ($ch) {
            '{' { $braceDepth++ }
            '}' { $braceDepth-- }
            '(' { $parenDepth++ }
            ')' { $parenDepth-- }
            '[' { $bracketDepth++ }
            ']' { $bracketDepth-- }
        }
    }
    
    if ($braceDepth -ne $prevBrace -or $parenDepth -ne $prevParen -or $bracketDepth -ne $prevBracket) {
        Write-Host ("L" + ($lineNum + 1) + ": {" + $braceDepth + "} (" + $parenDepth + ") [" + $bracketDepth + "]")
    }
    
    if ($braceDepth -lt 0 -or $parenDepth -lt 0 -or $bracketDepth -lt 0) {
        Write-Host ("*** NEGATIVE DEPTH at line " + ($lineNum + 1) + " ***")
    }
}

Write-Host ""
Write-Host ("Final at line 185: {" + $braceDepth + "} (" + $parenDepth + ") [" + $bracketDepth + "]")
