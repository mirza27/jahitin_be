-- name: ListOrderItemsByOrderID :many
SELECT
    oi.id,
    oi.order_id,
    oi.category_id,
    oi.service_type_id,
    oi.clothes_for,
    oi.custom_service_name,
    oi.notes,
    oi.price,
    oi.status,
    oi.updated_at,
    oi.created_at,
    oi.finished_at,
    cc.name AS category_name,
    st.name AS service_type_name
FROM order_items oi
LEFT JOIN clothes_categories cc ON cc.id = oi.category_id
LEFT JOIN service_types st ON st.id = oi.service_type_id
WHERE oi.order_id = $1
ORDER BY oi.created_at DESC
LIMIT $2 OFFSET $3;


-- name: ListOrderItemsByMultipleOrderIDs :many
SELECT
    oi.id,
    oi.order_id,
    oi.category_id,
    oi.service_type_id,
    oi.clothes_for,
    oi.custom_service_name,
    oi.notes,
    oi.price,
    oi.status,
    oi.updated_at,
    oi.created_at,
    oi.finished_at,
    cc.name AS category_name,
    st.name AS service_type_name
FROM order_items oi
LEFT JOIN clothes_categories cc ON cc.id = oi.category_id
LEFT JOIN service_types st ON st.id = oi.service_type_id
WHERE oi.order_id = ANY($1::int[]);

-- name: CreateOrderItem :one
INSERT INTO order_items (
    order_id, category_id, service_type_id, clothes_for, custom_service_name, notes, price
) VALUES
    ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;