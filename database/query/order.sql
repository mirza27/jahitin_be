
-- name: ListUserOrdersFiltered :many
SELECT o.*, c.name AS customer_name
FROM orders o
JOIN customers c ON c.id = o.customer_id
WHERE o.user_id = $1
  AND (NULLIF($2, '')::text IS NULL OR o.status = $2)
  AND (
    NULLIF($3, '')::text IS NULL OR
    o.name ILIKE '%' || $3 || '%' OR
    c.name ILIKE '%' || $3 || '%'
  )
ORDER BY o.created_at DESC
LIMIT $4 OFFSET $5;


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
