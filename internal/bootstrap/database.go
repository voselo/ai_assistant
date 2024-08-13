package bootstrap

import (
	"ai_assistant/config"
	closer "ai_assistant/pkg/util"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func InitDB(cfg *config.Config) *sqlx.DB {

	pgDNS := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)

	db, err := sqlx.Connect("pgx", pgDNS)
	if err != nil {
		panic(fmt.Sprintf("err init db: %v", err))
	}

	closer.Add(db.Close)

	return db.Unsafe()
}
