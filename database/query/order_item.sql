-- name: ListOrderItemsByOrderID :many
SELECT * FROM order_items
WHERE order_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;


-- name: ListOrderItemsByMultipleOrderIDs :many
SELECT * FROM order_items
WHERE order_id = ANY($1::int[]);

-- name: CreateOrderItem :one
INSERT INTO order_items (
    order_id, category_id, service_type_id, clothes_for, custom_service_name, notes, price, status
) VALUES
    ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;