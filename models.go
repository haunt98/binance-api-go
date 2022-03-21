package binanceapi

type GetCandlestickRequest struct {
	Symbol      string `url:"symbol"`
	Interval    string `url:"interval"`
	StartTimeMs int64  `url:"startTime,omitempty"`
	EndTimeMs   int64  `url:"endTime,omitempty"`
	Limit       int64  `url:"limit,omitempty"`
}

type GetCandlestickResponse struct {
	Candlesticks []Candlestick
}

type Candlestick struct {
	OpenTimeMs               int64
	Open                     string
	High                     string
	Low                      string
	Close                    string
	Volume                   string
	CloseTimeMs              int64
	QuoteAssetVolume         string
	NumberOfTrades           int64
	TakerBuyBaseAssetVolume  string
	TakerBuyQuoteAssetVolume string
}
