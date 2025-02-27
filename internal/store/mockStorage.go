package store

import (
	"context"
	"github.com/Ditta1337/RemitlyInternshipTask2025/internal/model"
)

func NewMockStorage() Storage {
	headquarterSWIFTCode := "ABCDEFGHXXX"
	return Storage{
		Banks: &MockBankStore{
			banks: []model.Bank{
				{
					SWIFTCode:     headquarterSWIFTCode,
					BankName:      "Headquarter bank PL",
					CountryISO2:   "PL",
					CountryName:   "Poland",
					IsHeadquarter: true,
				},
				{
					SWIFTCode:            "ABCDEFGH123",
					BankName:             "Branch bank PL",
					CountryISO2:          "PL",
					CountryName:          "Poland",
					IsHeadquarter:        false,
					HeadquarterSWIFTCode: &headquarterSWIFTCode,
				},
			},
		},
	}
}

type MockBankStore struct {
	banks []model.Bank
}

func (m *MockBankStore) Create(ctx context.Context, bank *model.Bank) error {
	for _, existingBank := range m.banks {
		if existingBank.SWIFTCode == bank.SWIFTCode {
			return ErrAlreadyExists
		}
	}

	headquarterSwiftCode, err := m.findHeadquarterSwiftCode(ctx, bank.SWIFTCode)
	if err != nil {
		return err
	}
	bank.HeadquarterSWIFTCode = headquarterSwiftCode

	m.banks = append(m.banks, *bank)
	return nil
}

func (m *MockBankStore) findHeadquarterSwiftCode(ctx context.Context, swiftCode string) (*string, error) {
	swiftCodeToFind := swiftCode[:8] + "XXX"

	for _, bank := range m.banks {
		if bank.SWIFTCode == swiftCodeToFind {
			return &bank.SWIFTCode, nil
		}
	}

	return nil, nil
}
func (m *MockBankStore) GetBySWIFTCode(ctx context.Context, swiftCode string) ([]model.Bank, error) {
	var headquarter model.Bank
	var branches []model.Bank

	for _, bank := range m.banks {
		if bank.SWIFTCode == swiftCode {
			headquarter = bank
		} else if bank.HeadquarterSWIFTCode != nil && *bank.HeadquarterSWIFTCode == swiftCode {
			branches = append(branches, bank)
		}
	}

	if headquarter.SWIFTCode == "" {
		return nil, ErrNotFound
	}

	return append([]model.Bank{headquarter}, branches...), nil
}

func (m *MockBankStore) GetAllByCountryISO2(ctx context.Context, countryISO2 string) ([]model.Bank, error) {
	var banks []model.Bank

	for _, bank := range m.banks {
		if bank.CountryISO2 == countryISO2 {
			banks = append(banks, bank)
		}
	}

	if len(banks) == 0 {
		return nil, ErrNotFound
	}

	return banks, nil
}

func (m *MockBankStore) Delete(ctx context.Context, swiftCode string) error {
	for i, bank := range m.banks {
		if bank.SWIFTCode == swiftCode {
			m.banks = append(m.banks[:i], m.banks[i+1:]...)
			return nil
		}
	}
	return ErrNotFound
}
