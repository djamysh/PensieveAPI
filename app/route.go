package app

import (
	"github.com/gorilla/mux"

	"github.com/djamysh/PensieveAPI/services"
)

func RegisterRoutes(r *mux.Router) {

	// Add the routes
	r.HandleFunc("/activities", services.CreateActivityHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/activities/{id}", services.UpdateActivityHandler).Methods("PUT", "OPTIONS")
	r.HandleFunc("/activities/{id}", services.DeleteActivityHandler).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/activities", services.GetActivitiesHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/activities/{id}", services.GetActivityHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/activities/ByName/{name}", services.GetActivityByNameHandler).Methods("GET", "OPTIONS")

	r.HandleFunc("/properties", services.CreatePropertyHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/properties/{id}", services.UpdatePropertyHandler).Methods("PUT", "OPTIONS")
	r.HandleFunc("/properties/{id}", services.DeletePropertyHandler).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/properties", services.GetPropertiesHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/properties/{id}", services.GetPropertyHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/properties/ByName/{name}", services.GetPropertyByNameHandler).Methods("GET", "OPTIONS")

	r.HandleFunc("/events", services.CreateEventHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/events/{id}", services.GetEventHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/events/by/{activityID}", services.GetEventsByActivityID).Methods("GET", "OPTIONS")
	r.HandleFunc("/events", services.GetEventsHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/events/{id}", services.DeleteEventHandler).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/events/{id}", services.UpdateEventHandler).Methods("PUT", "OPTIONS")
}
