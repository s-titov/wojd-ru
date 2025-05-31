package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type FileName string
type TranslationKey string
type Translations map[FileName]map[TranslationKey]string

func main() {
	sourceDir := "patch\\cn\\~Ru_Patch_P\\ZhuxianClient\\gamedata\\client\\FormatString"
	targetDir := "patch\\tw\\~Ru_Patch_P\\ZhuxianClient\\gamedata\\client\\FormatString"

	translations, err := prepareTxtTranslations(sourceDir)
	if err != nil {
		panic(err)
	}

	err = translateTxtTarget(targetDir, translations)
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully translated")
}

func prepareTxtTranslations(sourceDir string) (Translations, error) {
	translations := map[FileName]map[TranslationKey]string{}
	err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if filepath.Ext(path) == ".txt" {
				var translationMap map[TranslationKey]string
				translationMap, err = parseTranslationFile(path)
				if err != nil {
					return err
				}

				translations[FileName(strings.TrimPrefix(path, sourceDir))] = translationMap
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return translations, nil
}

func parseTranslationFile(filePath string) (map[TranslationKey]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	result := make(map[TranslationKey]string)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// skip empty lines
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := TranslationKey(strings.TrimSpace(parts[0]))
			value := strings.TrimSpace(parts[1])
			result[key] = value
		}
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func translateTxtTarget(targetDir string, translations Translations) error {
	var fileName FileName
	var file *os.File
	var tempFile *os.File

	err := filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".txt" {
			fileName = FileName(strings.TrimPrefix(path, targetDir))

			// if we have translations for this file
			if fileTranslations, ok := translations[fileName]; ok {
				// match to translate
				file, err = os.Open(path)
				if err != nil {
					return err
				}
				defer func() { _ = file.Close() }()

				targetPath := file.Name()
				tempPath := file.Name() + ".tmp"
				tempFile, err = os.Create(tempPath)
				if err != nil {
					return err
				}
				defer func() { _ = tempFile.Close() }()

				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					line := scanner.Text()

					// transfer as is
					line = strings.TrimSpace(line)
					if line == "" || strings.HasPrefix(line, "#") {
						_, err = tempFile.WriteString(line)
						if err != nil {
							return err
						}

						continue
					}

					// mappings translate
					parts := strings.SplitN(line, "=", 2)
					if len(parts) == 2 {
						key := TranslationKey(strings.TrimSpace(parts[0]))
						value := strings.TrimSpace(parts[1]) + " //TODO: Translate"
						if translation, ok := fileTranslations[key]; ok {
							value = translation
						}

						newLine := fmt.Sprintf("%s = %s\n", key, value)
						_, err = tempFile.WriteString(newLine)
					}
				}

				err = tempFile.Close()
				if err != nil {
					return err
				}
				err = file.Close()
				if err != nil {
					return err
				}

				// rewrite original file
				err = os.Rename(tempPath, targetPath)
				if err != nil {
					return err
				}
			}

		}
		return nil
	})

	return err
}
