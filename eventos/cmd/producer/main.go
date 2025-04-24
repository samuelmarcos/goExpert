package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/samuelmarscos/eventos/pkg/rabbitmq"
)

const (
	defaultMessageCount = 100000
	defaultExchange     = "amq.direct"
)

func main() {
	// Create a context that will be canceled on interrupt
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Initialize logger
	logger := log.New(os.Stdout, "[PRODUCER] ", log.LstdFlags)

	// Open RabbitMQ channel
	ch, err := rabbitmq.OpenChannel()
	if err != nil {
		logger.Fatalf("Failed to open channel: %v", err)
	}
	defer ch.Close()

	// Create a ticker for logging progress
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	// Start publishing messages
	logger.Printf("Starting to publish %d messages to exchange '%s'", defaultMessageCount, defaultExchange)

	var publishedCount int
	for i := 0; i < defaultMessageCount; i++ {
		select {
		case <-ctx.Done():
			logger.Printf("Received shutdown signal. Published %d messages", publishedCount)
			return
		case <-ticker.C:
			logger.Printf("Progress: %d/%d messages published", publishedCount, defaultMessageCount)
		default:
			msg := fmt.Sprintf("Message %d: Hello World!", i)
			if err := rabbitmq.Publish(ch, msg, defaultExchange); err != nil {
				logger.Printf("Failed to publish message %d: %v", i, err)
				continue
			}
			publishedCount++
		}
	}

	logger.Printf("Successfully published all %d messages", publishedCount)
}
