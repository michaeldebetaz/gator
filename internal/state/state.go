package state

import (
	"gator/internal/config"
	"gator/internal/database"
)

type State struct {
	Cfg *config.Config
	Db  *database.Queries
}
