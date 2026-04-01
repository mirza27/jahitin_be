-- name: CreateBaseUser :one
INSERT INTO users (
    name, device_id, user_type, phone
) VALUES (
    $1, $2, $3, $4
) RETURNING *;


-- name: CreateUser :one
INSERT INTO users (
    name, username, email, password_hash, device_id, user_type, phone
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;


-- name: GetUserByDeviceID :one
SELECT * FROM users
WHERE device_id = $1;


-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;


-- name: UpdateUser :one
UPDATE users
SET name = $2, username = $3, email = $4, password_hash = $5, device_id = $6, user_type = $7, phone = $8
WHERE id = $1
RETURNING *;