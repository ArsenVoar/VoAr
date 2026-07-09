// Package logger provides application-specific logging helpers.
// It centralizes log formatting for requests, transactions and cache events.
package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func Init() error {
	file, err := os.OpenFile(
		"logs/app.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)
	if err != nil {
		return err
	}

	writer := io.MultiWriter(os.Stdout, file)

	log.SetOutput(writer)
	log.SetFlags(0)

	return nil
}

func Info(message string) {
	timestamp := time.Now().Format(time.RFC3339)

	log.Printf("INFO %s %s", timestamp, message)
}
func Error(message string) {
	timestamp := time.Now().Format(time.RFC3339)

	log.Printf("ERROR %s %s", timestamp, message)
}
func Request(requestID string, method string, path string, status int, duration time.Duration) {
	log.Printf(
		"REQUEST request_id=%s method=%s path=%s status=%d duration=%s", requestID, method, path, status, duration)
}

func TransactionStarted(requestID string) {
	message := fmt.Sprintf("request_id=%s tx_started", requestID)
	Info(message)
}
func TransactionCommitted(requestID string) {
	message := fmt.Sprintf("request_id=%s tx_committed", requestID)
	Info(message)
}
func TransactionRollback(requestID string, err error) {
	message := fmt.Sprintf("request_id=%s tx_rollback err=%s", requestID, err)
	Info(message)
}
func TransactionFailed(requestID string, err error) {
	message := fmt.Sprintf("request_id=%s tx_failed err=%s", requestID, err)
	Error(message)
}

func CacheHit(requestID string, key string) {
	message := fmt.Sprintf("request_id=%s cache_hit key=%s", requestID, key)
	Info(message)
}
func CacheDeleted(requestID string, key string) {
	message := fmt.Sprintf("request_id=%s cache_deleted key=%s", requestID, key)
	Info(message)
}
func CacheSet(requestID string, key string) {
	message := fmt.Sprintf("request_id=%s cache_set key=%s", requestID, key)
	Info(message)
}
func CacheMiss(requestID string, key string) {
	message := fmt.Sprintf("request_id=%s cache_miss key=%s", requestID, key)
	Info(message)
}
func CacheError(requestID string, key string, err error) {
	message := fmt.Sprintf("request_id=%s cache_error key=%s err=%s", requestID, key, err)
	Error(message)
}
