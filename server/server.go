package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ExchangeRate struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

var db *gorm.DB

func main() {
	initDB()
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", getExchangeHandler)

	fmt.Println("Servidor iniciado na porta 8080")
	http.ListenAndServe(":8080", mux)
}

func getExchangeHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	rate, err := getExchangeRate(ctx)
	if err != nil {
		http.Error(w, "Erro ao obter a cotação do dólar", http.StatusInternalServerError)
		fmt.Println("Erro ao obter a cotação do dólar:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rate)
}

func getExchangeRate(ctx context.Context) (*ExchangeRate, error) {
	resp, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data map[string]ExchangeRate
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	rate, ok := data["USDBRL"]
	if !ok {
		return nil, fmt.Errorf("Cotação não encontrada na resposta da API")
	}

	if err := saveExchangeRate(&rate); err != nil {
		fmt.Println("Erro ao salvar a cotação no banco de dados:", err)
	}

	return &rate, nil
}

func initDB() {
	database, err := gorm.Open(sqlite.Open("cotacoes.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("Erro ao abrir o banco de dados:", err)
		return
	}

	database.AutoMigrate(&ExchangeRate{})

	db = database
}

func saveExchangeRate(rate *ExchangeRate) error {
	result := db.Create(&rate)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
