package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Ditta1337/RemitlyInternshipTask2025/internal/model"
	"time"
)

var (
	ErrNotFound          = errors.New("resource not found")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Banks interface {
		Create(context.Context, *model.Bank) error
		GetBySWIFTCode(context.Context, string) ([]model.Bank, error)
		GetAllByCountryISO2(context.Context, string) ([]model.Bank, error)
		Delete(context.Context, string) error
	}
}

func NewPostgresStorage(db *sql.DB) Storage {
	return Storage{
		Banks: &BankStore{db},
	}
}
