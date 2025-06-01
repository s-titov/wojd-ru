package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"unicode"

	"translate_ai/adapter"
)

func main() {
	ctx := context.Background()

	fileToTranslate := "patch/tw/~Ru_Patch_P/ZhuxianClient/gamedata/client/FormatString/ai.txt"
	err := translateTxtByAI(ctx, fileToTranslate)
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully translated")
}

func translateTxtByAI(ctx context.Context, filePath string) error {
	chatGPT := adapter.NewAdapter()

	var tempFile *os.File

	// match to translate
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	tempPath := file.Name() + ".tmp"
	tempFile, err = os.Create(tempPath)
	if err != nil {
		return err
	}
	defer func() { _ = tempFile.Close() }()

	scanner := bufio.NewScanner(file)
	translationsLimit := 10
	translationsCount := 0
	for scanner.Scan() {
		if translationsCount >= translationsLimit {
			break
		}

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

		var value string
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])

			value = strings.TrimSpace(parts[1])
			if needTranslate(value) {
				// TODO: Translate map
				value, err = chatGPT.Translate(ctx, value)
				translationsCount++
			}
			if err != nil {
				return err
			}

			newLine := fmt.Sprintf("%s = %s\n", key, value)
			_, err = tempFile.WriteString(newLine)
			if err != nil {
				return err
			}
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
	//err = os.Rename(tempPath, file.Name())
	//if err != nil {
	//	return err
	//}

	return err
}

// needTranslate проверяет, состоит ли строка только из цифр, латиницы и знаков препинания
func needTranslate(s string) bool {
	for _, r := range s {
		switch {
		case unicode.IsDigit(r):
			continue
		case unicode.IsPunct(r):
			continue
		case unicode.IsSpace(r):
			continue
		case unicode.IsLetter(r) && r <= unicode.MaxASCII: // latin
			continue
		case r >= 0x0400 && r <= 0x04FF: // cyrillic
			continue
		default:
			return true
		}
	}
	return false
}
