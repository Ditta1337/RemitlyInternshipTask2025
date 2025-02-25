-- +goose Up
-- +goose StatementBegin
CREATE TABLE banks
(
    swiftCode            varchar(11) PRIMARY KEY,
    address              varchar(255),
    bankName             varchar(255) NOT NULL,
    countryISO2          varchar(2)   NOT NULL,
    countryName          varchar(255) NOT NULL,
    isHeadquarter        boolean      NOT NULL,
    headquarterSwiftCode varchar(11),
    CONSTRAINT fk_headquarter FOREIGN KEY (headquarterSwiftCode) REFERENCES banks(swiftCode) ON DELETE SET NULL
);

CREATE INDEX idx_headquarter ON banks(headquarterSwiftCode);
CREATE INDEX idx_country ON banks(countryISO2);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE banks
    DROP CONSTRAINT fk_headquarter;

DROP TABLE IF EXISTS banks;
-- +goose StatementEnd
