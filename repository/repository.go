package repository

import (
	mdgeotrack "md-geo-track/iface"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) mdgeotrack.Repository {
	return &repository{
		db: db,
	}
}
