package main

import (
	"context"
	"log"
	"net/http"

	df "github.com/kompiuter/dialogflow-go-webhook"
)

func main() {
	h := df.NewHandler()
	h.Register("room.get", roomGetHandler)
	h.Register("room.reserve", roomReserveHandler)

	log.Fatal(http.ListenAndServe(":5000", h))
}

func roomGetHandler(ctx context.Context, dfr *df.Request) (*df.Fulfillment, int) {
	// retrieve underlying http request if necessary
	// req := df.HTTPRequest(ctx)

	return &df.Fulfillment{
		FulfillmentText: "Your room is ...",
	}, http.StatusOK
}

func roomReserveHandler(ctx context.Context, dfr *df.Request) (*df.Fulfillment, int) {
	// retrieve room number param from dfr, perform your logic

	return &df.Fulfillment{
		FulfillmentText: "Your room has been reserved!",
	}, http.StatusOK
}
