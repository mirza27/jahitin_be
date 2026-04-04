package jahitin_be_db

import "context"

// Store is the application-level database contract.
// It reuses all sqlc-generated query methods from Querier.
type Store interface {
	Querier
	ExecTx(ctx context.Context, fn func(*Queries) error) error
}
