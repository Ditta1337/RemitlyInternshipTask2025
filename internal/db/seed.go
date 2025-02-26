package db

import (
	"context"
	"encoding/csv"
	"github.com/Ditta1337/RemitlyInternshipTask2025/internal/model"
	"github.com/Ditta1337/RemitlyInternshipTask2025/internal/store"
	"os"
	"path/filepath"
	"strings"
)

type BankRecord struct {
	CountryISO2 string
	SWIFTCode   string
	CodeType    string
	Name        string
	Address     *string
	TownName    string
	CountryName string
	TimeZone    string
}

func Seed(store store.Storage) error {
	ctx := context.Background()

	seedFilePath, err := filepath.Abs("internal/db/seed/SWIFT_CODES.tsv")
	if err != nil {
		return err
	}

	f, err := os.Open(seedFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.Comma = '\t'

	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// skip the header row
	records = records[1:]

	headquartersMap := make(map[string]string)

	// first pass, fill headquarters
	for _, record := range records {
		parsedRecord := parseRecord(record)

		if parsedRecord.SWIFTCode[8:] == "XXX" {

			var bank model.Bank

			bank.IsHeadquarter = true
			bank.SWIFTCode = parsedRecord.SWIFTCode
			bank.CountryISO2 = parsedRecord.CountryISO2
			bank.BankName = parsedRecord.Name
			bank.Address = parsedRecord.Address
			bank.CountryName = parsedRecord.CountryName
			bank.HeadquarterSWIFTCode = nil

			headquartersMap[parsedRecord.SWIFTCode[:8]] = parsedRecord.SWIFTCode

			if err := store.Banks.Create(ctx, &bank); err != nil {
				return err
			}
		}
	}

	// second pass, fill branches
	for _, record := range records {
		parsedRecord := parseRecord(record)

		if parsedRecord.SWIFTCode[8:] != "XXX" {

			var bank model.Bank

			bank.IsHeadquarter = false
			bank.SWIFTCode = parsedRecord.SWIFTCode
			bank.CountryISO2 = parsedRecord.CountryISO2
			bank.BankName = parsedRecord.Name
			bank.Address = parsedRecord.Address
			bank.CountryName = parsedRecord.CountryName

			if hqCode, ok := headquartersMap[parsedRecord.SWIFTCode[:8]]; ok {
				bank.HeadquarterSWIFTCode = &hqCode
			} else {
				bank.HeadquarterSWIFTCode = nil
			}

			if err := store.Banks.Create(ctx, &bank); err != nil {
				return err
			}
		}
	}

	return nil
}

func parseRecord(record []string) BankRecord {
	var address *string
	if record[4] == "" || len(strings.TrimSpace(record[4])) == 0 {
		address = nil
	} else {
		address = &record[4]
	}

	return BankRecord{
		CountryISO2: record[0],
		SWIFTCode:   record[1],
		CodeType:    record[2],
		Name:        record[3],
		Address:     address,
		TownName:    record[5],
		CountryName: record[6],
		TimeZone:    record[7],
	}
}
