param (
    [string]$unrealLocresExe,
    [string]$repakExe,
    [string]$gamePath,
    [string]$version
)

[Console]::OutputEncoding = [System.Text.Encoding]::UTF8

if (-not $unrealLocresExe) {
    Write-Error "Не передан аргумент unrealLocresExe"
    exit 1
}
if (-not $repakExe) {
    Write-Error "Не передан аргумент repakExe"
    exit 1
}
if (-not $gamePath) {
    Write-Error "Не передан аргумент gamePath"
    exit 1
}
if (-not $version) {
    Write-Error "Не передан аргумент version"
    exit 1
}
if ($version -notin @("cn", "tw")) {
    Write-Error "Доступные version: cn, tw"
    exit 1
}

# Set paths
$root = (Split-Path $PSScriptRoot -Parent)
$gameCsv = Join-Path $root "patch\$version\Locres\Game.csv"
$locresOriginal = Join-Path $root "patch\$version\Locres\OriginalGame.locres"
$locresNew = "$locresOriginal.new"
$pakFolder = Join-Path $root "patch\$version\~Ru_Patch_P"
$pakOutput = "$pakFolder.pak"
$hashFile = "$PSScriptRoot\hashes\build_patch_hash_$version.txt"
$pakFinal = "$gamePath\ZXSJ\Game\ZhuxianClient\Content\Paks\~Ru_Patch_P.pak"

$locresDir = "zh-Hans"
if ($version -eq "tw") {
    $locresDir = "zh-Hant"
}
$locresTarget = Join-Path $root "patch\$version\~Ru_Patch_P\ZhuxianClient\Content\Localization\Game\$locresDir\Game.locres"

# Calculate current CSV hash
$currentHash = Get-FileHash -Algorithm SHA256 $gameCsv | Select-Object -ExpandProperty Hash

# Load previous hash if exists
$previousHash = ""
if (Test-Path $hashFile) {
    $previousHash = Get-Content $hashFile -Raw
}

# Step 1-2: Only if hash changed
if ($currentHash -ne $previousHash) {
    Write-Host "··· Game.csv изменился, перезаписываем locres..."

    & $unrealLocresExe import "$locresOriginal" "$gameCsv" -f csv

    if (Test-Path $locresNew) {
        Copy-Item -Path $locresNew -Destination $locresTarget -Force
    } else {
        Write-Error "··· Файл $locresNew не найден после импорта"
        exit 1
    }

    $currentHash | Out-File -Encoding ASCII -NoNewline $hashFile
} else {
    Write-Host "··· Game.csv не изменился, пропускаем импорт locres"
}

# Step 3: Упаковка
Write-Host "··· Упаковка .pak..."
& $repakExe pack "$pakFolder"

# Step 4: Копирование .pak
if (Test-Path $pakOutput) {
    Copy-Item -Path $pakOutput -Destination $pakFinal -Force
    Write-Host "··· Файл .pak успешно скопирован в папку игры"
} else {
    Write-Error "··· Файл $pakOutput не найден после упаковки"
    exit 1
}

Write-Host "··· Работа скрипта завершена"
