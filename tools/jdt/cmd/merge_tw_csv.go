package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: program <source_csv> <target_csv>")
		os.Exit(1)
	}

	srcCsv := os.Args[1]
	targetCsv := os.Args[2]

	data, err := collectData(srcCsv)
	if err != nil {
		panic(err)
	}

	addedLines, err := mergeData(targetCsv, data)
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully merged, added lines:", addedLines)
}

func collectData(srcCsv string) (map[string]string, error) {
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

func mergeData(targetCsv string, data map[string]string) (int, error) {
	origFile, err := os.Open(targetCsv)
	if err != nil {
		return 0, err
	}
	defer func() { _ = origFile.Close() }()

	reader := csv.NewReader(origFile)
	reader.FieldsPerRecord = -1

	var record []string
	for {
		record, err = reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil || len(record) < 3 {
			continue
		}

		// if already existed
		if _, ok := data[record[0]]; ok {
			delete(data, record[0])
		}
	}

	origFileWrite, err := os.OpenFile(targetCsv, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer func() { _ = origFileWrite.Close() }()

	writer := csv.NewWriter(origFileWrite)
	for key, source := range data {
		rec := []string{key, source}
		err = writer.Write(rec)
		if err != nil {
			return 0, err
		}
	}
	writer.Flush()

	return len(data), nil
}
