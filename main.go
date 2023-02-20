package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	WalletId string `json:"wallet_id"`
	Wallet   Wallet `json:"wallet"`
}

type Wallet struct {
	Id     string `json:"id"`
	Amount int64  `json:"amount"`
}

type Deposit struct {
	Id       string `json:"id"`
	WalletId string `json:"wallet_id"`
	Amount   int64  `json:"amount"`
}

type Transaction struct {
	Id       string `json:"id"`
	Amount   int64  `json:"amount"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
}

func main() {
	users := []User{}
	transactions := []Transaction{}
	wallets := []Wallet{}
	deposits := []Deposit{}

	engine := gin.Default()
	engine.POST("/user/registration", func(c *gin.Context) {
		id, err := uuid.NewRandom()
		if err != nil {
			fmt.Println(err)
			return
		}
		walletId, err := uuid.NewRandom()
		if err != nil {
			fmt.Println(err)
			return
		}
		uu1 := id.String()
		uu2 := walletId.String()
		var request User
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if request.Username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "request is bad"})
			return
		}
		for _, user := range users {
			if request.Username == user.Username {
				c.JSON(http.StatusBadRequest, gin.H{"error": "user is already exit"})
				return
			}
		}
		request.Id = uu1
		request.WalletId = uu2

		var newWallet Wallet

		newWallet.Id = uu2
		newWallet.Amount = 0

		wallets = append(wallets, newWallet)
		users = append(users, request)

		request.Wallet = newWallet

		c.JSON(200, gin.H{
			"user": request,
		})
	})

	engine.GET("/user/:username", func(c *gin.Context) {
		var user *User
		username := c.Param("username")
		for _, v := range users {
			if v.Username == username {
				user = &v
				break
			}
		}
		for _, v := range wallets {
			if v.Id == user.WalletId {
				user.Wallet = v
				break
			}
		}
		if wallets == nil {
			c.JSON(400, gin.H{
				"message": "Something's wrong",
			})
			return
		}
		if user == nil {
			c.JSON(400, gin.H{
				"message": "user not found",
			})
			return
		}
		c.JSON(200, user)
	})

	engine.GET("/user/wallet/:walletId", func(c *gin.Context) {
		var wallet *Wallet
		walletId := c.Param("walletId")
		for _, v := range wallets {
			if v.Id == walletId {
				wallet = &v
				break
			}
		}
		if wallet == nil {
			c.JSON(400, gin.H{
				"message": "user not found",
			})
			return
		}
		c.JSON(200, wallet)
	})

	engine.POST("/transactions", func(c *gin.Context) {
		var request Transaction
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Println(request)
		for i, v := range wallets {
			if v.Id == request.Sender {
				fmt.Println("I found the account")
				if v.Amount < request.Amount {
					fmt.Println("Out of balance")
					c.JSON(http.StatusBadRequest, gin.H{
						"message": "out of balance",
					})
					return
				}
				v.Amount -= request.Amount
				wallets[i] = v
			}
			if v.Id == request.Receiver {
				v.Amount += request.Amount
				wallets[i] = v
			}
		}
		transactions = append(transactions, request)
		c.JSON(http.StatusOK, gin.H{"is_succeed": true})
	})

	engine.POST("/user/wallet/deposit", func(c *gin.Context) {
		var request Deposit
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		for i, v := range wallets {
			if request.WalletId == v.Id {
				v.Amount += request.Amount
				wallets[i] = v
				break
			}
		}
		deposits = append(deposits, request)
		c.JSON(http.StatusOK, gin.H{"is_succeed": true})
	})

	engine.Run(":3000")
}
