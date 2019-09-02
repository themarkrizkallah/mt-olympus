# matcher
My initial take on the matching engine, the heart and soul of the exchange.

## Overview
Decided to initially write it in **Go** for quick iteration. Once fully functional, it will be rewritten in **C++**.

The matching engine listens for messages from Kafka on the `orders` topic, processes the order, then sends a message with the completed trade(s) (if applicable) to the `trades` topic.
 
Supports the following operations:
1. **Create** limit orders
2. **Modify** limit orders (**TODO**)
3. **Cancel** limit orders (**TODO**)

It works, but still got a lot of work to do on this one :)