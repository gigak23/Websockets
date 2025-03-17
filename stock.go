package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go/v2"
	"github.com/joho/godotenv"
	"golang.org/x/net/websocket"
)

type server struct {
	conns map[*websocket.Conn]bool
}

func NewStockServer() *server {
	return &server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *server) getStockInfo(ws *websocket.Conn, finnhubClient *finnhub.DefaultApiService) {
	for {
		stockSymbols, _, err := finnhubClient.StockSymbols(context.Background()).Exchange("US").Execute()
		if err != nil {
			log.Println("Could not get stock symbols", err)
		}

		data, err := json.Marshal(&stockSymbols)
		fmt.Printf("%+v\n", stockSymbols[0:5])
		ws.Write(data)
	}

}

func allowWebsocketOrigin(handler websocket.Handler) websocket.Server {
	return websocket.Server{
		Handshake: func(c *websocket.Config, req *http.Request) error {
			return nil
		},
		Handler: handler,
	}
}

func (s *server) test() {
	return
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("Could not load API Key")
	}
	key := os.Getenv("API_KEY")

	cfg := finnhub.NewConfiguration()
	cfg.AddDefaultHeader("X-Finnhub-Token", key)
	finnhubClient := finnhub.NewAPIClient(cfg).DefaultApi

	stockServer := NewStockServer()
	http.Handle("/ws/stocks", allowWebsocketOrigin(func(ws *websocket.Conn) {
		stockServer.getStockInfo(ws, finnhubClient)
	}))

	log.Fatal(http.ListenAndServe(":3000", nil))
}
