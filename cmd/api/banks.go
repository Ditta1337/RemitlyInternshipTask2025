package main

import (
	"errors"
	"github.com/Ditta1337/RemitlyInternshipTask2025/internal/dto/requests"
	"github.com/Ditta1337/RemitlyInternshipTask2025/internal/dto/responses"
	"github.com/Ditta1337/RemitlyInternshipTask2025/internal/model"
	"github.com/Ditta1337/RemitlyInternshipTask2025/internal/store"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
)

func (app *application) createBankHandler(w http.ResponseWriter, r *http.Request) {
	var payload requests.BankPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var address *string
	if payload.Address == "" || len(strings.TrimSpace(payload.Address)) == 0 {
		address = nil
	} else {
		address = &payload.Address
	}

	if payload.IsHeadquarter && payload.SWIFTCode[8:] != "XXX" {
		app.badRequestResponse(w, r, errors.New("isHeadquarters set to true, while swift code says otherwise"))
	}

	bank := &model.Bank{
		SWIFTCode:     payload.SWIFTCode,
		Address:       address,
		BankName:      payload.BankName,
		CountryISO2:   payload.CountryISO2,
		CountryName:   payload.CountryName,
		IsHeadquarter: payload.IsHeadquarter,
	}

	ctx := r.Context()

	if err := app.store.Banks.Create(ctx, bank); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.writeJSONResponse(w, http.StatusCreated, responses.Message{Message: "successfully added bank to database"}); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getBankBySWIFTCodeHandler(w http.ResponseWriter, r *http.Request) {
	swiftCode := strings.ToUpper(chi.URLParam(r, "swift-code"))

	ctx := r.Context()

	banks, err := app.store.Banks.GetBySWIFTCode(ctx, swiftCode)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	firstBank := banks[0]
	if firstBank.IsHeadquarter {
		bankHeadquarter := mapBankToBankHeadquarter(firstBank, banks[1:])

		if err := app.writeJSONResponse(w, http.StatusOK, bankHeadquarter); err != nil {
			app.internalServerError(w, r, err)
			return
		}
	} else {
		bankBranch := mapBankToBankBranch(firstBank)

		if err := app.writeJSONResponse(w, http.StatusOK, bankBranch); err != nil {
			app.internalServerError(w, r, err)
			return
		}
	}
}

func (app *application) getAllBanksByCountryISO2Handler(w http.ResponseWriter, r *http.Request) {
	countryISO := strings.ToUpper(chi.URLParam(r, "countryISO2code"))
	if len(countryISO) != 2 {
		app.badRequestResponse(w, r, errors.New("incorrect country iso2 code"))
	}

	ctx := r.Context()

	banks, err := app.store.Banks.GetAllByCountryISO2(ctx, countryISO)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):

			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	firstBank := banks[0]
	allBanks := responses.AllBanks{
		CountryISO2: firstBank.CountryISO2,
		CountryName: firstBank.CountryName,
		SwiftCodes:  mapBanksToBankShorts(banks),
	}

	if err := app.writeJSONResponse(w, http.StatusOK, allBanks); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) deleteBankHandler(w http.ResponseWriter, r *http.Request) {
	swiftCode := strings.ToUpper(chi.URLParam(r, "swift-code"))

	ctx := r.Context()

	if err := app.store.Banks.Delete(ctx, swiftCode); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
	}

	if err := app.writeJSONResponse(w, http.StatusCreated, responses.Message{Message: "successfully deleted bank from database"}); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func mapBankToBankHeadquarter(bank model.Bank, branches []model.Bank) responses.BankHeadquarter {
	return responses.BankHeadquarter{
		SWIFTCode:     bank.SWIFTCode,
		BankName:      bank.BankName,
		Address:       bank.Address,
		CountryISO2:   bank.CountryISO2,
		CountryName:   bank.CountryName,
		IsHeadquarter: bank.IsHeadquarter,
		Branches:      mapBanksToBankShorts(branches),
	}
}

func mapBankToBankBranch(bank model.Bank) responses.BankBranch {
	return responses.BankBranch{
		SWIFTCode:     bank.SWIFTCode,
		Address:       bank.Address,
		BankName:      bank.BankName,
		CountryISO2:   bank.CountryISO2,
		CountryName:   bank.CountryName,
		IsHeadquarter: bank.IsHeadquarter,
	}
}

func mapBanksToBankShorts(banks []model.Bank) []responses.BankShort {
	var shortBanks []responses.BankShort

	for _, bank := range banks {
		var shortBank responses.BankShort

		shortBank.SWIFTCode = bank.SWIFTCode
		shortBank.Address = bank.Address
		shortBank.CountryName = bank.CountryName
		shortBank.CountryISO2 = bank.CountryISO2
		shortBank.IsHeadquarter = bank.IsHeadquarter

		shortBanks = append(shortBanks, shortBank)
	}

	return shortBanks
}
