package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
	_ "go-api/docs" // Import gerado automaticamente pelo swag
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
)

// Estrutura para o item
type Item struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

var items []Item
var nextID int = 1

// Handlers

// @Summary Get all items
// @Description Retrieve all items in the list
// @Tags items
// @Produce json
// @Success 200 {array} Item
// @Router /items [get]
func getItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// @Summary Get an item by ID
// @Description Retrieve a single item by ID
// @Tags items
// @Produce json
// @Param id path int true "Item ID"
// @Success 200 {object} Item
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Item not found"
// @Router /items/{id} [get]
func getItemByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for _, item := range items {
		if item.ID == id {
			json.NewEncoder(w).Encode(item)
			return
		}
	}

	http.Error(w, "Item not found", http.StatusNotFound)
}

// @Summary Create an item
// @Description Add a new item to the list
// @Tags items
// @Accept json
// @Produce json
// @Param item body Item true "New Item"
// @Success 201 {object} Item
// @Router /items [post]
func createItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newItem Item
	json.NewDecoder(r.Body).Decode(&newItem)

	newItem.ID = nextID
	nextID++

	items = append(items, newItem)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newItem)
}

// @Summary Update an item
// @Description Update an existing item by ID
// @Tags items
// @Accept json
// @Produce json
// @Param id path int true "Item ID"
// @Param item body Item true "Updated Item"
// @Success 200 {object} Item
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Item not found"
// @Router /items/{id} [put]
func updateItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updatedItem Item
	json.NewDecoder(r.Body).Decode(&updatedItem)
	for i, item := range items {
		if item.ID == id {
			updatedItem.ID = item.ID
			items[i] = updatedItem
			json.NewEncoder(w).Encode(updatedItem)
			return
		}
	}

	http.Error(w, "Item not found", http.StatusNotFound)
}

// @Summary Delete an item
// @Description Remove an item by ID
// @Tags items
// @Param id path int true "Item ID"
// @Success 204 {string} string "No Content"
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Item not found"
// @Router /items/{id} [delete]
func deleteItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for i, item := range items {
		if item.ID == id {
			items = append(items[:i], items[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Item not found", http.StatusNotFound)
}

// @title Item API
// @version 1.0
// @description This is a simple API for managing items.
// @host localhost:8080
// @BasePath /
func main() {
	r := mux.NewRouter()

	// Swagger
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Roteamento
	r.HandleFunc("/items", getItems).Methods("GET")
	r.HandleFunc("/items/{id:[0-9]+}", getItemByID).Methods("GET")
	r.HandleFunc("/items", createItem).Methods("POST")
	r.HandleFunc("/items/{id:[0-9]+}", updateItem).Methods("PUT")
	r.HandleFunc("/items/{id:[0-9]+}", deleteItem).Methods("DELETE")

	// Exibindo a mensagem no terminal
	log.Println("Servidor rodando na porta :8080")

	// Iniciando o servidor
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("Erro ao iniciar o servidor: ", err)
	}
}