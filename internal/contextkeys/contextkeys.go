package contextkeys

import (
	"context"
)

type ContextKey string

const RequestIDKey ContextKey = "request_id"

func GetRequestID(ctx context.Context) string {
	value := ctx.Value(RequestIDKey)

	requestID, ok := value.(string)
	if !ok {
		return "unknown"
	}

	return requestID
}
