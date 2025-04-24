package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/samuelmarscos/eventos/pkg/rabbitmq"
	"github.com/streadway/amqp"
)

const (
	defaultQueueName = "orders"
	prefetchCount    = 10 // Number of messages to prefetch
)

func main() {
	// Create a context that will be canceled on interrupt
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Initialize logger
	logger := log.New(os.Stdout, "[CONSUMER] ", log.LstdFlags)

	// Open RabbitMQ channel
	ch, err := rabbitmq.OpenChannel()
	if err != nil {
		logger.Fatalf("Failed to open channel: %v", err)
	}
	defer ch.Close()

	// Set QoS to control how many messages are consumed at once
	err = ch.Qos(
		prefetchCount, // prefetch count
		0,             // prefetch size
		false,         // global
	)
	if err != nil {
		logger.Fatalf("Failed to set QoS: %v", err)
	}

	// Create message channel
	msgs := make(chan amqp.Delivery)

	// Start consuming messages
	logger.Printf("Starting to consume messages from queue '%s'", defaultQueueName)
	go rabbitmq.Consume(ch, msgs, defaultQueueName)

	// Create a ticker for logging
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// Process messages
	var processedCount int
	for {
		select {
		case <-ctx.Done():
			logger.Printf("Received shutdown signal. Processed %d messages", processedCount)
			return
		case <-ticker.C:
			logger.Printf("Consumer is running. Processed %d messages so far", processedCount)
		case msg, ok := <-msgs:
			if !ok {
				logger.Println("Message channel closed")
				return
			}

			// Process the message
			if err := processMessage(logger, msg); err != nil {
				logger.Printf("Failed to process message: %v", err)
				// Reject the message and don't requeue it
				if err := msg.Reject(false); err != nil {
					logger.Printf("Failed to reject message: %v", err)
				}
				continue
			}
			processedCount++

			// Acknowledge the message only if processing was successful
			if err := msg.Ack(false); err != nil {
				logger.Printf("Failed to acknowledge message: %v", err)
			}
		}
	}
}

func processMessage(logger *log.Logger, msg amqp.Delivery) error {
	// Add your message processing logic here
	logger.Printf("Received message: %s", string(msg.Body))

	// Simulate some processing time
	time.Sleep(100 * time.Millisecond)

	// Return nil if processing was successful
	// Return an error if processing failed
	return nil
}
