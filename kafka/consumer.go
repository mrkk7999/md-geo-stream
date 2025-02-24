package kafka

import (
	"context"
	mdgeotrack "md-geo-track/iface"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

// StartConsumer initializes and starts a Kafka consumer group
func StartConsumer(brokers []string, topic string, groupID string, service mdgeotrack.Service, log *logrus.Logger) {
	// Create a new Sarama configuration
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	// Create a new consumer group
	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		log.Errorf("Failed to start Sarama consumer group: %v", err)
		return
	}

	// Initialize the consumer
	consumer := Consumer{
		ready:   make(chan bool),
		service: service,
		log:     log,
	}

	// Create a new context
	ctx := context.Background()

	// Start the consumer group in a new goroutine
	go func() {
		for {
			// Consume messages from the topic
			if err := consumerGroup.Consume(ctx, []string{topic}, &consumer); err != nil {
				log.Errorf("Error from consumer: %v", err)
				return
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
	log.Info("Sarama consumer up and running!")
}

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	ready   chan bool
	service mdgeotrack.Service
	log     *logrus.Logger
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Signal that the consumer is ready
	close(consumer.ready)
	consumer.log.Info("Consumer session setup complete")
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	consumer.log.Info("Consumer session cleanup complete")
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// Loop over the messages in the claim
	for message := range claim.Messages() {
		// consumer.log.Infof("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
		consumer.log.Infof("Message claimed")

		// Process the message
		consumer.service.ProcessLocationData(message)
		// Mark the message as processed
		session.MarkMessage(message, "")
	}
	return nil
}
