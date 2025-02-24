package implementation

import (
	"encoding/json"
	logger "log"
	"md-geo-track/request_response/location"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

// ProcessLocationData
// Unmarshal message consumed, sends it to third party
func (s *service) ProcessLocationData(message *sarama.ConsumerMessage) {
	var locReq location.LocationReq
	if err := json.Unmarshal(message.Value, &locReq); err != nil {
		s.log.WithError(err).Error("could not unmarshal message")
		return
	}

	s.log.WithFields(logrus.Fields{
		"location": locReq,
	}).Info("received location data")

	// Simulate sending data to a third party for analysis
	s.sendToThirdParty(locReq)
}

func (s *service) sendToThirdParty(locReq location.LocationReq) {
	// Simulate sending data to a third party
	s.log.Info("sending data to third party")
	logger.Println("---------------------------------------------")
}
