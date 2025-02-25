package requests

type BankPayload struct {
	SWIFTCode     string `json:"swiftCode" validate:"required,len=11,alphanum"`
	Address       string `json:"address" validate:"max=255"`
	BankName      string `json:"bankName" validate:"required,max=255"`
	CountryISO2   string `json:"countryISO2" validate:"required,len=2,iso3166_1_alpha2"`
	CountryName   string `json:"countryName" validate:"required,max=255"`
	IsHeadquarter bool   `json:"isHeadquarter" validate:"required"`
}
