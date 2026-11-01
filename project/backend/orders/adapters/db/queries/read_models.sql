-- name: ListMenuItemsWithRestaurant :many
SELECT
    r.name AS restaurant_name,
    mi.restaurant_menu_item_uuid AS menu_item_uuid,
    mi.name AS menu_item_name,
    r.currency,
    r.restaurant_uuid, 
    mi.gross_price
FROM
    orders.restaurant_menu_items mi
JOIN
    orders.restaurants r ON mi.restaurant_uuid = r.restaurant_uuid
WHERE
    mi.is_archived = false
ORDER BY
    r.name,
    mi.ordering;