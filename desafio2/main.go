package main

import (
	"context"
	"encoding/json"
	"fmt"
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
	url  string
	body CepResponse
}

func main() {
	urls := []string{
		"https://brasilapi.com.br/api/cep/v1/01153000 + 35402176",
		"http://viacep.com.br/ws/35402176/json/",
	}

	chanResults := make(chan *Response, len(urls))
	chanErros := make(chan error, len(urls))
	var wg sync.WaitGroup
	wg.Add(len(urls))

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	go makeRequest(urls[0], chanResults, chanErros, ctx, &wg)
	go makeRequest(urls[1], chanResults, chanErros, ctx, &wg)

	wg.Wait()

	select {
	case res := <-chanResults:
		fmt.Printf("Url %s obteve response primeiro, com resposta %v :", res.url, res.body)
		cancel()
	case err := <-chanErros:
		fmt.Printf("Error %v !! ", err)
		return

	case <-time.After(1 * time.Second):
		fmt.Printf("Timeout atingido ao realizar request")
		return

	case <-ctx.Done():
		fmt.Printf("Timeout atingido ao realizar request")
		return
	}
}

func makeRequest(url string, chanResults chan<- *Response, chanErros chan<- error, ctx context.Context, wg *sync.WaitGroup) {

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		chanErros <- err
	}

	defer resp.Body.Close()

	defer wg.Done()

	var cepResp CepResponse

	err = json.NewDecoder(resp.Body).Decode(&cepResp)
	if err != nil {
		chanErros <- err
	}
	response := &Response{
		url:  url,
		body: cepResp,
	}
	chanResults <- response

}
