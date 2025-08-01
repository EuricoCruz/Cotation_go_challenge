package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Cotacao struct {
	Bid string `json:"bid"`
}

func main() {
	fmt.Println("Cotation Go Challenge Client")
	ctx, cancel := context.WithTimeout(context.Background(), 300 *time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error response from server:", resp.Status)
		return
	}
	var cotacao Cotacao
	if err := json.NewDecoder(resp.Body).Decode(&cotacao); err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}
	
	text := fmt.Sprintf("Dólar: %s", cotacao.Bid)
	err = os.WriteFile("cotacao.txt", []byte(text), 0644)
	if err != nil {
		fmt.Printf("Erro ao escrever arquivo: %v\n", err)
	}

	fmt.Println("Cotação salva com sucesso!")
}