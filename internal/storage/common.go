package storage

import "log"

func logInfo(msg string, args ...any) {
	log.Printf("[DBINFO] "+msg, args...)
}

func logError(msg string, args ...any) {
	log.Printf("[DBERROR] "+msg, args...)
}
