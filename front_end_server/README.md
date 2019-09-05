# front_end_server
This REST Go server is where it all starts!

## Overview
This server acts as the entry point to the platform. Its main purpose is to handle user signup and authentication as well as acting as as a proxy between the front-end and the gRPC server ([grpc_server](../grpc_server/README.md)) responsible for all interactions with the matching engine.

## API Specs
By default, the server will be listening on [localhost:8080](http://localhost:8080/). This can be modified here: [docker-compose.yml](../docker-compose.yml)

### `POST /orders`
Endpoint to create a new limit order. The payload should have the following fields:

```json
{
	"user_id": "1",
	"order_id": "1",
	"amount": 2,
	"price": 100,
	"side": 1,
	"type": 0
}
```

**Field Descriptions**:
- `user_id`
  - `string` corresponding to the authenticated user's id
  - **Note**: This will be generated on the server side once user signups and authentication have been implemented
- `order_id`
  - `string` corresponding to the order id of the authenticated user
  - **Note**: This will be generated on the server side once user signups and authentication have been implemented
- `amount`:
  - `uint64` representing the number of units of the object to buy/sell
- `price`:
  - `uint64` representing the price/unit of the object at hand
  - **TODO**: Determine tick size
- `side`:
  - `int` an enum corresponding to whether the order is on the buy (1) or sell (1) side
- `type`:
  - `int` an enum corresponding to whether the order is a limit order (0, default), a market order (1), or a stop order (2)
