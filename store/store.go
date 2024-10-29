package store

import (
	"database/sql"

	"github.com/Open-Code-Zone/cms/config"
	"github.com/Open-Code-Zone/cms/internal/database"
)

type Storage struct {
	DB          *sql.DB
	Queries     *database.Queries
	Collections *config.CollectionConfig
}

func NewStore(db *sql.DB, queries *database.Queries, collections *config.CollectionConfig) *Storage {
	return &Storage{
		DB:          db,
		Queries:     queries,
		Collections: collections,
	}
}
