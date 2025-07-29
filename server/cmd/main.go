package main

import (
	"database/sql"
	"fmt"
	"net/http"
	_ "github.com/mattn/go-sqlite3"

)


func main() {

//	url := "https://economia.awesomeapi.com.br/json/last/USD-BRL"

	db, err := sql.Open("sqlite3", "./cotation.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	createTable(db)
	fmt.Println("Cotation Go Challenge Server")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path == "/" {
			fmt.Fprintf(w, "Welcome to the Cotation Go Challenge!")
		}
	})
	http.HandleFunc("/cotacao", handleCotation)
	http.ListenAndServe(":8080", nil)

}

func handleCotation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintf(w, "Cotation endpoint is under construction")
}

func createTable(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS cotacoes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		bid TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		fmt.Printf("Erro ao criar tabela: %v", err)
	}
}