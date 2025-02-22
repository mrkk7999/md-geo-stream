package http

import (
	"md-geo-track/controller"
	"net/http"

	"github.com/gorilla/mux"
)

func SetUpRouter(controller *controller.Controller) http.Handler {
	var (
		router = mux.NewRouter()
	)
	router.HandleFunc("/api/v1/heartbeat", controller.HeartBeatHandler).Methods("GET")

	return router
}
