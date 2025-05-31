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

func main() {
	sourceDir := "patch\\cn\\~Ru_Patch_P\\ZhuxianClient\\gamedata\\client\\FormatString"
	//targetDir := "patch\\tw\\~Ru_Patch_P\\ZhuxianClient\\gamedata\\client\\FormatString"

	// TODO: prepare translations map
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
		panic(err)
	}

	fmt.Println(translations)
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
