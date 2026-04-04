-- name: ListCustomersByUserID :many
SELECT * FROM customers
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;


-- name: CreateCustomer :one
INSERT INTO customers (
    name, user_id, phone, notes
) VALUES (
    $1, $2, $3, $4
) RETURNING *;


-- name: UpdateCustomerNotes :one
UPDATE customers
SET notes = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;


-- name: GetCustomerByID :one
SELECT * FROM customers
WHERE id = $1;