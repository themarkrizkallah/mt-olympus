# zeus
This REST Go server is where it all starts!

## Overview
This server acts as the entry point to the platform. Its main purpose is to handle user signup and authentication as well as acting as as a proxy between the front-end and the gRPC server ([grpc_server](../grpc_server/README.md)) responsible for all interactions with the matching engine.

## API Specs
By default, the server will be listening on [localhost:8080](http://localhost:8080/). This can be modified here: [docker-compose.yml](../docker-compose.yml)

### `POST /signup`
Endpoint to create a new user. The payload should have the following fields:

```json
{
  "user_name": "user",
  "email": "user@domain.com",
  "password": "$tR0Ng P4$$w0rD",
  "first_name": "John",
  "last_name": "Smith",
  "address": "1 Boulevard Street"
}
```

where `user_name` and `email` are unique.

### `POST /login`
Endpoint to log an user in. The payload should have the following fields:

```json
{
  "user_name": "user",
  "email": "user@domain.com",
  "password": "$tR0Ng P4$$w0rD"
}
```
**Note**: You only need to provide one of {`user_name`, `email`}. If both are provided, `user_name` is used to search
MongoDB. Currently, this endpoint does not verify that the  provided`user_name` and `email` are related.

The response body should return a cookie (`exchange_userCookie`) on success that can be used for other endpoints.

### `GET /users`
Endpoint to retrieve a json of all users. This endpoint requires a valid `exchange_userCookie` cookie to be passed in 
with the request.

The response body should look like:
```json
{
  "response": [
    {
      "_id": "5d77c5cba06a043fe7e5063e",
      "user_name": "user",
      "email": "user@domain.com",
      "password": "$argon2id$v=19$m=65536,t=3,p=2$Ity3IxYmiTFWpFnNbY1/BQ$2q8jkC1VrJ9hzEAn6n4waq51E+yGrcCytXaXeojTmrY",
      "first_name": "John",
      "last_name": "Smith",
      "address": "1 Boulevard Street"
    }
  ]
}
```

### `POST /orders`
Endpoint to create a new order. The payload should have the following fields:

```json
{
  "amount": 15,
  "price": 100000,
  "side": 1,
  "type": 0
}
```

**Field Descriptions**:
- `amount`:
  - `uint64` representing the number of units of the object to buy/sell
- `price`:
  - `uint64` representing the price/unit of the object at hand
  - **TODO**: Determine tick size
- `side`:
  - `int` an enum corresponding to whether the order is on the buy (1) or sell (1) side
- `type`:
  - `int` an enum corresponding to whether the order is a limit order (0, default), a market order (1), or a stop order (2)
