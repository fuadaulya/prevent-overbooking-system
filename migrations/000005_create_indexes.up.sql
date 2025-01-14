-- Create index on products name
CREATE INDEX idx_products_name ON products (name);

-- Create index on cart_items cart_id
CREATE INDEX idx_cart_items_cart_id ON cart_items (cart_id);

-- Create index on cart_items product_id
CREATE INDEX idx_cart_items_product_id ON cart_items (product_id);

-- Create unique index on cart_items cart_id and product_id
CREATE UNIQUE INDEX idx_cart_items_cart_id_product_id ON cart_items (cart_id, product_id);

-- Create index on carts user_id
CREATE INDEX idx_carts_user_id ON carts (user_id);

-- Create index on carts status
CREATE INDEX idx_carts_status ON carts (status);