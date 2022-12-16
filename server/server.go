package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

type Cotacao struct {
	Usdbrl struct {
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
	} `json:"USDBRL"`
}

type RestClient struct {
	Url string
}

func main() {

	mux := CreateNewMuxAndAttachHandlers()
	port := ":8080"

	log.Println("Server is running on port", port)
	log.Fatal(http.ListenAndServe(port, mux))
}

func CreateNewMuxAndAttachHandlers() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", cotacaoHandler)
	return mux
}

func cotacaoHandler(w http.ResponseWriter, r *http.Request) {
	restClient := NewRestClient("https://economia.awesomeapi.com.br/json/last/USD-BRL")

	body, statusCode := restClient.Get()

	if statusCode != http.StatusOK {
		w.WriteHeader(statusCode)
		return
	}

	var cotacao Cotacao
	err := json.Unmarshal(body, &cotacao)

	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	db, err := sql.Open("sqlite3", "../db/cotacoes.db")

	if err != nil {
		log.Println(err.Error())
		panic(err)
	}

	defer db.Close()

	err = insertCotacao(db, &cotacao)

	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)
	w.Write([]byte(cotacao.Usdbrl.Bid))
}

func NewRestClient(url string) *RestClient {
	return &RestClient{Url: url}
}

func (r RestClient) Get() ([]byte, int) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", r.Url, nil)

	if err != nil {
		log.Println(err.Error())
		return nil, http.StatusInternalServerError
	}

	res, err := http.DefaultClient.Do(req)

	select {
	case <-ctx.Done():
		log.Println("[STATUS 408 (REQUEST TIMEOUT)]: The server is taking to long to respond.")
		return nil, http.StatusRequestTimeout
	case <-time.After(100 * time.Millisecond):
		if err != nil {
			log.Println(err.Error())
			return nil, http.StatusInternalServerError
		}

		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)

		if err != nil {
			log.Println(err.Error())
			return nil, http.StatusInternalServerError
		}
		log.Println("[STATUS 200 (OK)]: Request successfully completed.")
		return body, http.StatusOK
	}
}

func insertCotacao(db *sql.DB, cotacao *Cotacao) error {
	query := "INSERT INTO cotacoes(id, code, codein, name, high, low, varBid, pct_change, bid, ask, timestamp, create_date) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?,?, ?)"

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()

	_, err := db.ExecContext(
		ctx,
		query,
		uuid.New().String(),
		cotacao.Usdbrl.Code,
		cotacao.Usdbrl.Codein,
		cotacao.Usdbrl.Name,
		cotacao.Usdbrl.High,
		cotacao.Usdbrl.Low,
		cotacao.Usdbrl.VarBid,
		cotacao.Usdbrl.PctChange,
		cotacao.Usdbrl.Bid,
		cotacao.Usdbrl.Ask,
		cotacao.Usdbrl.Timestamp,
		cotacao.Usdbrl.CreateDate,
	)

	if err != nil {
		return err
	}

	return nil
}
