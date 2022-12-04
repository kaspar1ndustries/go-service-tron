package api

import (
	"encoding/json"
	"errors"
	"fmt"
	log "go-tron-wallet/logger"
	"io/ioutil"
	"net/http"
)

type TronGrid struct {
	Host string
	Key  string
}

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
	Timestamp int64  `json:"timestamp"`
	Direction string `json:"direction"`
}

const (
	MAINNET = "https://api.trongrid.io"
	TESTNET = "https://nile.trongrid.io"
)

// init testnet without api key
func NewTronGridTestnet() TronGrid {
	return TronGrid{
		Host: TESTNET,
	}
}

// api key is required for mainnet
func NewTronGridMainnet(apiKey string) TronGrid {
	return TronGrid{
		Host: MAINNET,
		Key:  apiKey,
	}
}

func (tg TronGrid) Balance(address string) (uint64, error) {
	return 0, nil
}
func (tg TronGrid) BalanceTrc20(address string, tokenName string) (uint64, error) {
	return 0, nil
}

// API DOCS https://developers.tron.network/reference/get-account-info-by-address
func (tg TronGrid) Account(address string) (AccountInfo, error) {

	url := fmt.Sprintf("%s/v1/accounts/%s", tg.Host, address)
	data, err := tg.apiRequest(url)
	if err != nil {
		return AccountInfo{}, err
	}
	log.Info("got API data")

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
func (tg TronGrid) TransactionsTrc20(address string) ([]Tx, error) {

	url := fmt.Sprintf("%s/v1/accounts/%s/transactions/trc20", tg.Host, address)
	data, err := tg.apiRequest(url)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	resp := struct {
		Data []struct {
			Txid     string `json:"transaction_id"`
			Token    Token  `json:"token_info"`
			From     string `json:"from"`
			To       string `json:"to"`
			DateTime int64  `json:"block_timestamp"`
			Type     string `json:"type"`
			Value    string `json:"value"`
		} `json:"data"`
	}{}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	txs := []Tx{}
	for _, tx := range resp.Data {
		if tx.Type != "Transfer" {
			continue
		}
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
				Timestamp: tx.DateTime,
				Direction: direction,
			})
		}
	}
	return txs, nil
}

func (tg TronGrid) apiRequest(url string) ([]byte, error) {
	log.Info(fmt.Sprintf("requesting %s", url))
	req, _ := http.NewRequest("GET", url, nil)
	// req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	// req.Header.Add("TRON-PRO-API-KEY", tg.Key)
	res, _ := http.DefaultClient.Do(req)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer res.Body.Close()
	log.Info("got API response")
	// decoding api response to basic struct
	// check for success and error keys
	resp := struct {
		Success bool   `json:"success"`
		Error   string `json:"error"`
	}{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if !resp.Success {
		// map error to resp.Error
		log.Error(errors.New(resp.Error))
		return nil, fmt.Errorf(resp.Error)
	}
	log.Info(fmt.Sprintf("api success, returning %d bytes", len(body)))

	return body, err
}
