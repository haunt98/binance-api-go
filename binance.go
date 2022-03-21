package binanceapi

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/buger/jsonparser"
	"github.com/google/go-querystring/query"
)

const (
	// https://binance-docs.github.io/apidocs/spot/en/#general-info
	baseURL           = "https://api.binance.com"
	getCandlestickURL = "/api/v3/klines"

	defaultLimit = 500
)

type Service interface {
	GetCandlestick(ctx context.Context, req GetCandlestickRequest) (GetCandlestickResponse, error)
}

type service struct {
	httpClient *http.Client
}

func NewService(httpClient *http.Client) Service {
	return &service{
		httpClient: httpClient,
	}
}

// https://binance-docs.github.io/apidocs/spot/en/#kline-candlestick-data
func (s *service) GetCandlestick(ctx context.Context, req GetCandlestickRequest) (GetCandlestickResponse, error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+getCandlestickURL, nil)
	if err != nil {
		return GetCandlestickResponse{}, fmt.Errorf("failed to new http request: %w", err)
	}

	urlValues, err := query.Values(req)
	if err != nil {
		return GetCandlestickResponse{}, fmt.Errorf("failed to parse query values: %w", err)
	}
	httpReq.URL.RawQuery = urlValues.Encode()

	httpRsp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return GetCandlestickResponse{}, fmt.Errorf("failed to do http request: %w", err)
	}
	defer httpRsp.Body.Close()

	if httpRsp.StatusCode != http.StatusOK {
		return GetCandlestickResponse{}, fmt.Errorf("http request status not OK: %d", httpRsp.StatusCode)
	}

	data, err := io.ReadAll(httpRsp.Body)
	if err != nil {
		return GetCandlestickResponse{}, fmt.Errorf("failed to read all http response body: %w", err)
	}

	limit := req.Limit
	if limit == 0 {
		limit = defaultLimit
	}
	candlestickes := make([]Candlestick, 0, limit)

	if _, err := jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		candlestick := Candlestick{}

		candlestick.OpenTime, _ = jsonparser.GetInt(value, "[0]")
		candlestick.Open, _ = jsonparser.GetString(value, "[1]")
		candlestick.High, _ = jsonparser.GetString(value, "[2]")
		candlestick.Low, _ = jsonparser.GetString(value, "[3]")
		candlestick.Close, _ = jsonparser.GetString(value, "[4]")
		candlestick.Volume, _ = jsonparser.GetString(value, "[5]")
		candlestick.CloseTime, _ = jsonparser.GetInt(value, "[6]")
		candlestick.QuoteAssetVolume, _ = jsonparser.GetString(value, "[7]")
		candlestick.NumberOfTrades, _ = jsonparser.GetInt(value, "[8]")
		candlestick.TakerBuyBaseAssetVolume, _ = jsonparser.GetString(value, "[9]")
		candlestick.TakerBuyQuoteAssetVolume, _ = jsonparser.GetString(value, "[10]")

		candlestickes = append(candlestickes, candlestick)
	}); err != nil {
		return GetCandlestickResponse{}, fmt.Errorf("failed to parse json: %w", err)
	}

	return GetCandlestickResponse{
		Candlesticks: candlestickes,
	}, nil
}
