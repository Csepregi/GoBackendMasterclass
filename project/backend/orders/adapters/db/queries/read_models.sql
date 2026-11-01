-- name: ListMenuItemsWithRestaurant :many
-- Lists menu items with optional restaurant name filter and dynamic ordering.
-- Uses CASE WHEN to support multiple ordering options in a single query.
SELECT
    mi.restaurant_menu_item_uuid AS menu_item_uuid,
    mi.name AS menu_item_name,
    mi.gross_price,
    r.currency,
    r.restaurant_uuid,
    r.name AS restaurant_name
FROM orders.restaurant_menu_items mi
JOIN orders.restaurants r ON mi.restaurant_uuid = r.restaurant_uuid
WHERE mi.is_archived = false
  AND (sqlc.narg(restaurant_name_filter)::text IS NULL
       OR LOWER(r.name) LIKE LOWER('%' || sqlc.narg(restaurant_name_filter)::text || '%'))
ORDER BY
    CASE WHEN (sqlc.narg(order_by)::text IS NULL OR sqlc.narg(order_by)::text = 'default')
         THEN r.name END ASC,
    CASE WHEN (sqlc.narg(order_by)::text IS NULL OR sqlc.narg(order_by)::text = 'default')
         THEN mi.ordering END ASC,
    CASE WHEN sqlc.narg(order_by)::text = 'price_asc' THEN mi.gross_price END ASC,
    CASE WHEN sqlc.narg(order_by)::text = 'price_desc' THEN mi.gross_price END DESC,
    CASE WHEN sqlc.narg(order_by)::text = 'name_asc' THEN mi.name END ASC,
    CASE WHEN sqlc.narg(order_by)::text = 'name_desc' THEN mi.name END DESC;
