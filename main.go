package main

import (
	"log"
	"net/http"

	"github.com/djamysh/TracerApp/app"
	"github.com/gorilla/mux"
)

// expected formats
// string -> string
// percentage -> x element of [0,100]
// integer -> x element of Z
// float -> x element of R
// timelings -> map of timestamps with tags

// each activity must have the property timelings. As a default if the activity has
// a start and end time specific then it must have two key-pair in the timelings map,
// which are {'start':timestamp of start, 'end': timestamp of end}, if the activity data
// does not rely on start and end, then it must have 'instant' timestamp, which corresponds
// to the moment which data submitted. For example recording the cigarette consumption
// can be instantenous timeling while recording the thermodynamics study duration being
// an interval of start and end. This is for the default timelings property. Each activity
// will have either 'IntervalTimelings' property or 'InstantTimeling' property.
// After a brainstorm I think defining as Instant or Interval is restricting the flexability.
// *I think every Activity must have Default timeling property but not restricted as
// instant or interval.
// ** After a bit of coding I decided that Interval Timelings is a bit of optional, user additionally
// can also define it, I will automatically define the default timeling while creating the event in the collection.

// each activity must have the default Note property. I am not sure, it may be redundant.

// Get (Property&Activity) by name functions must be implemented

// I am thinking of removing percentage datatype, because if you consider
// all of the other datatypes they are well defined datatypes(timelings as map[string]int64)
// however percentage requires additional constraint of x element of [0,100]. It is is similar
// to that start-end timelings is some constraint of timelings. I will remove the percentage
// datatype, in future I am planning to add some built-in constraints for percentage, start-end
// intervals or similar usefull built-in features.

func main() {

	/*
		var err error
		// Check if default timelings properties are created
		// models.DefaultTimelingsPropertyID, err = models.DefaultTimelingsProperty()

		if err != nil {
			log.Fatal("ERROR while trying to obtain Default Timelings Property ID ")
			log.Fatal(err)

		}
	*/

	// Define the router
	r := mux.NewRouter()

	r.Use(app.SetHeaders)
	app.RegisterRoutes(r)

	// Start the server
	log.Fatal(http.ListenAndServe(":8000", r))
}
