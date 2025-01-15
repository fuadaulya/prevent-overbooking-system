# prevent-overbooking-system

Stock management system to prevent overbooking with a stock reservation mechanism using Go, PostgreSQL, and Redis.

## Case Study:

```
We are members of the engineering team of an online store. When we look at ratings for our online store application, we received the following facts:
   1. Customers were able to add items to their carts, check out, and then pay. After several days, however, many of our customers received calls from our Customer Service department stating that their orders had been canceled due to stock unavailability.
   2. These bad reviews generally come within a week after our 12.12 event, in which we held a large flash sale and set up other major discounts to promote our store.
After checking in with our Customer Service and Order Processing departments, we received the following additional facts:
   1. Our inventory quantities are often misreported, and some items even go as far as having a negative inventory quantity.
   2. The misreported items are those that performed very well on our 12.12 event.
   3. Because of these misreported inventory quantities, the Order Processing department was unable to fulfill a lot of orders, and thus requested help from our Customer Service department to call our customers and notify them that we have had to cancel their orders.
Based on the stated facts above, please do the following things:
   1. Describe what you think happened that caused those bad reviews during our 12.12 event and why it happened. Put this in a section in your README.md file.
   2. Based on your analysis, propose a solution that will prevent the incidents from occurring again. Put this in a section in your README.md file.
   3. Based on your proposed solution, build a Proof of Concept that demonstrates technically how your solution will work.
The technical requirements for your Proof of Concept are as follows:
   1. The PoC must be in the form of an API containing as many endpoints as needed.
      a. The API must use JSON as its message format.
      b. The API must use proper response codes and error messages.
   2. The PoC must be runnable locally.
   3. The PoC must be able to capture and process reasonably detailed information that can be found in an online store’s backend systems.
   4. The PoC must demonstrate its ability to prevent further incidents from occurring.
   5. The PoC must contain at least one functional test that can be run from the command line or share postman collection as a file in your repository, which demonstrates the API’s ability to prevent further incidents.
Aspects of the PoC that we will evaluate also include, but are not limited to:
   1. Database schema and entity design.
   2. API endpoints design.
   3. Logging and error handling.
```

## Analysis: Causes of Bad Reviews During Flash Sale Events (e.g., 12.12 Event)

### What Happened

During the 12.12 event, our online store experienced an influx of orders due to the significant discounts and promotions offered. Customers successfully added items to their carts, checked out, and paid for their orders. However, after a few days, many customers were informed by our Customer Service team that their orders were canceled due to stock unavailability. This situation led to frustration among customers and resulted in numerous bad reviews for our online store.

### Why It Happened

1. **Lack of Real-Time Stock Management:**  
   During the event, inventory was not updated in real-time during the checkout process. Multiple customers could reserve or purchase the same item simultaneously, leading to overbooking.

2. **Concurrent Order Processing:**  
   High traffic during the event overwhelmed the system, causing delays in processing orders. This allowed orders to exceed the actual available stock.

3. **No Stock Reservation Mechanism:**  
   Items added to the cart were not reserved. As a result, customers could proceed to checkout without ensuring the stock was still available for their orders.

4. **Inaccurate Stock Reporting:**  
   The lack of synchronization between the inventory system and the order processing system led to misreported stock quantities. In some cases, inventory quantities dropped to negative values.

5. **Delayed Order Processing:**  
   Delays in processing orders compounded the issue, as orders were confirmed and processed despite the actual stock being depleted.

### Impact

1. **Customer Dissatisfaction:**  
   Customers were disappointed when their orders were canceled, despite having successfully completed the checkout and payment processes.

2. **Operational Inefficiency:**  
   The Order Processing and Customer Service teams had to spend significant time and effort handling cancellations and complaints.

3. **Reputation Damage:**  
   Negative reviews hurt the store’s reputation and reduced customer trust, potentially leading to lost revenue in the future.

### Conclusion

The root cause of the bad reviews was the absence of a robust stock management system capable of handling high traffic and concurrent orders during peak sales events. Implementing real-time stock reservation and better synchronization mechanisms would prevent these issues from occurring in the future.

## Proposed Solution

1. The system will reduce stock when users add products to the cart.
2. Utilizing a locking mechanism to prevent concurrent access and atomic operations to ensure data consistency during the stock reduction process in database transactions.
3. Implementing a reserved stock mechanism to temporarily block stock during the add-to-cart process, thereby preventing overbooking.
4. Implementing Redis to enhance performance and reduce database load when users add items to the cart.
5. Stock will only be permanently reduced if the checkout process is successful.
6. Providing a stock restoration mechanism if users do not complete the checkout process or if an error occurs during checkout.

## Proof of Concept

1.  **Add to Cart API**

- **Function**: Ensures the cart exists, retrieves stock (from Redis or Database), validates stock availability, reserves the stock, updates stock in Redis, and adds the item to the cart.

2. **Checkout API**

- **Function**: Validates and confirms reserved stock, then reduces stock permanently upon successful checkout.

3. **Restore Stock**

- **Function**: Retrieves reserved stock for the cart, updates the `reserved_stock` status to 'cancelled', and returns the reserved stock back to the available inventory.

## API Endpoints

1. **Add to Cart API**

- **Endpoint Name**: Add to Cart API
- **HTTP Method**: POST
- **URL Path**: `/cart/:cartID/item`
- **Request Body**:

```bash
curl --location 'http://localhost:8080/cart/12/item' \
--header 'Content-Type: application/json' \
--header 'X-User-ID: 123' \
--data '{
"product_id": 7,
"quantity": 1
}'
```

- **Response Body**:

```bash
Item successfully added to cart
```

2. **Checkout API**

- **Endpoint Name**: Checkout API
- **HTTP Method**: POST
- **URL Path**: `/order/checkout`
- **Request Body**:

```bash
curl --location 'http://localhost:8080/order/checkout' \
--header 'Content-Type: application/json' \
--data '{
  "cart_id": 3
}'
```

- **Response Body**:

```bash
Checkout successful
```

3. **HTTP response status codes**

- `200 OK`
- `400 Bad Request`
- `401 Unauthorized`
- `500 Internal Server Error`

## Database Schema

1. **`products`**

- Stores product information and stock.

```bash
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    price INT NOT NULL,
    stock INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

2. **`carts`**

- Stores user shopping cart information.

```bash
CREATE TABLE carts (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    status VARCHAR(50) NOT NULL DEFAULT 'active'
);
```

3. **`cart_items`**

- Stores the list of items in a shopping cart.

```bash
CREATE TABLE cart_items (
    id SERIAL PRIMARY KEY,
    cart_id INT NOT NULL,
    product_id INT NOT NULL,
    quantity INT NOT NULL CHECK (quantity > 0),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (cart_id, product_id),
    FOREIGN KEY (cart_id) REFERENCES carts(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE RESTRICT
);
```

4. **`reserved_stock`**

- Stores the list of items in a shopping cart.

```bash
CREATE TABLE reserved_stock (
    id SERIAL PRIMARY KEY,
    product_id INT NOT NULL,
    quantity INT NOT NULL,
    status VARCHAR(20) DEFAULT 'reserved', -- reserved, confirmed, canceled, etc.
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (product_id) REFERENCES products(id)
);
```

## Setup and Usage

1. **Clone repository**

```bash
git clone <repo-url>
cd <repo-folder>
```

2. **Setup Environment Variables**

```bash
# Postgres environment
DBUSER=postgres
DBPASSWORD=postgrespassword
DBNAME=mydb
DBHOST=127.0.0.1
PORT=5432
SSL=disable

# Redis environment
REDIS_HOST=localhost  # For local environment
REDIS_PORT=6379       # Default Redis port
```

3. **Run Docker Compose**

```bash
docker-compose up --build
```

4. **Run `golang-migrate`**

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
migrate -path ./migrations -database postgres://myuser:mypassword@localhost:5432/url-shortener?sslmode=disable up
```

## Limitations

1. **_Limited Scope_**: The system only handles processes from add-to-cart to checkout, assuming that the checkout is successful.
2. **_Concurrency Management_**: Although locking mechanisms and atomic operations are implemented, the system may face performance challenges if the number of concurrent transactions is very high.

## Future Work

1. **_Cron Job for Stock Restoration_**: Adding a cron job feature to ensure reserved stock that is not used within a specific timeframe is automatically restored to the inventory.
2. **_Payment Integration_**: Providing integration with a payment system to ensure the checkout transaction is fully completed.
3. **_Real-Time Stock Synchronization_**: Implementing real-time stock synchronization to improve user experience.
4. **_User Behavior_**: Storing information about who made the reservation and implementing stock limitation rules.
