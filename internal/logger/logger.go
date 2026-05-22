package logger

import (
	"fmt"
	"log"
	"time"
)

func Info(message string) {
	timestamp := time.Now().Format(time.RFC3339)

	log.Println(
		fmt.Sprintf(
			"INFO %s %s",
			timestamp,
			message,
		),
	)
}

func Error(message string) {
	timestamp := time.Now().Format(time.RFC3339)

	log.Println(
		fmt.Sprintf(
			"ERROR %s %s",
			timestamp,
			message,
		),
	)
}

func Request(requestID string, method string, path string, status int, duration time.Duration) {
	log.Println(
		fmt.Sprintf(
			"REQUEST request_id=%s method=%s path=%s status=%d duration=%s",
			requestID,
			method,
			path,
			status,
			duration.String(),
		),
	)
}
