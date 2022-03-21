package main

import (
	"context"
	"log"
	"net/http"
	"time"

	binanceapi "github.com/haunt98/binance-api-go"
)

func main() {
	httpClient := &http.Client{
		Timeout: time.Second * 5,
	}

	service := binanceapi.NewService(httpClient)

	rsp, err := service.GetCandlestick(context.Background(), binanceapi.GetCandlestickRequest{
		Symbol:   "BNBUSDT",
		Interval: "5m",
		Limit:    2,
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("%+v", rsp)
}
