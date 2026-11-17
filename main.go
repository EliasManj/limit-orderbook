package main

import (
	"log"
	"net/http"

	"github.com/EliasManj/orderbook/api"
	"github.com/gorilla/mux"
)

func main() {
	api.SeedOrders()

	r := mux.NewRouter()
	r.HandleFunc("/bids/", api.GetBids)
	r.HandleFunc("/asks/", api.GetAsks)
	r.HandleFunc("/order/", api.CreateOrder).Methods("POST")
	r.HandleFunc("/clear/", api.ClearOrders).Methods("POST")

	// Serve static HTML file
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", r))
}
