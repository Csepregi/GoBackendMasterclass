-- name: AddCourier :exec
INSERT INTO orders.couriers (
    courier_uuid,
    name,
    phone_number,
    city
)
VALUES ($1, $2, $3, $4);

-- name: GetCourierByUUID :one
SELECT *
FROM orders.couriers
WHERE courier_uuid = $1;

-- name: ListCouriers :many
SELECT *
FROM orders.couriers
ORDER BY name DESC;


-- name: UpdateCourier :exec
UPDATE orders.couriers
SET name = $2, phone_number = $3, city = $4
WHERE courier_uuid = $1;

-- name: DeleteCourier :exec
DELETE FROM orders.couriers
WHERE courier_uuid = $1;