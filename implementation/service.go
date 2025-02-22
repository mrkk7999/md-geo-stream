package implementation

import mdgeotrack "md-geo-track/iface"

type service struct {
	repository mdgeotrack.Repository
}

func New(repository mdgeotrack.Repository) mdgeotrack.Service {
	return &service{
		repository: repository,
	}
}
