param (
    [string]$quickbmsExe,
    [string]$gamePath,
    [string]$outputDir
)

if (-not $quickbmsExe) {
    Write-Error "Не передан аргумент quickbmsExe"
    exit 1
}
if (-not $gamePath) {
    Write-Error "Не передан аргумент gamePath"
    exit 1
}
if (-not $outputDir) {
    Write-Error "Не передан аргумент outputDir"
    exit 1
}

$paksDir = "$gamePath\ZXSJ\Game\ZhuxianClient\Content\Paks"
$bmsScript = "$PSScriptRoot/third_party/unreal_tournament_4_0.4.27e_zhuxian_world_ue5.bms"

# Обеспечим вывод в UTF-8
chcp 65001 > $null
[Console]::OutputEncoding = [System.Text.Encoding]::UTF8

# Создаем выходную директорию, если не существует
if (-not (Test-Path $outputDir)) {
    New-Item -ItemType Directory -Path $outputDir | Out-Null
}

# Обработка всех .pak файлов, кроме содержащих Ru_Patch
$paks = Get-ChildItem -Path $paksDir -Filter *.pak | Where-Object {
    $_.Name -notmatch "Ru_Patch"
}

$total = $paks.Count
$count = 0

foreach ($pak in $paks) {
    $count++
    $pakFile = $pak.FullName

    # Обновление прогресс-бара
    Write-Progress -Activity "Обработка .pak файлов" `
                   -Status "Файл $count из ${total}: $($pak.Name)" `
                   -PercentComplete (($count / $total) * 100)

    # Запуск команды без вывода
    & $quickbmsExe -o -Y $bmsScript $pakFile $outputDir
}
