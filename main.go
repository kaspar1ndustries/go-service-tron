package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	// "io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	log "go-tron-wallet/logger"
)

const (
	TronMainNet       = "https://api.trongrid.io"
	TronTestNetShasta = "https://api.shasta.trongrid.io"
	TronTestNetNile   = "https://nile.trongrid.io"
)

const API_HOST = TronTestNetNile

type AccountInfo struct {
	Address string `json:"address"`
	Balance int64  `json:"balance"`
}

type Token struct {
	Symbol  string `json:"symbol"`
	Address string `json:"address"`
}

type Tx struct {
	Token     string `json:"token"`
	Amount    string `json:"amount"`
	Direction string `json:"direction"`
}

func apiRequest(url string) ([]byte, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")
	res, _ := http.DefaultClient.Do(req)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	defer res.Body.Close()
	// decoding api response to basic struct
	// check for success and error keys
	resp := struct {
		Success bool   `json:"success"`
		Error   string `json:"error"`
	}{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Error(fmt.Sprintf("api wrapper err: %s", err.Error()))
		return nil, err
	}

	if !resp.Success {
		log.Error(fmt.Sprintf("API request got error: %s", resp.Error))
		return nil, fmt.Errorf(resp.Error)
	}
	log.Info(fmt.Sprintf("api success, returning %d bytes", len(body)))

	return body, err
}

func main() {
	r := gin.Default()

	// Recovery middleware
	r.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {

			log.Error(fmt.Sprintf("panic: %s", err))
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	r.GET("/panic", func(c *gin.Context) {

		panic("test_PANIC")
	})

	r.GET("/address/:address/info", func(c *gin.Context) {
		address := c.Param("address")

		account, err := acctountInfo(address)
		if err != nil {
			log.Error(fmt.Sprintf("api error: %s", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		log.Info(fmt.Sprintf("address info: %+v", account))

		c.JSON(http.StatusOK, gin.H{
			"address": account.Address,
			"balance": account.Balance,
		})

	})

	r.GET("/address/:address/transactions", func(c *gin.Context) {
		address := c.Param("address")
		txs, err := transactionsTrc20(address)
		if err != nil {
			log.Error(fmt.Sprintf("api error: %s", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"address":      address,
			"transactions": txs,
		})
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	serverHostPort := os.Getenv("SERVER_PORT")
	log.Info("Starting server on " + serverHostPort)
	r.Run(serverHostPort)
}

// TODO move to package
// API DOCS https://developers.tron.network/reference/get-account-info-by-address
func acctountInfo(address string) (AccountInfo, error) {

	url := fmt.Sprintf("%s/v1/accounts/%s", API_HOST, address)
	data, err := apiRequest(url)
	if err != nil {
		return AccountInfo{}, err
	}

	resp := struct {
		Data []AccountInfo `json:"data"`
	}{}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return AccountInfo{}, err
	}

	if len(resp.Data) == 0 {
		return AccountInfo{}, errors.New("no data")
	}

	return resp.Data[0], nil
}

// API DOCS https://developers.tron.network/reference/get-trc20-transaction-info-by-account-address
func transactionsTrc20(address string) ([]Tx, error) {

	url := fmt.Sprintf("%s/v1/accounts/%s/transactions/trc20", API_HOST, address)
	data, err := apiRequest(url)
	if err != nil {
		log.Error(fmt.Sprintf("Api request: %s", err))
		return nil, err
	}

	resp := struct {
		Data []struct {
			Txid      string `json:"transaction_id"`
			Token     Token  `json:"token_info"`
			From      string `json:"from"`
			To        string `json:"to"`
			BlockTime int64  `json:"block_timestamp"`
			Type      string `json:"type"`
			Value     string `json:"value"`
		} `json:"data"`
	}{}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		log.Error(fmt.Sprintf("Error unmarshalling data: %s", err))
		return nil, err
	}
	txs := []Tx{}
	for _, tx := range resp.Data {
		var direction string
		if tx.From == address {
			direction = "out"
		} else {
			direction = "in"
		}
		if tx.Token.Symbol != "" {
			txs = append(txs, Tx{
				Token:     tx.Token.Symbol,
				Amount:    tx.Value,
				Direction: direction,
			})
		}
	}
	return txs, nil
}
