-- name: ListOrdersHeaderByUserID :many
SELECT o.*, oi.*, c.name AS customer_name FROM orders o
JOIN order_items oi ON o.id = oi.order_id 
JOIN customers c ON o.customer_id = c.id
WHERE o.user_id = $1
ORDER BY o.created_at DESC
LIMIT $2 OFFSET $3;


-- name: ListOrdersByUserId :many
SELECT * FROM orders
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;


-- name: GetOrderDetailsByOrderID :one
SELECT o.id, o.name, o.deadline, o.status, c.id, c.name AS customer_name, c.phone AS customer_phone FROM orders o
JOIN customers c ON o.customer_id = c.id
WHERE o.id = $1;


-- name: CreateOrder :one
INSERT INTO orders (
    user_id, name, customer_id, deadline, status
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;
