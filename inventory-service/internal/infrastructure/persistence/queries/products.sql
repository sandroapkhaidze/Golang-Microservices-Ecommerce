-- name: CreateProduct :exec
INSERT INTO products (
    id, name, description, price, stock_quantity, reserved_stock, is_active
) VALUES (
             $1, $2, $3, $4, $5, $6, $7
         );

-- name: GetProductByID :one
SELECT id, name, description, price, stock_quantity, reserved_stock,
       is_active, created_at, updated_at
FROM products
WHERE id = $1;

-- name: UpdateProduct :exec
UPDATE products
SET name = $2,
    description = $3,
    price = $4,
    stock_quantity = $5,
    reserved_stock = $6,
    is_active = $7
WHERE id = $1;

-- name: DeleteProduct :exec
UPDATE products
SET is_active = false
WHERE id = $1;

-- name: ListProducts :many
SELECT id, name, description, price, stock_quantity, reserved_stock,
       is_active, created_at, updated_at
FROM products
WHERE is_active = true
ORDER BY created_at DESC
    LIMIT $1 OFFSET $2;

-- name: GetActiveProducts :many
SELECT id, name, description, price, stock_quantity, reserved_stock,
       is_active, created_at, updated_at
FROM products
WHERE is_active = true
ORDER BY name ASC;

-- name: GetProductsByIDs :many
SELECT id, name, description, price, stock_quantity, reserved_stock,
       is_active, created_at, updated_at
FROM products
WHERE id = ANY($1::uuid[]);