package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

func main() {
	source := "patch/Locres/GameTwRu.csv"
	target := "patch/Locres/Game.csv"

	translations, err := prepareTranslations(source)
	if err != nil {
		panic(err)
	}

	notTranslatedInc, err := translateTarget(target, translations)
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully translated")
	fmt.Println(fmt.Sprintf("Not translated count: %d", notTranslatedInc))
}

func prepareTranslations(srcCsv string) (map[string]string, error) {
	data := make(map[string]string)

	file, err := os.Open(srcCsv)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	// skip first row
	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

	var record []string
	for {
		record, err = reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(record) < 2 {
			continue
		}
		key := record[0]
		source := record[1]
		data[key] = source
	}

	return data, nil
}

func translateTarget(targetCsv string, translations map[string]string) (int, error) {
	tempPath := targetCsv + ".tmp"

	origFile, err := os.Open(targetCsv)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = origFile.Close()
	}()

	reader := csv.NewReader(origFile)
	reader.FieldsPerRecord = -1

	tempFile, err := os.Create(tempPath)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = tempFile.Close()
		_ = os.Remove(tempPath)
	}()

	writer := csv.NewWriter(tempFile)

	notTranslatedInc := 0
	var record []string
	for {
		record, err = reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil || len(record) < 3 {
			continue
		}
		// if "target" column is empty
		if record[2] == "" {
			if val, ok := translations[record[0]]; ok {
				record[2] = val
				record[3] = "tw-translate"
			} else {
				notTranslatedInc++
			}
		}
		err = writer.Write(record)
		if err != nil {
			return 0, err
		}
	}
	writer.Flush()

	err = tempFile.Close()
	if err != nil {
		return 0, err
	}
	err = origFile.Close()
	if err != nil {
		return 0, err
	}

	// rewrite original file
	err = os.Rename(tempPath, targetCsv)
	if err != nil {
		return 0, err
	}

	return notTranslatedInc, nil
}
