// Package notify provides notification implementations.
package notify

import (
	"context"
	"fmt"
)

// Notifier is the interface for sending notifications.
type Notifier interface {
	Send(ctx context.Context, message string) error
}

// ConsoleNotifier prints notifications to stdout.
type ConsoleNotifier struct{}

// NewConsole returns a new ConsoleNotifier.
func NewConsole() *ConsoleNotifier {
	return &ConsoleNotifier{}
}

func (c *ConsoleNotifier) Send(ctx context.Context, message string) error {
	fmt.Printf("\n[NOTIFICATION] %s\n\n", message)
	return nil
}
