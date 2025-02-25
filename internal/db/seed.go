package db

import (
	//"context"
	"encoding/csv"
	"github.com/Ditta1337/RemitlyInternshipTask2025/internal/store"
	"log"
	"os"
	"path/filepath"
)

func Seed(store store.Storage) error {
	//ctx := context.Background()

	seedFilePath, err := filepath.Abs("internal/db/seed/SWIFT_CODES.csv")
	if err != nil {
		return err
	}

	f, err := os.Open(seedFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	reader := csv.NewReader(f)

	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// skip the header row
	records = records[1:]

	for _, record := range records {
		log.Printf("srocessing record: %s", record)
		// todo: parse rows
	}

	return nil
}
