package repository

import (
	"github.com/jmoiron/sqlx"
)

type Factory struct {
	CustomersRepository *CustomersRepository
	WazzupRepository    *WazzupRepository
}

func NewFactory(db *sqlx.DB) *Factory {
	return &Factory{
		CustomersRepository: NewCustomerRepository(db),
		WazzupRepository:    NewWazzupRepository(),
	}
}
