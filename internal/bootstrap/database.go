package bootstrap

import (
	"fmt"
	"messages_handler/config"
	closer "messages_handler/pkg/util"

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
