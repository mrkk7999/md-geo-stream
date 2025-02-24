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
	err := s.sendToThirdParty(locReq)
	if err != nil {
		s.log.Error("Error sending data to third party")
		return
	}
	// will do it as queue/log based
	// will process this too asynchronously
	err = s.repository.UpdateLocationStatus(locReq)
	if err != nil {
		s.log.Error("Error updating location status in db")
		return
	}
	s.log.Info("Location end-to-end operation complete")
}

func (s *service) sendToThirdParty(locReq location.LocationReq) error {
	// Simulate sending data to a third party
	s.log.Info("sending data to third party")
	logger.Println("---------------------------------------------")
	return nil
}
