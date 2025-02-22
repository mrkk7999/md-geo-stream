package mdgeotrack

type Repository interface{
	// Health check
	HeartBeat() map[string]string
}
