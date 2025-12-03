package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
)

type quote struct {
	Id     int    `json:"id"`
	Text   string `json:"Text"`
	Author string `json:"Author"`
}

var quotes = []quote{
	{1, "believe yourself", "saranath"},
	{2, "mind your silence", "rajaram"},
	{3, "everyone is not powerful", "rambha"},
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
	nextID := 1

	for _, q := range quotes {
		if q.Id >= nextID {
			nextID = q.Id + 1
		}
	}
	newquote.Id = nextID

	quotes = append(quotes, newquote)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Quote added successfully!",
	})

}

func delete_quote(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "DELETE only available ", http.StatusMethodNotAllowed)
		return
	}

	newQuotes = []quotes{}
	id := r.URL.Query().Get("Id")
	NumID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "conversion problem from string to int ", http.StatusBadRequest)
		return
	}
	for _, q := range quotes {
		if NumID != q.id {
			newQuotes := append(newQuotes, quotes)

		}
	}

}

func home_hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "this is Hello Page of this server")
}

func home_quote(w http.ResponseWriter, r *http.Request) {

	random_q := quotes[rand.Intn(len(quotes))]
	json.NewEncoder(w).Encode(random_q)
	json.NewEncoder(w).Encode(quotes)
	fmt.Fprintln(w)
}

//main function to start the server
// it will listen on port 9094 and handle requests to the root path

func main() {
	fmt.Println("localhost is running on http://localhost:9094")
	http.HandleFunc("/", home_handler)
	http.HandleFunc("/hello", home_hello)
	http.HandleFunc("/home_quote", home_quote)
	http.HandleFunc("/add_quote", add_quote)
	http.HandleFunc("/delete_quote", delete_quote)

	err := http.ListenAndServe(":9094", nil)

	if err != nil {
		fmt.Println("server failed", err)
	}
	fmt.Println("server started sucessfully")

}
