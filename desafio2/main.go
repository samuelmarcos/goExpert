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

func main() {
	urls := []string{
		"https://brasilapi.com.br/api/cep/v1/01153000 + 35402176",
		"http://viacep.com.br/ws/35402176/json/",
	}

	chanResults := make(chan *CepResponse, len(urls))
	chanErros := make(chan error, len(urls))
	var wg sync.WaitGroup
	wg.Add(len(urls))

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	for _, url := range urls {
		go makeRequest(url, chanResults, chanErros, ctx, &wg)
	}

	wg.Wait()

	for _, url := range urls {
		select {
		case res := <-chanResults:
			fmt.Printf("Url %s obteve response primeiro :", url)
			fmt.Println(res)
			return
		case err := <-chanErros:
			fmt.Printf("Url %s obteve response de error :", url)
			fmt.Println(err)
			return

		case <-ctx.Done():
			fmt.Printf("Timeout atingido ao realizar request em %s :", url)
			return
		}

	}
}

func makeRequest(url string, chanResults chan<- *CepResponse, chanErros chan<- error, ctx context.Context, wg *sync.WaitGroup) {

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
	chanResults <- &cepResp

}
