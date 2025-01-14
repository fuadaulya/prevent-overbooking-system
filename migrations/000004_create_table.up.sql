-- Create table carts
CREATE TABLE carts (
    id SERIAL PRIMARY KEY,                      -- Unique ID for each cart
    user_id INT NOT NULL,                       -- ID of the user owning the cart
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),-- Time the cart was created
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),-- Time the cart was last updated
    status VARCHAR(50) NOT NULL DEFAULT 'active' -- Status of the cart (e.g., active, checked_out)
);

-- Create table cart_items
CREATE TABLE cart_items (
    id SERIAL PRIMARY KEY,                      -- Unique ID for each cart item
    cart_id INT NOT NULL,                       -- Cart ID, references the carts table
    product_id INT NOT NULL,                    -- Product ID, references the products table
    quantity INT NOT NULL CHECK (quantity > 0), -- Quantity of the product, must be greater than 0
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),-- Time the item was added
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),-- Time the item was updated
    UNIQUE (cart_id, product_id),               -- Unique combination of cart_id and product_id
    FOREIGN KEY (cart_id) REFERENCES carts(id) ON DELETE CASCADE,    -- Relation to carts table
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE RESTRICT -- Relation to products table
);