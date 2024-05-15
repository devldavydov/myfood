package storage

import (
	"database/sql"
	"time"
)

const StorageOperationTimeout = 15 * time.Second

func isDriverRegistered(drvName string) bool {
	for _, name := range sql.Drivers() {
		if name == drvName {
			return true
		}
	}
	return false
}
