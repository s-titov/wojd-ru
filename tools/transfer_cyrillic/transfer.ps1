[Console]::OutputEncoding = [System.Text.Encoding]::UTF8

# Параметры
$fontforge = "D:\Programs\FontForgeBuilds\bin\fontforge.exe"
$script = "transfer_cyrillic.py"
$donor = "D:\JD_Russian\wojd-ru\tools\transfer_cyrillic\donor\NotoSansDisplay-Regular.ttf"
$target = "D:\JD_Russian\wojd-ru\instructions\fonts\FZShengSKSJW_Zhong.ttf"
$output = "target_cyrillic.ttf"

# Запуск FontForge
& "$fontforge" -lang=py -script $script $donor $target $output

$finalPath = "D:\JD_Russian\wojd-ru\patch\Ru_Patch_Strings_Main_P\ZhuxianClient\Content\UI\UI_Texture\UI_ziti\FZShengSKSJW_Zhong.ufont"

# Проверка успешности
if (Test-Path $output) {
    # Перенос и переименование
    Copy-Item -Path $output -Destination $finalPath -Force
    Write-Host "✔ Шрифт успешно создан и перенесён в:"
    Write-Host $finalPath
} else {
    Write-Error "❌ Не удалось создать файл $output"
}
