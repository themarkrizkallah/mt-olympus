package database

import (
	"database/sql"

	"matcher/types"
)

const (
	getAccountSql   = `select id, balance, holds from accounts where user_id = $1 and asset_id = $2 limit 1`
	putHoldSql      = `update accounts set holds = holds + $1 where asset_id = $2 and user_id = $3`
	transferHoldSql = `update accounts set balance = balance + $1, holds = holds + $2 where asset_id = $3 and user_id = $4`
	transferSql     = `update accounts set balance = balance + $1 where asset_id = $2 and user_id = $3`
)

func GetAccount(tx *sql.Tx, userId, assetId string) (types.Account, error) {
	var account types.Account

	stmt, err := tx.Prepare(getAccountSql)
	if err != nil {
		return account, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(userId, assetId).Scan(&account.Id, &account.Balance, &account.Holds)

	return account, err
}

func PutHold(tx *sql.Tx, userId, assetId string, amount int64) error {
	stmt, err := tx.Prepare(putHoldSql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(amount, assetId, userId)

	return err
}

// MUST EDIT THIS
func TransferValue(tx *sql.Tx, trade types.Trade, baseId, quoteId string) error {
	var err error = nil

	buyOrder := trade.Buy
	sellOrder := trade.Sell
	tradeMsg := trade.TradeMsg

	price := tradeMsg.GetPrice()
	amount := tradeMsg.GetAmount()

	// Reduce Buyer QUOTE hold and balance
	stmt, err := tx.Prepare(transferHoldSql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err = stmt.Exec(-(amount * price), -(amount * buyOrder.Price), quoteId, buyOrder.UserId); err != nil {
		return err
	}

	// Reduce Seller BASE hold and balance
	stmt, err = tx.Prepare(transferHoldSql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err = stmt.Exec(-amount, -amount, baseId, sellOrder.UserId); err != nil {
		return err
	}

	// Increase Buyer BASE balance
	stmt, err = tx.Prepare(transferSql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err = stmt.Exec(amount, baseId, buyOrder.UserId); err != nil {
		return err
	}

	// Increase Seller QUOTE balance
	stmt, err = tx.Prepare(transferSql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err = stmt.Exec(amount * tradeMsg.Price, quoteId, sellOrder.UserId); err != nil {
		return err
	}

	return err
}
