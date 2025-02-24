package implementation

import (
	mdgeotrack "md-geo-track/iface"

	"github.com/sirupsen/logrus"
)

type service struct {
	repository mdgeotrack.Repository
	log        *logrus.Logger
}

func New(repository mdgeotrack.Repository, log *logrus.Logger) mdgeotrack.Service {
	return &service{
		repository: repository,
		log:        log,
	}
}
