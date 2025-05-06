package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	r.GET("/ping/:id/:contract-address", func(c *gin.Context) {

		//panggil service

		id := c.Param("id")
		contract_address := c.Param("contract-address")

		url := "https://api.coingecko.com/api/v3/coins/" + id + "/contract/" + contract_address

		req, _ := http.NewRequest("GET", url, nil)

		req.Header.Add("accept", "application/json")

		res, _ := http.DefaultClient.Do(req)

		defer res.Body.Close()
		body, _ := io.ReadAll(res.Body)

		type response struct {
			ID         string          `json:"id"`
			Symbol     string          `json:"symbol"`
			Name       string          `json:"name"`
			Platforms  json.RawMessage `json:"platforms"`
			WebSlug    string          `json:"web_slug"`
			MarketData struct {
				CurrentPrice struct {
					USD float64 `json:"usd"`
				} `json:"current_price"`
			} `json:"market_data"`
		}

		// type response struct {
		// 	ID         string          `json:"id"`
		// 	Symbol     string          `json:"symbol"`
		// 	Name       string          `json:"name"`
		// 	Platforms  json.RawMessage `json:"platforms"`
		// 	WebSlug    string          `json:"web_slug"`
		// 	MarketData json.RawMessage `json:"market_data"`
		// }

		output := &response{}

		err := json.Unmarshal(body, &output)

		if err != nil {
			c.JSON(http.StatusInternalServerError, nil)
			return
		}

		c.JSON(http.StatusOK, output)
	})

	r.GET("/coins/all", func(c *gin.Context) {
		// Endpoint untuk get all coins
		url := "https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&order=market_cap_desc&per_page=100&page=1"

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("accept", "application/json")

		res, _ := http.DefaultClient.Do(req)
		defer res.Body.Close()

		body, _ := io.ReadAll(res.Body)

		type Coin struct {
			ID           string  `json:"id"`
			Symbol       string  `json:"symbol"`
			Name         string  `json:"name"`
			CurrentPrice float64 `json:"current_price"`
			MarketCap    float64 `json:"market_cap"`
			TotalVolume  float64 `json:"total_volume"`
		}

		var coins []Coin
		err := json.Unmarshal(body, &coins)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data"})
			return
		}

		c.JSON(http.StatusOK, coins)
	})

	r.GET("/coins/all/stream", func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Flush()

		type Coin struct {
			ID           string  `json:"id"`
			Symbol       string  `json:"symbol"`
			Name         string  `json:"name"`
			CurrentPrice float64 `json:"current_price"`
			MarketCap    float64 `json:"market_cap"`
			TotalVolume  float64 `json:"total_volume"`
		}

		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		// Kirim data pertama kali saat koneksi dibuka
		sendCoinData := func() {
			url := "https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&order=market_cap_desc&per_page=100&page=1"

			req, _ := http.NewRequest("GET", url, nil)
			req.Header.Add("accept", "application/json")

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Println("Request error:", err)
				return
			}
			defer res.Body.Close()

			body, _ := io.ReadAll(res.Body)

			var coins []Coin
			err = json.Unmarshal(body, &coins)
			if err != nil {
				log.Println("Unmarshal error:", err)
				return
			}

			data, _ := json.Marshal(coins)

			fmt.Fprintf(c.Writer, "data: %s\n\n", data)
			c.Writer.Flush()
		}
		sendCoinData()
		for {
			select {
			case <-c.Request.Context().Done():
				return
			case <-ticker.C:
				sendCoinData()
			}
		}
	})

	r.Run() //

}
