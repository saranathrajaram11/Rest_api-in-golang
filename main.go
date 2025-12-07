package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

type quote struct {
	Id     int    `json:"id"`
	Text   string `json:"Text"`
	Author string `json:"Author"`
}

var db *sql.DB

func initDb() {
	var err error
	constr := "host=localhost port=5432 user=postgres dbname=mydb password=postgres sslmode=disable"
	db, err = sql.Open("postgres", constr)

	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("✓ Successfully connected to PostgreSQL!")

	// Create table if not exists
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS quotes (
		id SERIAL PRIMARY KEY,
		text TEXT NOT NULL,
		author TEXT NOT NULL
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		panic(err)
	}
	fmt.Println("✓ Table 'quotes' ready!")
}

func home_handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello iam this server writerr")
}

func add_quote(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "POST request only allowed", http.StatusMethodNotAllowed)
		return
	}
	var newquote quote
	err := json.NewDecoder(r.Body).Decode(&newquote)

	if err != nil {
		http.Error(w, "invalid json response", http.StatusBadRequest)
		return
	}

	// Insert into database
	query := "INSERT INTO quotes (text, author) VALUES ($1, $2) RETURNING id"
	err = db.QueryRow(query, newquote.Text, newquote.Author).Scan(&newquote.Id)

	if err != nil {
		http.Error(w, "failed to insert quote", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Quote added successfully!",
		"id":      newquote.Id,
	})
}

func update_quote(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "Put method only available", http.StatusMethodNotAllowed)
		return
	}
	id := r.URL.Query().Get("id")
	NumID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "error while conversion string to int", http.StatusBadRequest)
		return
	}

	var updatedQuote quote
	err = json.NewDecoder(r.Body).Decode(&updatedQuote)
	if err != nil {
		http.Error(w, "invalid json response", http.StatusBadRequest)
		return
	}

	// Update in database
	query := "UPDATE quotes SET text = $1, author = $2 WHERE id = $3"
	result, err := db.Exec(query, updatedQuote.Text, updatedQuote.Author, NumID)

	if err != nil {
		http.Error(w, "failed to update quote", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "quote not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "quote updated successfully",
	})
}

func delete_quote(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "DELETE only available", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	NumID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "conversion problem from string to int", http.StatusBadRequest)
		return
	}

	// Delete from database
	query := "DELETE FROM quotes WHERE id = $1"
	result, err := db.Exec(query, NumID)

	if err != nil {
		http.Error(w, "failed to delete quote", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "id is not found please type the correct id bruhh!", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "quote deleted successfully",
	})
}

func home_hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "this is Hello Page of this server")
}

func home_quote(w http.ResponseWriter, r *http.Request) {
	// Get random quote from database
	query := "SELECT id, text, author FROM quotes ORDER BY RANDOM() LIMIT 1"
	var q quote
	err := db.QueryRow(query).Scan(&q.Id, &q.Text, &q.Author)

	if err == sql.ErrNoRows {
		http.Error(w, "no quotes found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "failed to fetch quote", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(q)
}

func get_all_quotes(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "GET method only", http.StatusMethodNotAllowed)
		return
	}

	query := "SELECT id, text, author FROM quotes"
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "failed to fetch quotes", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var allQuotes []quote
	for rows.Next() {
		var q quote
		err := rows.Scan(&q.Id, &q.Text, &q.Author)
		if err != nil {
			http.Error(w, "failed to scan quote", http.StatusInternalServerError)
			return
		}
		allQuotes = append(allQuotes, q)
	}

	json.NewEncoder(w).Encode(allQuotes)
}

//main function to start the server
// it will listen on port 9094 and handle requests to the root path

func main() {
	initDb()
	defer db.Close()
	fmt.Println("localhost is running on http://localhost:9094")
	http.HandleFunc("/", home_handler)
	http.HandleFunc("/hello", home_hello)
	http.HandleFunc("/home_quote", home_quote)
	http.HandleFunc("/quotes", get_all_quotes)
	http.HandleFunc("/add_quote", add_quote)
	http.HandleFunc("/update_quote", update_quote)
	http.HandleFunc("/delete_quote", delete_quote)

	err := http.ListenAndServe(":9094", nil)

	if err != nil {
		fmt.Println("server failed", err)
	}
	fmt.Println("server started sucessfully")

}
