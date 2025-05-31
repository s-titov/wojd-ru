# tools

## build_patch.ps1

Собирает патч и кладет в папку с игрой. Аргументы:
- versions - версия клиента ("tw" - тайвань или "cn" - китай)
- gamePath - путь до клиента игры
- unrealLocresExe - путь до UnrealLocres.exe ([скачать](https://github.com/akintos/UnrealLocres/releases))
- путь до repak.exe ([скачать](https://github.com/trumank/repak/releases))

Пример запуска:
```bash
.\build_patch.ps1 -version "tw" -gamePath "D:\Games\zxsjgt" -unrealLocresExe "D:\Programs\UnrealLocres\UnrealLocres.exe" -repakExe "D:\Programs\repak\repak.exe"
```

## unpack_game.ps1

Распаковывает все паки с игры (только .txt, .locres и .json файлы).
<br>При распаковке файлы весом больше 100мб пропускаются для ускорения распаковки.

Аргументы:
- quibckmsExe - путь до quickbms_4gb_files.exe ([скачать](https://github.com/LittleBigBug/QuickBMS/releases))
- gamePath - путь до игры
- outputDir - директория куда распаковывать файлы

Пример запуска:
```bash
.\unpack_game.ps1 -quickbmsExe "D:\Programs\QuickBMS\quickbms_4gb_files.exe" -gamePath "D:\Games\zxsjgt" -outputDir "D:\JD_Russian\JDUnpacked"
```

\* bms скрипт и ключ шифрования игры (AES key в bms скрипте) брались [отсюда](https://cs.rin.ru/forum/viewtopic.php?f=10&t=100672), если они больше не подходят можно покопаться в этом треде
