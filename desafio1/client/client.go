package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

type CotacaoResume struct {
	ID  int `gorm:"primaryKey"`
	Bid string
}

type DollarInfo struct {
	ID    string
	Dolar string
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
			dolar VARCHAR(10)
		)`,
	)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func main() {

	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create("cotacao.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var resume CotacaoResume
	err = json.Unmarshal(body, &resume)

	_, err = file.WriteString(fmt.Sprintf("Dólar: %s", resume.Bid))
	if err != nil {
		log.Fatal(err)
	}

	var dollarInfo = NewDolarInfo(resume.Bid)

	err = SaveDollarInfo(db, dollarInfo, ctx)

	fmt.Println(body)
}

func SaveDollarInfo(db *sql.DB, resume *DollarInfo, ctx context.Context) error {
	_, err := db.ExecContext(
		ctx,
		"insert into cotacoes(id,bid) values (?,?), resume.ID, resume.Dolar",
	)

	if err != nil {
		return err
	}

	return nil
}

func NewDolarInfo(bid string) *DollarInfo {
	return &DollarInfo{
		ID:    uuid.New().String(),
		Dolar: bid,
	}
}
