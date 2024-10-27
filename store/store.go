package store

import (
	"time"

	"github.com/Open-Code-Zone/cms/config"
	"github.com/Open-Code-Zone/cms/internal/database"
)

type Storage struct {
	DB          *database.Queries
	Collections *config.CollectionConfig
}

func NewStore(q *database.Queries, collections *config.CollectionConfig) *Storage {
	return &Storage{
		DB:          q,
		Collections: collections,
	}
}


