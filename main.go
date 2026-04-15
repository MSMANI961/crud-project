package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

var db *sql.DB

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

func initDB() {
	var err error
	connStr := "user=postgres password=sql1234 dbname=albumsbd sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
}

// GET all
func getProducts(w http.ResponseWriter, r *http.Request) {
	rows, _ := db.Query("SELECT * FROM products")
	var products []Product

	for rows.Next() {
		var p Product
		rows.Scan(&p.ID, &p.Name, &p.Price)
		products = append(products, p)
	}
	json.NewEncoder(w).Encode(products)
}

// POST
func createProduct(w http.ResponseWriter, r *http.Request) {
	var p Product
	
	json.NewDecoder(r.Body).Decode(&p)

	err := db.QueryRow(
		"INSERT INTO products(name, price) VALUES($1,$2) RETURNING id",
		p.Name, p.Price).Scan(&p.ID)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(p)
	
}

func updateProduct(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id := params["id"]

	var p Product
	json.NewDecoder(r.Body).Decode(&p)

	log.Println("Updating:", id, p.Name, p.Price)

	_, err := db.Exec(
		"UPDATE products SET name=$1, price=$2 WHERE id=$3",
		p.Name, p.Price, id,
	)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode("updated")
}

// DELETE
func deleteProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
fmt.Println("ID received:", id) // Log the ID to verify it's being received correctly
	_, err := db.Exec("DELETE FROM products WHERE id=$1", id)
	
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode("deleted")
}

func main() {
	initDB()

	r := mux.NewRouter()

	r.HandleFunc("/products", getProducts).Methods("GET")
	r.HandleFunc("/products", createProduct).Methods("POST")
	r.HandleFunc("/products/{id}", updateProduct).Methods("PUT")
	r.HandleFunc("/products/{id}", deleteProduct).Methods("DELETE")
	

c := cors.New(cors.Options{
	AllowedOrigins: []string{"http://localhost:4200"},
	AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	AllowedHeaders: []string{"*"},
})

handler := c.Handler(r)

http.ListenAndServe(":8080", handler)

	
}