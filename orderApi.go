package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Order struct {
	OrderId      string    `json:"orderId"`
	CustomerName string    `json:"customerName"`
	OrderedAt    time.Time `json:orderedAt`
	Items        []Item    `json:items`
}

type Item struct {
	ItemId      string `json:"itemId"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
}

var orders []Order
var prevOrderId = 0

func main() {
	router := mux.NewRouter()

	//Create an Order
	router.HandleFunc("/orders", createOder).Methods("POST")
	//Get Order by ID
	router.HandleFunc("orders/{orderId}", getOrder).Methods("GET")
	//Getting All Orders from the database
	router.HandleFunc("/orders", getOrders).Methods("GET")
	//Update Orders in the database
	router.HandleFunc("/orders{orderId}", updateOrder).Methods("PUT")
	//Delet an Order from the database
	router.HandleFunc("orders/{orderId}", deleteOrder).Methods("DELETE")
	//swagger
	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	log.Fatal(http.ListenAndServe(":9090", router))
}

func createOder(w http.ResponseWriter, r *http.Request) {
	var order Order
	json.NewDecoder(r.Body).Decode(&order)
	prevOrderId++
	order.OrderId = strconv.Itoa(prevOrderId)
	orders = append(orders, order)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func getOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)

}

func getOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	inputOrderId := params["orderId"]
	for _, order := range orders {
		if order.OrderId == inputOrderId {
			json.NewEncoder(w).Encode(order)
			return
		}
	}
}

func updateOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	inputOrderId := params["orderId"]
	for i, order := range orders {
		if inputOrderId == order.OrderId {
			orders = append(orders[:i], orders[i+1:]...)
			var updatedOrder Order
			json.NewDecoder(r.Body).Decode(&updatedOrder)
			orders = append(orders, updatedOrder)
			json.NewEncoder(w).Encode(updatedOrder)
			return
		}
	}
}

func deleteOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	inputOrderId := params["orderId"]
	for i, order := range orders {
		if inputOrderId == order.OrderId {
			orders = append(orders[:i], orders[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
}
