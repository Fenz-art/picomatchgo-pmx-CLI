$bytes = [System.IO.File]::ReadAllBytes('c:\Users\KIIT\OneDrive\Desktop\pmx\dashboard\app\page.js')
$content = [System.Text.Encoding]::UTF8.GetString($bytes)
$lines = $content.Split([char]10)
for ($i = 172; $i -lt 186; $i++) {
    $line = $lines[$i]
    $hex = ''
    foreach ($c in $line.ToCharArray()) {
        if ([int]$c -gt 127 -and [int]$c -ne 13) {
            $hex += ' [U+' + [string]::Format('{0:X4}', [int]$c) + ']'
        }
    }
    if ($hex) {
        Write-Host ("Line " + ($i+1) + ":" + $hex)
    }
}
Write-Host "Done"
