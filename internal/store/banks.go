package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Ditta1337/RemitlyInternshipTask2025/internal/model"
)

type BankStore struct {
	db *sql.DB
}

func (s *BankStore) Create(ctx context.Context, bank *model.Bank) error {
	existsQuery := "SELECT EXISTS (SELECT 1 FROM banks WHERE swiftCode = $1)"

	var exists bool
	err := s.db.QueryRowContext(ctx, existsQuery, bank.SWIFTCode).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return ErrAlreadyExists
	}

	headquarterSwiftCode, err := s.findHeadquarterSwiftCode(ctx, bank.SWIFTCode)
	if err != nil {
		return err
	}

	insertQuery := `
		INSERT INTO banks (swiftCode, address, bankName, countryISO2, countryName, isHeadquarter, headquarterSwiftCode)
		values ($1, $2, $3, $4, $5, $6, $7)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err = s.db.ExecContext(
		ctx,
		insertQuery,
		bank.SWIFTCode,
		bank.Address,
		bank.BankName,
		bank.CountryISO2,
		bank.CountryName,
		bank.IsHeadquarter,
		headquarterSwiftCode,
	)

	return err
}

// checks if headquarter of bank already exists in db
func (s *BankStore) findHeadquarterSwiftCode(ctx context.Context, swiftCode string) (*string, error) {
	swiftCodeToFind := swiftCode[:8] + "XXX"

	query := `
		SELECT swiftCode
		FROM banks
		WHERE swiftCode = $1
	`

	var foundSwiftCode string
	err := s.db.QueryRowContext(ctx, query, swiftCodeToFind).Scan(&foundSwiftCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &foundSwiftCode, nil
}

func (s *BankStore) GetBySWIFTCode(ctx context.Context, swiftCode string) ([]model.Bank, error) {
	headquarterQuery := `
		SELECT swiftCode, address, bankName, countryISO2, countryName, isHeadquarter, headquarterSwiftCode 
		FROM banks
		WHERE swiftCode = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var banks []model.Bank

	var headquarter model.Bank
	err := s.db.QueryRowContext(ctx, headquarterQuery, swiftCode).Scan(
		&headquarter.SWIFTCode,
		&headquarter.Address,
		&headquarter.BankName,
		&headquarter.CountryISO2,
		&headquarter.CountryName,
		&headquarter.IsHeadquarter,
		&headquarter.HeadquarterSWIFTCode,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	banks = append(banks, headquarter)

	branchesQuery := `
		SELECT swiftCode, address, bankName, countryISO2, countryName, isHeadquarter, headquarterSwiftCode 
		FROM banks
		WHERE headquarterSwiftCode = $1
	`

	rows, err := s.db.QueryContext(ctx, branchesQuery, swiftCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var branch model.Bank
		err := rows.Scan(
			&branch.SWIFTCode,
			&branch.Address,
			&branch.BankName,
			&branch.CountryISO2,
			&branch.CountryName,
			&branch.IsHeadquarter,
			&branch.HeadquarterSWIFTCode,
		)
		if err != nil {
			return nil, err
		}
		banks = append(banks, branch)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return banks, nil
}

func (s *BankStore) GetAllByCountryISO2(ctx context.Context, countryISO2 string) ([]model.Bank, error) {
	query := `
		SELECT swiftCode, address, bankName, countryISO2, countryName, isHeadquarter, headquarterSwiftCode 
		FROM banks
		WHERE countryISO2 = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var banks []model.Bank

	rows, err := s.db.QueryContext(ctx, query, countryISO2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var bank model.Bank
		err := rows.Scan(
			&bank.SWIFTCode,
			&bank.Address,
			&bank.BankName,
			&bank.CountryISO2,
			&bank.CountryName,
			&bank.IsHeadquarter,
			&bank.HeadquarterSWIFTCode,
		)
		if err != nil {
			return nil, err
		}
		banks = append(banks, bank)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(banks) == 0 {
		return nil, ErrNotFound
	}

	return banks, nil
}

func (s *BankStore) Delete(ctx context.Context, swiftCode string) error {
	query := `
		DELETE FROM banks
		WHERE swiftCode = $1
	`

	res, err := s.db.ExecContext(ctx, query, swiftCode)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}
