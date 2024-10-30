package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
)

const baseURL = "https://www.alphavantage.co/query"

type TimeSeries struct {
	Open   string `json:"1. open"`
	High   string `json:"2. high"`
	Low    string `json:"3. low"`
	Close  string `json:"4. close"`
	Volume string `json:"5. volume"`
}

type TimeSeriesDailyResponse struct {
	TimeSeries map[string]TimeSeries `json:"Time Series (Daily)"`
}

type Response struct {
	Symbol        string    `json:"symbol"`
	ClosingPrices []float64 `json:"closing_prices"`
	AverageClose  float64   `json:"average_close"`
}

func getClosingPrices(symbol string, numDays int) ([]float64, float64, error) {
	apiKey := os.Getenv("APIKEY")
	url := baseURL + "?function=TIME_SERIES_DAILY&symbol=" + symbol + "&apikey=" + apiKey

	resp, err := http.Get(url)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	var tsResponse TimeSeriesDailyResponse
	if err := json.NewDecoder(resp.Body).Decode(&tsResponse); err != nil {
		return nil, 0, err
	}

	closingPrices := make([]float64, 0)
	total := 0.0
	count := 0

	for _, data := range tsResponse.TimeSeries {
		if count >= numDays {
			break
		}

		closePrice, err := strconv.ParseFloat(data.Close, 64)
		if err == nil {
			closingPrices = append(closingPrices, closePrice)
			total += closePrice
			count++
		}
	}

	averageClose := 0.0
	if count > 0 {
		averageClose = total / float64(count)
	}

	return closingPrices, averageClose, nil
}

func closingPricesHandler(w http.ResponseWriter, r *http.Request) {
	symbol := os.Getenv("SYMBOL")
	numDaysStr := os.Getenv("NDAYS")

	if symbol == "" {
		http.Error(w, "Stock symbol environment variable (SYMBOL) is required.", http.StatusBadRequest)
		return
	}

	numDays, err := strconv.Atoi(numDaysStr)
	if err != nil || numDays <= 0 {
		numDays = 7 // Default to 7 days
	}

	prices, avgPrice, err := getClosingPrices(symbol, numDays)
	if err != nil {
		http.Error(w, "Failed to retrieve stock data.", http.StatusInternalServerError)
		return
	}

	response := Response{
		Symbol:        symbol,
		ClosingPrices: prices,
		AverageClose:  avgPrice,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/api/closing-prices", closingPricesHandler)
	port := ":8080"
	log.Printf("Starting server on port %s...\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
