package storage

import "database/sql"

func isDriverRegistered(drvName string) bool {
	for _, name := range sql.Drivers() {
		if name == drvName {
			return true
		}
	}
	return false
}
