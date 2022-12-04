package main

import (
	"errors"
	"fmt"

	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	api "go-tron-wallet/api"
	log "go-tron-wallet/logger"
)

func main() {

	var tg api.TronGrid
	if os.Getenv("NETWORK") == "mainnet" {
		tgApiKey := os.Getenv("API_KEY")
		if tgApiKey == "" {
			log.Fatal("API_KEY env var not set")
		}
		tg = api.NewTronGridMainnet(tgApiKey)
	} else {
		tg = api.NewTronGridTestnet()
	}



	// GIN
	r := gin.Default()
	// Recovery middleware
	r.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			log.Error(errors.New(err))
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	r.GET("/address/:address/info", func(c *gin.Context) {
		address := c.Param("address")
		account, err := tg.Account(address)
		if err != nil {
			log.Error(err)
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
		txs, err := tg.TransactionsTrc20(address)
		if err != nil {
			log.Error(err)
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
		c.JSON(http.StatusOK, "pong")
	})
	serverHostPort := os.Getenv("SERVER_PORT")
	log.Info("Starting server on " + serverHostPort)
	r.Run(serverHostPort)
}
