-- Drop index idx_products_name
DROP INDEX IF EXISTS idx_products_name;

-- Drop index idx_cart_items_cart_id
DROP INDEX IF EXISTS idx_cart_items_cart_id;

-- Drop index idx_cart_items_product_id
DROP INDEX IF EXISTS idx_cart_items_product_id;

-- Drop unique index idx_cart_items_cart_id_product_id
DROP INDEX IF EXISTS idx_cart_items_cart_id_product_id;

-- Drop index idx_carts_user_id
DROP INDEX IF EXISTS idx_carts_user_id;

-- Drop index idx_carts_status
DROP INDEX IF EXISTS idx_carts_status;