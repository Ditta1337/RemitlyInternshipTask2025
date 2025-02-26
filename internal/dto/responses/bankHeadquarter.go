package responses

type BankHeadquarter struct {
	SWIFTCode     string      `json:"swiftCode"`
	Address       *string     `json:"address"`
	BankName      string      `json:"bankName"`
	CountryISO2   string      `json:"countryISO2"`
	CountryName   string      `json:"countryName"`
	IsHeadquarter bool        `json:"isHeadquarter"`
	Branches      []BankShort `json:"branches"`
}
type BankShort struct {
	SWIFTCode     string  `json:"swiftCode"`
	Address       *string `json:"address"`
	CountryISO2   string  `json:"countryISO2"`
	CountryName   string  `json:"countryName"`
	IsHeadquarter bool    `json:"isHeadquarter"`
}
