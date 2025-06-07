package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/schollz/progressbar/v3"

	"translate_ai/adapter"
)

var alreadyTranslated = map[string]string{}
var bar = progressbar.NewOptions(-1,
	progressbar.OptionSetDescription("Translating rows"),
	progressbar.OptionShowCount(),
	progressbar.OptionSetElapsedTime(true),
	progressbar.OptionSetPredictTime(false),
	progressbar.OptionThrottle(100*time.Millisecond),
)

func main() {
	ctx := context.Background()

	path := "patch/tw/~Ru_Patch_P/ZhuxianClient/gamedata/client/FormatString/LoadingTips"
	info, err := os.Stat(path)
	if err != nil {
		panic(err)
	}

	if !info.IsDir() {
		err = translateTxtByAI(ctx, path)
		if err != nil {
			panic(err)
		}
	} else {
		err = translateDirByAI(ctx, path)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Successfully translated")
}

func translateTxtByAI(ctx context.Context, filePath string) error {
	chatGPT := adapter.NewAdapter()

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
					translate, err = chatGPT.Translate(ctx, value)
					if err != nil {
						return err
					}
					alreadyTranslated[value] = translate
					value = translate
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

func translateDirByAI(ctx context.Context, dirPath string) error {
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if filepath.Ext(path) == ".txt" {
				err = translateTxtByAI(ctx, path)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})

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
