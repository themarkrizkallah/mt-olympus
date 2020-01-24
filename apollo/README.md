# apollo
This is your entry point to the exchange.

## API Specs

### `POST /signup`
Endpoint to create a new user. The payload should have the following fields:

```json
{
  "email": "1337@domain.co",
  "password": "badpass"
}
```

The response body should return a UUID on success that corresponds to the user id:

```json
{
  "response": "7aef23ee-3e68-11ea-a0f1-0242ac180003" 
}
```

### `POST /login`
Endpoint to log an user in. The payload should have the following fields:

```json
{
  "email": "1337@domain.co",
  "password": "badpass"
}
```
A cookie (`exchange_userCookie`) should be attached to the response that can be used to authenticate in future requests.


### `POST /orders`
Endpoint to place an order. This endpoint requires a valid `exchange_userCookie` cookie to be passed in with the request.
The payload should have the following fields:

```json
{
  "amount": 10,
  "price":  1000,
  "side": 0,
  "type": 0
}
```
- `side`: an enum representing the side of the order
    - 0 (BUY)
    - 1 (SELL)
- `type`: an enum representing the order type
    - 0 (LIMIT) # Implemented
    - 1 (MARKET) # Todo
    - 2 (STOP) # Todo

The response body should look something like:
```json
{
  "response": {
    "user_id": "7aef23ee-3e68-11ea-a0f1-0242ac180003",
    "order_id": "d932f451-1bf8-45c8-b510-a2b73d05d599",
    "amount": 10,
    "price": 1000,
    "message": "Confirmed",
    "created_at": {
      "seconds": 1579842931,
      "nanos": 163333530
    }
  }
}
```