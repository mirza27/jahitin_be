-- name: ListCustomersByUserID :many
SELECT * FROM customers
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;


-- name: ListCustomersByUserIDAndName :many
SELECT * FROM customers
WHERE user_id = $1 AND name ILIKE '%' || $2 || '%'
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;


-- name: CreateCustomer :one
INSERT INTO customers (
    name, user_id, phone, notes, contact_key, country_code, formatted_phone
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetCustomerByUserIDAndPhone :one
SELECT * FROM customers
WHERE user_id = $1 AND formatted_phone = $2
LIMIT 1;

-- name: UpdateCustomerNotes :one
UPDATE customers
SET notes = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;


-- name: GetCustomerByID :one
SELECT * FROM customers
WHERE id = $1;