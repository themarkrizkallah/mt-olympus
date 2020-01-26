package accounts

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"apollo/database"
)

func GetUserAccounts(c *gin.Context) {
	var accounts []Account
	const getAccountsSql = `select id, asset_id, balance, holds, created_at from accounts where user_id = $1`

	// Get accounts
	db := database.GetDB()
	userId := c.GetString("user_id")
	rows, err := db.QueryContext(c, getAccountsSql, userId)
	if err != nil {
		log.Println("An error occurred getting accounts:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var account Account

		err := rows.Scan(&account.Id, &account.AssetId, &account.Balance, &account.Holds, &account.CreatedAt)
		if  err != nil {
			log.Println("An error occurred reading accounts:", err)
			return
		}

		accounts = append(accounts, account)
	}

	c.JSON(http.StatusOK, accounts)
}