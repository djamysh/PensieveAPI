package app

import (
	"github.com/gorilla/mux"

	"github.com/djamysh/PensieveAPI/services"
)

func RegisterRoutes(r *mux.Router) {

	// Add the routes
	r.HandleFunc("/activities", services.CreateActivityHandler).Methods("POST")
	r.HandleFunc("/activities/{id}", services.UpdateActivityHandler).Methods("PUT")
	r.HandleFunc("/activities/{id}", services.DeleteActivityHandler).Methods("DELETE")
	r.HandleFunc("/activities", services.GetActivitiesHandler).Methods("GET")
	r.HandleFunc("/activities/{id}", services.GetActivityHandler).Methods("GET")
	r.HandleFunc("/activities/ByName/{name}", services.GetActivityByNameHandler).Methods("GET")

	r.HandleFunc("/properties", services.CreatePropertyHandler).Methods("POST")
	r.HandleFunc("/properties/{id}", services.UpdatePropertyHandler).Methods("PUT")
	r.HandleFunc("/properties/{id}", services.DeletePropertyHandler).Methods("DELETE")
	r.HandleFunc("/properties", services.GetPropertiesHandler).Methods("GET")
	r.HandleFunc("/properties/{id}", services.GetPropertyHandler).Methods("GET")
	r.HandleFunc("/properties/ByName/{name}", services.GetPropertyByNameHandler).Methods("GET")

	r.HandleFunc("/events", services.CreateEventHandler).Methods("POST")
	r.HandleFunc("/events/{id}", services.GetEventHandler).Methods("GET")
	r.HandleFunc("/events/by/{activityID}", services.GetEventsByActivityID).Methods("GET")
	r.HandleFunc("/events", services.GetEventsHandler).Methods("GET")
	r.HandleFunc("/events/{id}", services.DeleteEventHandler).Methods("DELETE")
	r.HandleFunc("/events/{id}", services.UpdateEventHandler).Methods("PUT")
}
