package responses

type AllBanks struct {
	CountryISO2 string      `json:"countryISO2"`
	CountryName string      `json:"countryName"`
	SwiftCodes  []BankShort `json:"swiftCodes"`
}
