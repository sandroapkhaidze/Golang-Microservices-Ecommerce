-- name: CreateOrder :exec
INSERT INTO orders (
    id, user_id, status, total_amount, correlation_id, created_at, updated_at
) VALUES (
             $1, $2, $3, $4, $5, $6, $7
         );

-- name: CreateOrderItem :exec
INSERT INTO order_items (
    id, order_id, product_id, quantity, price
) VALUES (
             $1, $2, $3, $4, $5
         );

-- name: GetOrderByID :one
SELECT * FROM orders WHERE id = $1;

-- name: GetOrdersByUserID :many
SELECT * FROM orders
WHERE user_id = $1
ORDER BY created_at DESC
    LIMIT $2 OFFSET $3;

-- name: UpdateOrderStatus :exec
UPDATE orders
SET status = $2, updated_at = NOW()
WHERE id = $1;

-- name: GetOrderByCorrelationID :one
SELECT * FROM orders WHERE correlation_id = $1;

-- name: GetOrderItemsByOrderID :many
SELECT * FROM order_items WHERE order_id = $1;