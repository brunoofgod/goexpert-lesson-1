package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var logger = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)

	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if ctx != nil && ctx.Err() == context.DeadlineExceeded {
		logger.Println("Operação cancelada por timeout. Maximo aceito é 300ms")
		panic(ctx.Err())
	}

	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	err = saveExchangeIntoFile(res.Body)
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, res.Body)
}

func saveExchangeIntoFile(body io.Reader) error {

	f, err := os.OpenFile("cotacao.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		panic(err)
	}

	var result string
	err = json.NewDecoder(body).Decode(&result)
	if err != nil {
		return fmt.Errorf("Falha ao decodificar resposta: %w", err)
	}

	_, err = f.WriteString("Dólar: " + result + "\n")

	if err != nil {
		return fmt.Errorf("Falha ao salvar o arquivo: %w", err)
	}

	f.Close()
	return nil
}
