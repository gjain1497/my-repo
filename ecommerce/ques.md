Design a simplified e-commerce platform like Amazon with the following core requirements:
Functional Requirements:

Product Management

Add/update products with details (name, price, description, category, stock quantity)
Search products by name, category, or filters
Maintain product inventory


Inventory Management ⭐

Track stock levels in real-time
Handle low stock alerts
Update inventory when orders are placed/cancelled
Support inventory reservations (when item is in cart/order processing)
Handle concurrent inventory updates (race conditions)


Shopping Cart

Add/remove items to/from cart
Update item quantities
Calculate cart total with tax
Validate cart against inventory before checkout


Order Management

Place orders from cart
Track order status (PENDING, CONFIRMED, SHIPPED, DELIVERED, CANCELLED)
View order history
Handle order cancellations and refunds


Payment Processing ⭐

Support multiple payment methods (Credit Card, Debit Card, UPI, Wallet)
Integrate with different payment gateways (Stripe, PayPal, Razorpay)
Process payments with success/failure callbacks
Handle payment retries and refunds
Maintain payment transaction history


Shipping Integration ⭐

Support multiple shipping providers (FedEx, UPS, DHL, local courier)
Calculate shipping costs based on weight, distance, delivery speed
Generate tracking numbers
Update shipping status
Estimate delivery dates


Order Tracking ⭐

Real-time order status updates
Shipping notifications
Delivery confirmations
Track shipment location


User Management

Register/login users
Manage user profiles and multiple delivery addresses
Save payment methods