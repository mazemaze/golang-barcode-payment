package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	WalletId string `json:"wallet_id"`
	Wallet   Wallet `json:"wallet"`
}

type Wallet struct {
	Id        string `json:"id"`
	WalletKey string `json:"wallet_key"`
	Amount    int64  `json:"amount"`
}

type Deposit struct {
	Id       string `json:"id"`
	WalletId string `json:"wallet_id"`
	Amount   int64  `json:"amount"`
}

type Transaction struct {
	Id       string `json:"id"`
	Type     int    `json:"type"`
	Amount   int64  `json:"amount"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
}

type Claim struct {
	Id       string `json:"id"`
	Amount   int64  `json:"amount"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Accepted bool   `json:"accepted"`
}

func main() {
	users := []User{}
	transactions := []Transaction{}
	wallets := []Wallet{}
	deposits := []Deposit{}
	claims := []Claim{}

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
		c.JSON(200, request)
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

	engine.GET("/user/wallet/claim", func(ctx *gin.Context) {
		websocket.Handler(func(ws *websocket.Conn) {
			defer ws.Close()

			err := websocket.Message.Send(ws, "Server: Hello, Client")
			if err != nil {
				log.Println(err)
			}

			for {
				msg := ""
				err := websocket.Message.Receive(ws, &msg)
				if err != nil {
					log.Println(err)
				}

				err = websocket.Message.Send(ws, fmt.Sprintf("Server: \"%s\" received!", msg))
				if err != nil {
					log.Println(err)
				}
			}
		}).ServeHTTP(ctx.Writer, ctx.Request)
	})

	engine.GET("/transactions/claims/:username", func(ctx *gin.Context) {
		var newClaims []Claim
		username := ctx.Param("username")
		for _, v := range claims {
			if v.Id == username {
				newClaims = append(newClaims, v)
			}
		}
		ctx.JSON(http.StatusOK, newClaims)
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

		if request.Type == 1 {
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
		} else if request.Type == 2 {
			isExist := false
			for _, v := range users {
				if v.Id == request.Receiver {
					isExist = true
				}
			}
			if !isExist {
				c.JSON(http.StatusBadRequest, gin.H{"is_succeed": false})
				return
			}
			id, err := uuid.NewRandom()
			if err != nil {
				return
			}
			uu1 := id.String()
			claims = append(claims, Claim{Id: uu1, Amount: request.Amount, Sender: request.Sender, Receiver: request.Receiver, Accepted: false})
		}
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
