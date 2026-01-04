// Package notify provides notification implementations.
package notify

import (
	"context"
)

// Notifier is the interface for sending notifications.
type Notifier interface {
	Send(ctx context.Context, message string) error
}
