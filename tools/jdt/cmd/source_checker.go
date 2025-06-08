package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

func main() {
	sourceCsv := "original/tw/ZhuxianClient/Content/Localization/Game/zh-Hant/Game.csv"
	targetCsv := "patch/tw/Locres/Game.csv"

	keyMapSource, err := prepareKeySourceMap(sourceCsv)
	if err != nil {
		panic(err)
	}

	err = checkSource(targetCsv, keyMapSource)
	if err != nil {
		panic(err)
	}
}

func prepareKeySourceMap(srcCsv string) (map[string]string, error) {
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

func checkSource(targetCsv string, keySourceMap map[string]string) error {
	file, err := os.Open(targetCsv)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	// skip first row
	_, err = reader.Read()
	if err != nil {
		return err
	}

	var record []string
	for {
		record, err = reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if len(record) < 2 {
			continue
		}
		key := record[0]
		source := record[1]

		if origSource, ok := keySourceMap[key]; ok {
			if origSource != source {
				fmt.Println(fmt.Sprintf("Sources do not match: %s", key))
			}
		}
	}

	return nil
}
