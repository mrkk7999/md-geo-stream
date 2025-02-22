package kafka

import (
	"context"
	"log"
	mdgeotrack "md-geo-track/iface"

	"github.com/IBM/sarama"
)

// StartConsumer initializes and starts a Kafka consumer group
func StartConsumer(brokers []string, topic string, groupID string, service mdgeotrack.Service) {
	// Create a new Sarama configuration
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	// Create a new consumer group
	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		log.Fatalf("Failed to start Sarama consumer group: %v", err)
	}

	// Initialize the consumer
	consumer := Consumer{
		ready:   make(chan bool),
		service: service,
	}

	// Create a new context
	ctx := context.Background()

	// Start the consumer group in a new goroutine
	go func() {
		for {
			// Consume messages from the topic
			if err := consumerGroup.Consume(ctx, []string{topic}, &consumer); err != nil {
				log.Fatalf("Error from consumer: %v", err)
			}
			// Check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			// Reset the ready channel
			consumer.ready = make(chan bool)
		}
	}()

	// Wait until the consumer is ready
	<-consumer.ready
	log.Println("Sarama consumer up and running!")
}

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	ready   chan bool
	service mdgeotrack.Service
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Signal that the consumer is ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// Loop over the messages in the claim
	for message := range claim.Messages() {
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
		// Process the message
		consumer.service.ProcessLocationData(message)
		// Mark the message as processed
		session.MarkMessage(message, "")
	}
	return nil
}
