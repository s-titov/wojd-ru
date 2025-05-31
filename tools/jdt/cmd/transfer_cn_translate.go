package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

// transfer CN translates to TW .locres

type Translate struct {
	Target  string
	Comment string
}

func main() {
	source := "patch/cn/Locres/Game.csv"
	target := "patch/tw/Locres/Game.csv"

	translations, err := prepareCnTranslations(source)
	if err != nil {
		panic(err)
	}

	notTranslatedInc, err := translateTwTarget(target, translations)
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully translated")
	fmt.Println(fmt.Sprintf("Not translated: %d", notTranslatedInc))
}

func prepareCnTranslations(srcCsv string) (map[string]Translate, error) {
	data := make(map[string]Translate)

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
		if len(record) < 3 {
			continue
		}
		key := record[0]
		target := record[2]
		comment := record[3]
		data[key] = Translate{
			Target:  target,
			Comment: comment,
		}
	}

	return data, nil
}

func translateTwTarget(targetCsv string, translations map[string]Translate) (int, error) {
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
				record[2] = val.Target
				record[3] = val.Comment
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
