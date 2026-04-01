package jahitin_be_db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
}

// menyimpan semua fungsi to eksekusi db query dan transactions
type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *SQLStore {

	return &SQLStore{
		Queries: New(db),
		db:      db,
	}
}

// run db transactions
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {

			// return err transactions and rollback combined
			return fmt.Errorf("tx err : %v, rb err : %v", err, rbErr)
		}

		return err
	}

	return tx.Commit()

}
