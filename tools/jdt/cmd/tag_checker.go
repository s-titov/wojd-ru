package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"regexp"
	"slices"
)

func main() {
	srcCsv := "patch/tw/Locres/Game.csv"
	err := checkLinkTag(srcCsv)
	if err != nil {
		panic(err)
	}
}

func checkLinkTag(srcCsv string) error {
	file, err := os.Open(srcCsv)
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
		if len(record) < 3 {
			continue
		}
		key := record[0]
		source := record[1]
		target := record[2]
		if target != "" {
			reLinkOpen := regexp.MustCompile(`\(@@Link:[0-9]*\)`)
			sourceMatches := reLinkOpen.FindAllString(source, -1)

			// если в source есть Link
			if len(sourceMatches) > 0 {
				targetMatches := reLinkOpen.FindAllString(target, -1)

				if len(sourceMatches) != len(targetMatches) {
					fmt.Println(fmt.Sprintf("qty not matched: %s", key))
					continue
				}

				// сверяем, что линки матчатся с таргетом
				for _, sourceMatch := range sourceMatches {
					if !slices.Contains(targetMatches, sourceMatch) {
						fmt.Println(fmt.Sprintf("Links not match: %s", key))
						continue
					}
				}

				reLinkClose := regexp.MustCompile(`\(/Link\)`)
				targetCloseMatches := reLinkClose.FindAllString(target, -1)
				if len(targetCloseMatches) != len(targetMatches) {
					fmt.Println(fmt.Sprintf("(/Link) count incorrect: %s", key))
				}
			}
		}
	}

	return nil
}
