package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type CepResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Estado      string `json:"estado"`
	Regiao      string `json:"regiao"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type Response struct {
	URL  string
	Body CepResponse
}

func main() {
	urls := []string{
		"https://brasilapi.com.br/api/cep/v1/01153000",
		"http://viacep.com.br/ws/01153000/json/",
	}

	// Create buffered channels
	chanResults := make(chan *Response, len(urls))
	chanErrors := make(chan error, len(urls))

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(len(urls))

	for _, url := range urls {
		go makeRequest(url, chanResults, chanErrors, ctx, &wg)
	}

	go func() {
		wg.Wait()
		close(chanResults)
		close(chanErrors)
	}()

	select {
	case res, ok := <-chanResults:
		if !ok {
			fmt.Println("Channel closed without results")
			return
		}
		fmt.Printf("First response received from %s:\n", res.URL)
		fmt.Printf("CEP: %s\n", res.Body.Cep)
		fmt.Printf("Logradouro: %s\n", res.Body.Logradouro)
		fmt.Printf("Bairro: %s\n", res.Body.Bairro)
		fmt.Printf("Localidade: %s\n", res.Body.Localidade)
		fmt.Printf("UF: %s\n", res.Body.Uf)
		cancel()

	case err, ok := <-chanErrors:
		if !ok {
			fmt.Println("Channel closed without errors")
			return
		}
		fmt.Printf("Error occurred: %v\n", err)
		return

	case <-ctx.Done():
		fmt.Println("Timeout reached while waiting for responses")
		return
	}
}

func makeRequest(url string, chanResults chan<- *Response, chanErrors chan<- error, ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		chanErrors <- fmt.Errorf("error creating request for %s: %v", url, err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		chanErrors <- fmt.Errorf("error making request to %s: %v", url, err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		chanErrors <- fmt.Errorf("error reading response body from %s: %v", url, err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		chanErrors <- fmt.Errorf("unexpected status code %d from %s", resp.StatusCode, url)
		return
	}

	var cepResp CepResponse
	if err := json.Unmarshal(body, &cepResp); err != nil {
		chanErrors <- fmt.Errorf("error unmarshaling response from %s: %v", url, err)
		return
	}

	response := &Response{
		URL:  url,
		Body: cepResp,
	}
	chanResults <- response
}
