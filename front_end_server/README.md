# front_end_server
This **REST** Go server is where it all starts!

## Overview
This server acts as the entry point to the platform. 

Its main purpose is to handle user signup and authentication as well as to act as a proxy between the **React** front-end (**TODO**) and the gRPC server responsible for all interactions with the matching engine.



## API Specs
By default, the server will be listening on `localhost:8080`

### `POST /orders`
Endpoint to create a new limit order. The payload should have the following fields:

```json
{
	"id": "1",
	"amount": 1,
	"price": 3,
	"side": true,
	"created_at": "2014-03-31T14:11:29+02:00"
}
```

**Field Descriptions**:
- `id`
  - `string` corresponding to the authenticated user's id
  - **Note**: This will be generated on the server side once user signups and authentication have been implemented
- `amount`:
  - `uint64` representing the number of units of the object to buy/sell
- `price`:
  - `uint64` representing the price/unit of the object at hand
  - **TODO**: Determine tick size
- `side`:
  - `bool` corresponding to whether the order is on the buy (`true`) or sell (`false`) side
- `created_at`:
  - `string` representing the `rfc3339` timestamp in which the order was created
  - **Note**: This will soon be generated on the server side for obvious reasons
