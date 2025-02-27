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

// CreateBank godoc
//
//	@Summary		Creates a bank
//	@Description	Creates a bank
//	@Tags			banks
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		requests.BankPayload	true	"Bank payload"
//	@Success		201		{object}	responses.Message
//	@Failure		400		{object}	responses.Error
//	@Failure		500		{object}	responses.Error
//	@Router			/swift-codes [post]
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

// GetBankBySWIFTCode godoc
//
//	@Summary		Gets a bank by SWIFT code
//	@Description	Gets a bank by SWIFT code
//	@Tags			banks
//	@Accept			json
//	@Produce		json
//	@Param			swift-code	path		string		true	"SWIFT Code"
//	@Success		200			{object}	interface{}	"Returns either a BankHeadquarter or BankBranch. See the API documentation for details."
//	@Failure		400			{object}	responses.Error
//	@Failure		404			{object}	responses.Error
//	@Failure		500			{object}	responses.Error
//	@Router			/swift-codes/{swift-code} [get]
func (app *application) getBankBySWIFTCodeHandler(w http.ResponseWriter, r *http.Request) {
	swiftCode, err := parseSwiftCode(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

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

// GetAllBanksByCountryISO2Code godoc
//
//	@Summary		Gets all banks with given Country ISO2 Code
//	@Description	Gets all banks with given Country ISO2 Code
//	@Tags			banks
//	@Accept			json
//	@Produce		json
//	@Param			countryISO2code	path		string	true	"Country ISO2 Code"
//	@Success		200				{object}	responses.AllBanks
//	@Failure		400				{object}	responses.Error
//	@Failure		500				{object}	responses.Error
//	@Router			/swift-codes/country/{countryISO2code} [get]
func (app *application) getAllBanksByCountryISO2Handler(w http.ResponseWriter, r *http.Request) {
	countryISO := strings.ToUpper(chi.URLParam(r, "countryISO2code"))
	if len(countryISO) != 2 {
		app.badRequestResponse(w, r, errors.New("incorrect country ISO2 code length"))
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

// DeleteBank godoc
//
//	@Summary		Deletes a bank by SWIFT code
//	@Description	Deletes a bank by SWIFT code
//	@Tags			banks
//	@Accept			json
//	@Produce		json
//	@Param			swift-code	path		string	true	"SWIFT Code"
//	@Success		200			{object}	responses.Message
//	@Failure		404			{object}	responses.Error
//	@Failure		500			{object}	responses.Error
//	@Router			/swift-codes/{swift-code} [delete]
func (app *application) deleteBankHandler(w http.ResponseWriter, r *http.Request) {
	swiftCode, err := parseSwiftCode(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	ctx := r.Context()

	if err := app.store.Banks.Delete(ctx, swiftCode); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
	}

	if err := app.writeJSONResponse(w, http.StatusOK, responses.Message{Message: "successfully deleted bank from database"}); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func parseSwiftCode(r *http.Request) (string, error) {
	swiftCode := strings.ToUpper(chi.URLParam(r, "swift-code"))
	if len(swiftCode) != 11 {
		return "", errors.New("incorrect SWIFT code length")
	}

	return swiftCode, nil
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
