package handlers

import "log"

func logInfo(msg string, args ...interface{}) {
	log.Printf("[INFO] "+msg, args...)
}

func logError(msg string, args ...interface{}) {
	log.Printf("[ERROR] "+msg, args...)
}
