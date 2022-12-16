package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)

	if err != nil {
		fmt.Println("Error on request")
		return
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println("Something bad happened")
		return
	}

	defer res.Body.Close()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao ler a resposta: %v\n", err)
		return
	}

	if res.StatusCode != http.StatusOK {
		fmt.Println(res.StatusCode, "The server is taking to long to respond")
		return
	}

	file, err := os.Create("cotacoes.txt")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao criar o arquivo: %v\n", err)
	}

	defer file.Close()

	contentOfRequest, err := io.ReadAll(res.Body)

	file.WriteString(fmt.Sprintf("Dólar: %v", string(contentOfRequest)))
	fmt.Println("Arquivo criado com sucesso!")
	fmt.Println("Dólar: ", string(contentOfRequest))
}
