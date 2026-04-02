package jahitin_be_db

// Store is the application-level database contract.
// It reuses all sqlc-generated query methods from Querier.
type Store interface {
	Querier
}
