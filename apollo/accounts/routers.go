package accounts

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"apollo/database"
	"apollo/redis"
)

func GetUserAccounts(c *gin.Context) {
	var accounts []Account
	const getAccountsSql = `select id, asset_id, balance, holds, created_at from accounts where user_id = $1`

	assets, err := redis.GetAssetList(c)
	if  err != nil {
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
		if  err != nil {
			log.Println("An error occurred reading accounts:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred"})
			return
		}

		account.AssetTick = assetIdMap[account.AssetId]
		accounts = append(accounts, account)
	}

	c.JSON(http.StatusOK, accounts)
}