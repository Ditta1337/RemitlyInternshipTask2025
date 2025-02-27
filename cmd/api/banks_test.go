package main

import (
	"encoding/json"
	"github.com/Ditta1337/RemitlyInternshipTask2025/internal/dto/responses"
	"io"
	"log"
	"net/http"
	"strings"
	"testing"
)

func TestCreateBankHandler(t *testing.T) {
	app := newMockApplication(t)
	mux := app.mount()

	t.Run("should create headquarter bank", func(t *testing.T) {
		payload := `{
			"swiftCode": "FAKECODEXXX",
			"address": "Test Addr",
			"bankName": "Headquarter bank US",
			"countryISO2": "US",
			"countryName": "United States",
			"isHeadquarter": true
		}`
		req, err := http.NewRequest(http.MethodPost, app.config.apiVersion+"/swift-codes", strings.NewReader(payload))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rec := executeRequest(req, mux)

		body, err := io.ReadAll(rec.Body)
		if err != nil {
			t.Fatal(err)
		}

		var res responses.Message
		if err := json.Unmarshal(body, &res); err != nil {
			t.Fatalf("cannot unmarshal response to expected response.Message: %v", err)
		}

		expectedMessage := "successfully added bank to database"

		if res.Message != expectedMessage {
			t.Errorf("expected message %s, got %s", expectedMessage, res.Message)
		}

		checkResponseCode(t, http.StatusCreated, rec.Code)
	})

	t.Run("should create branch bank", func(t *testing.T) {
		payload := `{
			"swiftCode": "FAKECODE123",
			"address": "Test Addr",
			"bankName": "Branch bank US",
			"countryISO2": "US",
			"countryName": "United States",
			"isHeadquarter": false
		}`
		req, err := http.NewRequest(http.MethodPost, app.config.apiVersion+"/swift-codes", strings.NewReader(payload))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rec := executeRequest(req, mux)

		body, err := io.ReadAll(rec.Body)
		if err != nil {
			t.Fatal(err)
		}

		// Log the response body
		log.Println("Response Body:", string(body))

		var res responses.Message
		if err := json.Unmarshal(body, &res); err != nil {
			t.Fatalf("cannot unmarshal response: %v", err)
		}
		expectedMessage := "successfully added bank to database"

		if res.Message != expectedMessage {
			t.Errorf("expected message %q, got %q", expectedMessage, res.Message)
		}

		checkResponseCode(t, http.StatusCreated, rec.Code)
	})

	t.Run("invalid JSON payload", func(t *testing.T) {
		payload := `invalid json`
		req, err := http.NewRequest(http.MethodPost, app.config.apiVersion+"/swift-codes", strings.NewReader(payload))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rec := executeRequest(req, mux)
		checkResponseCode(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("swiftCode already exists in db", func(t *testing.T) {
		// bank with this swiftCode already exists in db
		payload := `{
			"swiftCode": "ABCDEFGHXXX",
			"address": "this bank already exists",
			"bankName": "Headquarter bank US",
			"countryISO2": "US",
			"countryName": "United States",
			"isHeadquarter": true
		}`

		req, err := http.NewRequest(http.MethodPost, app.config.apiVersion+"/swift-codes", strings.NewReader(payload))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rec := executeRequest(req, mux)
		checkResponseCode(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("validation error missing required field", func(t *testing.T) {
		// missing required bankName field
		payload := `{
			"swiftCode": "ABCDEFGHXXX",
			"address": "Test Addr",
			"countryISO2": "US",
			"countryName": "United States",
			"isHeadquarter": true
		}`
		req, err := http.NewRequest(http.MethodPost, app.config.apiVersion+"/swift-codes", strings.NewReader(payload))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rec := executeRequest(req, mux)
		checkResponseCode(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("validation error invalid swift code length", func(t *testing.T) {
		// swiftCode must be 11 alphanumeric characters
		payload := `{
			"swiftCode": "TOOSHORT",
			"address": "Test Addr",
			"bankName": "Invalid Swift Bank",
			"countryISO2": "US",
			"countryName": "United States",
			"isHeadquarter": false
		}`
		req, err := http.NewRequest(http.MethodPost, app.config.apiVersion+"/swift-codes", strings.NewReader(payload))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rec := executeRequest(req, mux)
		checkResponseCode(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("isHeadquarter mismatch error", func(t *testing.T) {
		// headquarter bank swift code must end with "XXX"
		payload := `{
			"swiftCode": "ABCDEFGH123",
			"address": "HQ Address",
			"bankName": "Mismatch HQ Bank",
			"countryISO2": "US",
			"countryName": "United States",
			"isHeadquarter": true
		}`
		req, err := http.NewRequest(http.MethodPost, app.config.apiVersion+"/swift-codes", strings.NewReader(payload))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rec := executeRequest(req, mux)
		checkResponseCode(t, http.StatusBadRequest, rec.Code)
	})
}

func TestGetBankBySWIFTCode(t *testing.T) {
	app := newMockApplication(t)
	mux := app.mount()

	t.Run("should return correct headquarter bank", func(t *testing.T) {
		swiftCode := "ABCDEFGHXXX"
		req, err := http.NewRequest(http.MethodGet, app.config.apiVersion+"/swift-codes/"+swiftCode, nil)
		if err != nil {
			t.Fatal(err)
		}

		rec := executeRequest(req, mux)

		body, err := io.ReadAll(rec.Body)
		if err != nil {
			t.Fatal(err)
		}

		var response responses.BankHeadquarter
		if err := json.Unmarshal(body, &response); err != nil {
			t.Fatalf("cannot unmarshal response to expected response.BankHeadquarter: %v", err)
		}

		expectedBankName := "Headquarter bank PL"
		expectedBankBranchesSize := 1

		if response.BankName != expectedBankName {
			t.Errorf("expected bank name: %s, got %s", expectedBankName, response.BankName)
		}

		if len(response.Branches) != expectedBankBranchesSize {
			t.Errorf("expected bank to have %d branch bank, got %d", expectedBankBranchesSize, len(response.Branches))
		}

		checkResponseCode(t, http.StatusOK, rec.Code)
	})

	t.Run("should return correct branch bank", func(t *testing.T) {
		swiftCode := "ABCDEFGH123"
		req, err := http.NewRequest(http.MethodGet, app.config.apiVersion+"/swift-codes/"+swiftCode, nil)
		if err != nil {
			t.Fatal(err)
		}

		rec := executeRequest(req, mux)

		body, err := io.ReadAll(rec.Body)
		if err != nil {
			t.Fatal(err)
		}

		var response responses.BankBranch
		if err := json.Unmarshal(body, &response); err != nil {
			t.Fatalf("cannot unmarshal response to expected response.BankBranch: %v", err)
		}

		expectedBankName := "Branch bank PL"

		if response.BankName != expectedBankName {
			t.Errorf("expected bank name: %s, got %s", expectedBankName, response.BankName)
		}

		checkResponseCode(t, http.StatusOK, rec.Code)
	})

	t.Run("nonexistent swiftCode", func(t *testing.T) {
		invalidSwiftCode := "INVALIDXXXX"
		req, err := http.NewRequest(http.MethodGet, app.config.apiVersion+"/swift-codes/"+invalidSwiftCode, nil)
		if err != nil {
			t.Fatal(err)
		}

		rec := executeRequest(req, mux)

		checkResponseCode(t, http.StatusNotFound, rec.Code)
	})

	t.Run("invalid swiftCode format", func(t *testing.T) {
		invalidSwiftCode := "TOOSHORT"
		req, err := http.NewRequest(http.MethodGet, app.config.apiVersion+"/swift-codes/"+invalidSwiftCode, nil)
		if err != nil {
			t.Fatal(err)
		}

		rec := executeRequest(req, mux)

		checkResponseCode(t, http.StatusBadRequest, rec.Code)
	})
}

func TestGetAllBanksByCountryISO2(t *testing.T) {
	app := newMockApplication(t)
	mux := app.mount()

	t.Run("should return banks with valid country ISO2", func(t *testing.T) {
		countryISO := "PL"
		req, err := http.NewRequest(http.MethodGet, app.config.apiVersion+"/swift-codes/country/"+countryISO, nil)
		if err != nil {
			t.Fatal(err)
		}

		rec := executeRequest(req, mux)

		body, err := io.ReadAll(rec.Body)
		if err != nil {
			t.Fatal(err)
		}

		var response responses.AllBanks
		if err := json.Unmarshal(body, &response); err != nil {
			t.Fatalf("cannot unmarshal response to expected response.AllBanks: %v", err)
		}

		expectedCountryISO := "PL"
		expectedCountryName := "Poland"
		expectedSwiftCodesSize := 2

		if response.CountryISO2 != expectedCountryISO {
			t.Errorf("expected country ISO code: %s, got %s", expectedCountryISO, response.CountryISO2)
		}

		if response.CountryName != expectedCountryName {
			t.Errorf("expected country name: %s, got %s", expectedCountryISO, response.CountryISO2)
		}

		if len(response.SwiftCodes) != expectedSwiftCodesSize {
			t.Errorf("expected %d swift codes, got %d", expectedSwiftCodesSize, len(response.SwiftCodes))
		}

		checkResponseCode(t, http.StatusOK, rec.Code)
	})

	t.Run("invalid country ISO2", func(t *testing.T) {
		countryISO := "TOOLONG"
		req, err := http.NewRequest(http.MethodGet, app.config.apiVersion+"/swift-codes/country/"+countryISO, nil)
		if err != nil {
			t.Fatal(err)
		}

		rec := executeRequest(req, mux)

		checkResponseCode(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("country ISO2 not in database", func(t *testing.T) {
		countryISO := "DE"
		req, err := http.NewRequest(http.MethodGet, app.config.apiVersion+"/swift-codes/country/"+countryISO, nil)
		if err != nil {
			t.Fatal(err)
		}

		rec := executeRequest(req, mux)

		checkResponseCode(t, http.StatusNotFound, rec.Code)
	})
}

func TestDeleteBankBySWIFTCode(t *testing.T) {
	app := newMockApplication(t)
	mux := app.mount()

	t.Run("should delete bank", func(t *testing.T) {
		swiftCode := "ABCDEFGHXXX"
		req, err := http.NewRequest(http.MethodDelete, app.config.apiVersion+"/swift-codes/"+swiftCode, nil)
		if err != nil {
			t.Fatal(err)
		}

		rec := executeRequest(req, mux)

		body, err := io.ReadAll(rec.Body)
		if err != nil {
			t.Fatal(err)
		}

		var res responses.Message
		if err := json.Unmarshal(body, &res); err != nil {
			t.Fatalf("cannot unmarshal response to expected response.Message: %v", err)
		}

		expectedMessage := "successfully deleted bank from database"

		if res.Message != expectedMessage {
			t.Errorf("expected message %s, got %s", expectedMessage, res.Message)
		}

		checkResponseCode(t, http.StatusOK, rec.Code)
	})

	t.Run("nonexistent swiftCode", func(t *testing.T) {
		invalidSwiftCode := "INVALIDXXXX"
		req, err := http.NewRequest(http.MethodDelete, app.config.apiVersion+"/swift-codes/"+invalidSwiftCode, nil)
		if err != nil {
			t.Fatal(err)
		}

		rec := executeRequest(req, mux)

		checkResponseCode(t, http.StatusNotFound, rec.Code)
	})

	t.Run("invalid swiftCode format", func(t *testing.T) {
		invalidSwiftCode := "TOOSHORT"
		req, err := http.NewRequest(http.MethodDelete, app.config.apiVersion+"/swift-codes/"+invalidSwiftCode, nil)
		if err != nil {
			t.Fatal(err)
		}

		rec := executeRequest(req, mux)

		checkResponseCode(t, http.StatusBadRequest, rec.Code)
	})
}
