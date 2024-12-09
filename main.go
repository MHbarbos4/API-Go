package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "go-api/docs" // Import gerado automaticamente pelo swag
	httpSwagger "github.com/swaggo/http-swagger"
)

// Estrutura para o item
type Item struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// Handlers

// @Summary Get all items
// @Description Retrieve all items in the list
// @Tags items
// @Produce json
// @Success 200 {array} Item
// @Router /items [get]
func getItems(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, price FROM items")
	if err != nil {
		http.Error(w, "Error fetching items", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Price); err != nil {
			http.Error(w, "Error reading items", http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}

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
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var item Item
	err = db.QueryRow("SELECT id, name, price FROM items WHERE id = ?", id).Scan(&item.ID, &item.Name, &item.Price)
	if err == sql.ErrNoRows {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error fetching item", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

// @Summary Create an item
// @Description Add a new item to the database
// @Tags items
// @Accept json
// @Produce json
// @Param item body Item true "New Item"
// @Success 201 {object} Item
// @Router /items [post]
func createItem(w http.ResponseWriter, r *http.Request) {
	var newItem Item
	if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	result, err := db.Exec("INSERT INTO items (name, price) VALUES (?, ?)", newItem.Name, newItem.Price)
	if err != nil {
		http.Error(w, "Error inserting item", http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	newItem.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
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
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updatedItem Item
	if err := json.NewDecoder(r.Body).Decode(&updatedItem); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	result, err := db.Exec("UPDATE items SET name = ?, price = ? WHERE id = ?", updatedItem.Name, updatedItem.Price, id)
	if err != nil {
		http.Error(w, "Error updating item", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	updatedItem.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedItem)
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
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	result, err := db.Exec("DELETE FROM items WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Error deleting item", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @title Item API
// @version 1.0
// @description This is a simple API for managing items using SQLite.
// @host localhost:8080
// @BasePath /
func main() {
	// Inicializar banco de dados
	initDB()
	defer db.Close()

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