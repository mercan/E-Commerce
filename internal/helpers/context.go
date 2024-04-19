package helpers

import (
	"context"
	"time"
)

// ContextWithTimeout returns a context with a timeout
func ContextWithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout*time.Second) // Süreyi saniye cinsinden çarp
}
