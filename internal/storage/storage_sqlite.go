package storage

import "database/sql"

type StorageSQLite struct {
	db *sql.DB
}

func NewStorageSQLite() *StorageSQLite {
	return &StorageSQLite{}
}
