package mdgeotrack

import (
	"github.com/IBM/sarama"
)

type Service interface {
	// Health check
	HeartBeat() map[string]string

	ProcessLocationData(message *sarama.ConsumerMessage)
}
