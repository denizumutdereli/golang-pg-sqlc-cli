-- name: CreateOrder :one
INSERT INTO orders (
    account_id,
    token_id,
    amount,
    side,
    price,
    position,
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetOrder :one
SELECT * FROM orders
WHERE id = $1 LIMIT 1;

-- name: ListOrders :many
SELECT * FROM orders
ORDER BY id
LIMIT $1
OFFSET $2;