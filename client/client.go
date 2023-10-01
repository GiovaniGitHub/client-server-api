package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	clientTimeout := 300 * time.Millisecond

	_, cancel := context.WithTimeout(context.Background(), clientTimeout)
	defer cancel()

	resp, err := http.Get("http://localhost:8080/cotacao")
	if err != nil {
		fmt.Println("Erro ao fazer a requisição HTTP:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erro ao ler a resposta da requisição:", err)
		return
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Erro ao decodificar a resposta JSON:", err)
		return
	}

	bid, ok := data["bid"].(string)
	if !ok {
		fmt.Println("Erro: campo 'bid' não encontrado na resposta")
		return
	}

	fileName := "cotacao.txt"
	fileContent := fmt.Sprintf("Dólar: %s\n", bid)
	err = os.WriteFile(fileName, []byte(fileContent), 0644)
	if err != nil {
		fmt.Println("Erro ao salvar o valor da cotação em cotacao.txt:", err)
		return
	}

	fmt.Println("Cotação do Dólar salva em cotacao.txt:", bid)
}
