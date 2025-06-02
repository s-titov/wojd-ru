package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/schollz/progressbar/v3"
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
	//chatGPT := adapter.NewAdapter()

	var tempFile *os.File

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

	bar := progressbar.NewOptions(-1,
		progressbar.OptionSetDescription("Translating rows"),
		progressbar.OptionShowCount(),
		progressbar.OptionSetElapsedTime(true),
		progressbar.OptionSetPredictTime(false),
		progressbar.OptionThrottle(100*time.Millisecond),
	)

	scanner := bufio.NewScanner(file)
	translationsLimit := 6000
	translationsCount := 0
	alreadyTranslated := map[string]string{}
	for scanner.Scan() {
		line := scanner.Text()

		// transfer as is
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			_, err = tempFile.WriteString(line)
			if err != nil {
				return err
			}

			_ = bar.Add(1)
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			var value, translate string
			var ok bool

			key := strings.TrimSpace(parts[0])
			value = strings.TrimSpace(parts[1])
			if needTranslate(value) {
				if translate, ok = alreadyTranslated[value]; ok {
					value = translate
				} else {
					if translationsCount >= translationsLimit {
						break
					}

					value = value + " //TODO: Translate"
					//translate, err = chatGPT.Translate(ctx, value)
					//if err != nil {
					//	return err
					//}
					//alreadyTranslated[value] = translate
					//value = translate
					translationsCount++
				}
			}

			newLine := fmt.Sprintf("%s = %s\n", key, value)
			_, err = tempFile.WriteString(newLine)
			if err != nil {
				return err
			}
		}

		_ = bar.Add(1)
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
	err = os.Rename(tempPath, file.Name())
	if err != nil {
		return err
	}

	return err
}

func needTranslate(s string) bool {
	for _, r := range s {
		switch {
		case r <= 127: // ASCII (latin, digit, etc)
			continue
		case r >= 0x0400 && r <= 0x04FF: // cyrillic
			continue
		case unicode.IsPunct(r):
			continue
		case unicode.IsSpace(r):
			continue
		default:
			return true
		}
	}
	return false
}
