package implementation

import (
	"encoding/json"
	"log"
	"md-geo-track/request_response/location"

	"github.com/IBM/sarama"
)

// ProcessLocationData
// Unmarshal message consumed, sends it to third party
func (s *service) ProcessLocationData(message *sarama.ConsumerMessage) {
	var locReq location.LocationReq
	if err := json.Unmarshal(message.Value, &locReq); err != nil {
		log.Printf("could not unmarshal message: %v", err)
		return
	}

	log.Printf("received location data: %+v", locReq)

	// Simulate sending data to a third party for analysis
	s.sendToThirdParty(locReq)
}

func (s *service) sendToThirdParty(locReq location.LocationReq) {
	// Simulate sending data to a third party
	log.Printf("sending data to third party: %+v", locReq)
}
