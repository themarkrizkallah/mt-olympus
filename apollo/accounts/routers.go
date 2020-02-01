package accounts

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"apollo/database"
	"apollo/redis"
)

const (
	accountIdParam = "account_id"
	uuidLen        = 36
)

func GetUserAccounts(c *gin.Context) {
	var accounts []Account
	const getAccountsSql = `select id, asset_id, balance, holds, created_at from accounts where user_id = $1`

	assets, err := redis.GetAssetList(c)
	if err != nil {
		log.Println("An error occurred retrieving asset list:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred"})
		return
	}
	assetIdMap := make(map[string]string, len(assets))
	for _, asset := range assets {
		assetIdMap[asset.Id] = asset.Tick
	}

	// Get accounts
	db := database.GetDB()
	userId := c.GetString("user_id")
	rows, err := db.QueryContext(c, getAccountsSql, userId)
	if err != nil {
		log.Println("An error occurred getting accounts:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var account Account

		err := rows.Scan(&account.Id, &account.AssetId, &account.Balance, &account.Holds, &account.CreatedAt)
		if err != nil {
			log.Println("An error occurred reading accounts:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred"})
			return
		}

		account.AssetTick = assetIdMap[account.AssetId]
		accounts = append(accounts, account)
	}

	c.JSON(http.StatusOK, accounts)
}

func Deposit(c *gin.Context) {
	var (
		userId         string
		accountId      string
		balance        int64
		depositPayload DepositPayload
	)

	const depositSql = `update accounts
	set balance = balance + $1
	where id = $2 and user_id = $3
returning balance`


	// Basic accountId verification
	accountId = c.Param(accountIdParam)
	if len(accountId) < uuidLen {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "account_id invalid or missing"})
		return
	}

	// Validate payload
	if err := c.BindJSON(&depositPayload); err != nil {
		log.Println("Error unmarshaling deposit payload:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred"})
		return
	} else if depositPayload.Amount <= 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Amount cannot be <= 0"})
		return
	}

	userId = c.GetString("user_id")
	db := database.GetDB()
	err := db.QueryRowContext(c, depositSql, depositPayload.Amount, accountId, userId).Scan(&balance)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid account_id"})
		} else {
			log.Println("Error depositing to account:", err)
			c.JSON(http.StatusForbidden, gin.H{"error": "An error occurred"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": balance})
}

func Withdraw(c *gin.Context) {
	var (
		userId          string
		accountId       string
		balance         int64
		holds           int64
		withdrawPayload WithdrawPayload
	)

	const getAccountSql = `select balance, holds from accounts where id = $1 and user_id = $2`
	const withdrawSql = `update accounts
	set balance = balance - $1
	where id = $2
returning balance`

	userId = c.GetString("user_id")

	// Basic accountId verification
	accountId = c.Param(accountIdParam)
	if 	len(accountId) < uuidLen {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "account_id invalid or missing"})
		return
	}

	// Validate payload
	if err := c.BindJSON(&withdrawPayload); err != nil {
		log.Println("Error unmarshaling deposit payload:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred"})
		return
	} else if withdrawPayload.Amount <= 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Amount cannot be <= 0"})
		return
	}

	// Get current balance and holds
	db := database.GetDB()
	err := db.QueryRowContext(c, getAccountSql, accountId, userId).Scan(&balance, &holds)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid account_id"})
		} else {
			log.Println("Error retrieving account:", err)
			c.JSON(http.StatusForbidden, gin.H{"error": "An error occurred"})
		}
		return
	}

	// Check if user has sufficient funds available funds to withdraw
	if (balance - holds) < withdrawPayload.Amount {
		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient funds"})
		return
	}

	// Execute the withdrawal
	err = db.QueryRowContext(c, withdrawSql, withdrawPayload.Amount, accountId).Scan(&balance)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid account_id"})
		} else {
			log.Println("Error withdrawing from account:", err)
			c.JSON(http.StatusForbidden, gin.H{"error": "An error occurred"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": balance})
}
