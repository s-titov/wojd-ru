param (
    [string]$unrealLocresExe,
    [string]$unpackedDir,
    [string]$hashAlgorithm = "SHA256"
)

if (-not $unrealLocresExe) {
    Write-Error "unrealLocresExe arg not passed"
    exit 1
}
if (-not $unpackedDir) {
    Write-Error "unpackedDir arg not passed"
    exit 1
}

$hantLocresOrig = Join-Path $unpackedDir "\ZhuxianClient\Content\Localization\Game\zh-Hant\Game.locres"
$hansLocresOrig = Join-Path $unpackedDir "\ZhuxianClient\gamedata\client\ZCTranslateData\Game\zh-Hans\Game.locres"
$enLocresOrig = Join-Path $unpackedDir "\ZhuxianClient\gamedata\client\ZCTranslateData\Game\en\Game.locres"
$ruLocresOrig = Join-Path $unpackedDir "\ZhuxianClient\gamedata\client\ZCTranslateData\Game\ru\Game.locres"

$root = (Split-Path $PSScriptRoot -Parent)
$locresPartsDir = Join-Path $root "patch\tw\Locres\parts"
$hantLocresPatch = "$locresPartsDir\Hant.locres"
$hansLocresPatch = "$locresPartsDir\Hans.locres"
$enLocresPatch = "$locresPartsDir\En.locres"
$ruLocresPatch = "$locresPartsDir\Ru.locres"

function Get-FileHashValue {
    param (
        [string]$path,
        [string]$algorithm
    )
    return (Get-FileHash -Path $path -Algorithm $algorithm).Hash
}

# hashes
$hantLocresOrigHash = Get-FileHashValue -Path $hantLocresOrig -Algorithm $hashAlgorithm
$hantLocresPatchHash = Get-FileHashValue -Path $hantLocresPatch -Algorithm $hashAlgorithm

$hansLocresOrigHash = Get-FileHashValue -Path $hansLocresOrig -Algorithm $hashAlgorithm
$hansLocresPatchHash = Get-FileHashValue -Path $hansLocresPatch -Algorithm $hashAlgorithm

$enLocresOrigHash = Get-FileHashValue -Path $enLocresOrig -Algorithm $hashAlgorithm
$enLocresPatchHash = Get-FileHashValue -Path $enLocresPatch -Algorithm $hashAlgorithm

$ruLocresOrigHash = Get-FileHashValue -Path $ruLocresOrig -Algorithm $hashAlgorithm
$ruLocresPatchHash = Get-FileHashValue -Path $ruLocresPatch -Algorithm $hashAlgorithm

$isChanged = $false

if ($hantLocresOrigHash -eq $hantLocresPatchHash) {
    Write-Host "✅  Hant locres file is actual!"
} else {
    $isChanged = $true
    Write-Host "❌  Hant locres file has been changed!"
}

if ($hansLocresOrigHash -eq $hansLocresPatchHash) {
    Write-Host "✅  Hans locres file is actual!"
} else {
    $isChanged = $true
    Write-Host "❌  Hans locres file has been changed!"
}

if ($enLocresOrigHash -eq $enLocresPatchHash) {
    Write-Host "✅  En locres file is actual!"
} else {
    $isChanged = $true
    Write-Host "❌  En locres file has been changed!"
}

if ($ruLocresOrigHash -eq $ruLocresPatchHash) {
    Write-Host "✅  Ru locres file is actual!"
} else {
    $isChanged = $true
    Write-Host "❌  Ru locres file has been changed!"
}

if ($isChanged) {
    # copy files
    Write-Host "------"
    Write-Host "Copying Locres to patch dir..."
    Copy-Item -Path $hantLocresOrig -Destination $hantLocresPatch -Force
    Copy-Item -Path $hansLocresOrig -Destination $hansLocresPatch -Force
    Copy-Item -Path $enLocresOrig -Destination $enLocresPatch -Force
    Copy-Item -Path $ruLocresOrig -Destination $ruLocresPatch -Force

    # rebuild locres
    Write-Host "Merging Locres to OriginalGame.locres..."
    $locresOriginal = Join-Path $root "patch\tw\Locres\OriginalGame.locres"
    & $unrealLocresExe merge "$locresOriginal" "$hantLocresPatch" -o $locresOriginal
    & $unrealLocresExe merge "$locresOriginal" "$hansLocresPatch" -o $locresOriginal
    & $unrealLocresExe merge "$locresOriginal" "$enLocresPatch" -o $locresOriginal
    & $unrealLocresExe merge "$locresOriginal" "$ruLocresPatch" -o $locresOriginal

    Write-Host "Locres successfully updated"
}
