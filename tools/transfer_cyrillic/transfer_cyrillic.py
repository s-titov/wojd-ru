import fontforge
import psMat
import sys
import math

# Входные аргументы
donor_path = sys.argv[1]
target_path = sys.argv[2]
output_path = sys.argv[3]

# Загружаем шрифты
donor = fontforge.open(donor_path)
target = fontforge.open(target_path)

# Диапазон Unicode кириллицы
cyr_start = 0x0400
cyr_end = 0x052F

for codepoint in range(cyr_start, cyr_end + 1):
    if codepoint not in donor:
        continue

    donor_glyph = donor[codepoint]
    glyph_name = donor_glyph.glyphname

    # Удалим старый глиф в целевом шрифте
    if codepoint in target:
        target.removeGlyph(target[codepoint].glyphname)

    # Создаем новый глиф
    target.createChar(codepoint, glyph_name)
    target_glyph = target[codepoint]

    # Копируем контур
    target_glyph.foreground = donor_glyph.foreground

    # Выравнивание по baseline оригинала
    _, donor_ymin, _, _ = donor_glyph.boundingBox()
    _, target_ymin, _, _ = target_glyph.boundingBox()
    dy = target_ymin - donor_ymin
    target_glyph.transform(psMat.translate(0, -dy))

    # Масштаб
    scale_factor = 0.28
    target_glyph.transform(psMat.scale(scale_factor))

    # AUTO WIDTH вручную: вычисляем bbox и центрируем
    target_glyph.width = int(donor_glyph.width * scale_factor)



# Сохраняем итог
target.generate(output_path)
print("✔ Готово:", output_path)
