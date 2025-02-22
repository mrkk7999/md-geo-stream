package implementation

func (s *service) HeartBeat() map[string]string {
	return s.repository.HeartBeat()
}
