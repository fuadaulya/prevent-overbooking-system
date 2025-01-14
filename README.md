# prevent-overbooking-system
Stock management system to prevent overbooking with a stock reservation mechanism using Go, PostgreSQL, and Redis.

## Analysis: Causes of Bad Reviews During Flash Sale Events (e.g., 12.12 Event)
### What Happened
During the 12.12 event, our online store experienced an influx of orders due to the significant discounts and promotions offered. Customers successfully added items to their carts, checked out, and paid for their orders. However, after a few days, many customers were informed by our Customer Service team that their orders were canceled due to stock unavailability. This situation led to frustration among customers and resulted in numerous bad reviews for our online store.

### Why It Happened
1. Lack of Real-Time Stock Management:

      During the event, inventory was not updated in real-time during the checkout process. Multiple customers could reserve or purchase the same item simultaneously, leading to overbooking.

2. Concurrent Order Processing:

      High traffic during the event overwhelmed the system, causing delays in processing orders. This allowed orders to exceed the actual available stock.

3. No Stock Reservation Mechanism:

      Items added to the cart were not reserved. As a result, customers could proceed to checkout without ensuring the stock was still available for their orders.

4. Inaccurate Stock Reporting:

      The lack of synchronization between the inventory system and the order processing system led to misreported stock quantities. In some cases, inventory quantities dropped to negative values.

5. Delayed Order Processing:

      Delays in processing orders compounded the issue, as orders were confirmed and processed despite the actual stock being depleted.
### Impact
1. Customer Dissatisfaction:
   
      Customers were disappointed when their orders were canceled, despite having successfully completed the checkout and payment processes.
   
2. Operational Inefficiency:
   
    The Order Processing and Customer Service teams had to spend significant time and effort handling cancellations and complaints.
   
3. Reputation Damage:
   
    Negative reviews hurt the store’s reputation and reduced customer trust, potentially leading to lost revenue in the future.

### Conclusion
The root cause of the bad reviews was the absence of a robust stock management system capable of handling high traffic and concurrent orders during peak sales events. Implementing real-time stock reservation and better synchronization mechanisms would prevent these issues from occurring in the future.
