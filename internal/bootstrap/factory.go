package bootstrap

import (
	customersRepo "ai_assistant/internal/customers/repository"
	wazzupRepo "ai_assistant/internal/wazzup/repository"

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
