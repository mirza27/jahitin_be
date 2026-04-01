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