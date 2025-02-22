package repository

func (r *repository) HeartBeat() map[string]string {
	return map[string]string{
		"message": "Service is up and running",
	}
}
