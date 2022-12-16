package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"
)

type RestClient struct {
	Url string
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
