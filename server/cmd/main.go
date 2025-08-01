package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Cotacao struct {
	Bid string `json:"bid"`
}

type ResponseAPI struct {
	USDBRL Cotacao `json:"USDBRL"`
}

var db *sql.DB
var err error


func main() {
	
	db, err = sql.Open("sqlite3", "./cotation.db")
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
	url := "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	ctx, cancel := context.WithTimeout(r.Context(),  200 * time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "Erro ao obter cotação", http.StatusGatewayTimeout)
		fmt.Println("Erro na requisição externa:", err)
		return
	}

	var apiResp ResponseAPI
	if err := json.NewDecoder(res.Body).Decode(&apiResp); err != nil {
		http.Error(w, "Erro ao decodificar resposta", http.StatusInternalServerError)
		fmt.Println("Erro ao decodificar JSON:", err)
		return
	}
	fmt.Println("Cotação obtida:", apiResp.USDBRL.Bid)
	defer res.Body.Close()

	ctxDB, cancelDB := context.WithTimeout(r.Context(), 10 * time.Millisecond)
	defer cancelDB()
	err = saveCotacao(ctxDB, apiResp.USDBRL.Bid)
	if err != nil {
		fmt.Println("Erro ao salvar no banco:", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"bid": apiResp.USDBRL.Bid})
	
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

func saveCotacao(ctx context.Context, bid string) error {
	stmt, err := db.PrepareContext(ctx, "INSERT INTO cotacoes(bid) VALUES(?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, bid)
	return err
}
