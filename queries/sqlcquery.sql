-- name: InsertCartItem :exec
INSERT INTO items (item_id, item_name, price, quantity) 
VALUES ($1, $2, $3, $4);

-- name: GetUniqueItemCountInCart :one
SELECT COUNT(DISTINCT item_id) AS total_unique_items
FROM items;

-- name: DeleteCartItem :exec
DELETE FROM items WHERE item_id = $1;

-- name: ClearCartItems :exec
DELETE FROM items;

-- name: CheckItemInCart :one
SELECT EXISTS (
    SELECT 1
    FROM items
    WHERE item_id = $1
) AS item_exists;

-- name: ModifyItemQuantity :exec
UPDATE items
SET quantity = $2
WHERE item_id = $1;

-- name: ApplyPercentageDiscount :exec
UPDATE items
SET price = GREATEST(price - (price * $1 / 100), 0)  
WHERE item_id = $2;

-- name: ApplyFlatDiscount :exec
UPDATE items
SET price = GREATEST(price - $1, 0)
WHERE item_id = $2;


-- name: GetCartItems :many
SELECT item_id, item_name, price, quantity
FROM items;

-- name: CompleteCheckout :exec
DELETE FROM items;

-- name: CheckCartIsEmpty :one
SELECT COUNT(*) FROM items;

-- name: InsertProduct :exec
INSERT INTO products (product_id, product_name, price, stock)
VALUES ($1, $2, $3, $4);

-- name: FetchProductByID :one
SELECT product_id, product_name, price, stock
FROM products
WHERE product_id = $1;

-- name: UpdateProductStockLevel :exec
UPDATE products
SET stock = $2
WHERE product_id = $1;

-- name: GetItemByID :one
SELECT item_id, item_name, price, quantity
FROM items
WHERE item_id = $1;
