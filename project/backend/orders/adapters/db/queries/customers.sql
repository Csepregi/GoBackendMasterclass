-- name: InsertCustomer :exec
INSERT INTO orders.customers (customer_uuid, name, phone_number, address)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetCustomerByUUID :one
SELECT * FROM orders.customers WHERE customer_uuid = $1;