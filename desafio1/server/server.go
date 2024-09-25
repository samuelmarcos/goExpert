package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

type Cotacao struct {
	Code        string `json:"code"`
	Codein      string `json:"codein"`
	Name        string `json:"name"`
	High        string `json:"high"`
	Low         string `json:"low"`
	VarBid      string `json:"varBid"`
	PctChange   string `json:"pctChange"`
	Bid         string `json:"bid"`
	Ask         string `json:"ask"`
	Timestamp   string `json:"timestamp"`
	Create_date string `json:"create_date"`
}

type CotacaoResume struct {
	ID  string
	Bid string
}

func main() {
	http.HandleFunc("/cotacao", BuscaCotacaoHandle)
	http.ListenAndServe(":8080", nil)

}

func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "cotacao.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.ExecContext(
		context.Background(),
		`DROP TABLE IF EXISTS cotacoes;
		 CREATE TABLE cotacoes (
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			bid VARCHAR(10)
		)`,
	)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func NewCotacaoResume(bid string) *CotacaoResume {
	return &CotacaoResume{
		ID:  uuid.New().String(),
		Bid: bid,
	}
}

func BuscaCotacaoHandle(w http.ResponseWriter, r *http.Request) {

	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}

	cotacao, err := BuscaCotacao()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	resume := NewCotacaoResume(cotacao.Bid)

	ctx := context.Background()
	err = SaveCotacao(db, resume, ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cotacao)

}

func SaveCotacao(db *sql.DB, resume *CotacaoResume, ctx context.Context) error {
	_, err := db.ExecContext(
		ctx,
		"insert into cotacoes(id,bid) values (?,?), resume.ID, resume.Bid",
	)

	if err != nil {
		return err
	}

	return nil
}

func BuscaCotacao() (*Cotacao, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)

	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}
	var cotacao Cotacao
	err = json.Unmarshal(body, &cotacao)

	if err != nil {
		return nil, err
	}

	return &cotacao, nil
}
