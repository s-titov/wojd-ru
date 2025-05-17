# tools

## build_patch.ps1

Собирает патч и кладет в папку с игрой. Аргументы:
- gamePath - путь до клиента игры
- unrealLocresExe - путь до UnrealLocres.exe ([скачать](https://github.com/akintos/UnrealLocres/releases))
- путь до repak.exe ([скачать](https://github.com/trumank/repak/releases))

Пример запуска:
```bash
.\build_patch.ps1 -gamePath "D:\Games\ZXSJclient" -unrealLocresExe "D:\Programs\UnrealLocres\UnrealLocres.exe" -repakExe "D:\Programs\repak\repak.exe"
```
