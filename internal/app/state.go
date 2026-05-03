package app

import (
	"github.com/rigofekete/gator/internal/config"
	"github.com/rigofekete/gator/internal/database"
)

type State struct {
	DB  *database.Queries
	Cfg *config.Config
}
