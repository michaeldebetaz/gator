package state

import (
	"database/sql"
	"gator/internal/config"
	"gator/internal/database"
)

type State struct {
	Cfg     *config.Config
	Db      *sql.DB
	Queries *database.Queries
}
