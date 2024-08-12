package bootstrap

import (
	customersRepo "messages_handler/internal/customers/repository"
	wazzupRepo "messages_handler/internal/wazzup/repository"

	"github.com/jmoiron/sqlx"
)

type Factory struct {
	CustomersRepository *customersRepo.CustomersRepository
	WazzupRepository    *wazzupRepo.WazzupRepository
}

func NewFactory(db *sqlx.DB) *Factory {
	return &Factory{
		CustomersRepository: customersRepo.New(db),
		WazzupRepository:    wazzupRepo.New(),
	}
}
