package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type CurrencyExchange struct {
	USDBRL CurrencyDetails `json:"USDBRL"`
}

type CurrencyDetails struct {
	ID         uint   `gorm:"primaryKey"`
	Code       string `json:"code"`
	CodeIn     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
	gorm.Model
}

var db *gorm.DB
var logger = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

func main() {
	var err error

	db, err = gorm.Open(sqlite.Open("currency_exchange.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Erro ao conectar ao banco de dados:", err)
	}

	db.AutoMigrate(&CurrencyDetails{})

	http.HandleFunc("/cotacao", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	currencyExchange, err := getDollarExchangeRate(ctx)
	logContextTimeout(ctx.Err(), "A requisição foi cancelada devido ao timeout de 200ms")

	if err != nil {
		logger.Println("Ocorreu um erro ao consultar os dados na API: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	err = saveIntoDatabase(ctx, currencyExchange.USDBRL)

	if err != nil {
		logger.Println("Ocorreu um erro ao gravar os dados: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(currencyExchange.USDBRL.Bid)
}

func getDollarExchangeRate(ctx context.Context) (*CurrencyExchange, error) {
	ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bodyRes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var currencyExchange CurrencyExchange
	err = json.Unmarshal(bodyRes, &currencyExchange)
	if err != nil {
		return nil, err
	}
	return &currencyExchange, nil
}

func logContextTimeout(ctxError error, message string) {

	if ctxError != nil && ctxError == context.DeadlineExceeded {
		logger.Println(message)
	}
}

func saveIntoDatabase(ctx context.Context, currency CurrencyDetails) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()
	err := db.WithContext(ctx).Create(&currency).Error

	logContextTimeout(ctx.Err(), "Operação de gravar no banco de dados expirou devido ao timeout de 10ms")

	return err
}
