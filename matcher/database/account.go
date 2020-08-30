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


type MatchMetadata struct {
	Size, Price          int64
	BidPrice             int64
	BidUserId, AskUserId string
}
func TransferValue(tx *sql.Tx, meta MatchMetadata, baseId, quoteId string) error {
	var err error = nil

	// Reduce Buyer QUOTE hold and balance
	stmt, err := tx.Prepare(transferHoldSql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err = stmt.Exec(-(meta.Size * meta.Price), -(meta.Size * meta.Price), quoteId, meta.BidUserId); err != nil {
		return err
	}

	// Reduce Seller BASE hold and balance
	stmt, err = tx.Prepare(transferHoldSql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err = stmt.Exec(-meta.Size, -meta.Size, baseId, meta.AskUserId); err != nil {
		return err
	}

	// Increase Buyer BASE balance
	stmt, err = tx.Prepare(transferSql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err = stmt.Exec(meta.Size, baseId, meta.BidUserId); err != nil {
		return err
	}

	// Increase Seller QUOTE balance
	stmt, err = tx.Prepare(transferSql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err = stmt.Exec(meta.Size * meta.Price, quoteId, meta.AskUserId); err != nil {
		return err
	}

	return err
}
