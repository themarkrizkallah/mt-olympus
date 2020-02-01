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
  "user_id": "7aef23ee-3e68-11ea-a0f1-0242ac180003" 
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

### `GET /accounts`
Endpoint to retrieve a list of a user's accounts. 
This endpoint requires a valid `exchange_userCookie` cookie to be passed in with the request.

The response body should looking something like:
```json
[
  {
    "account_id": "99d2f35a-406e-11ea-b1c0-0242ac1a0003",
    "asset_id": "e69fb1a6-406d-11ea-99c1-0242ac1a0003",
    "tick": "USD",
    "balance": 0,
    "holds": 0,
    "created_at": "2020-01-26T19:03:55.938624Z"
  },
  {
    "account_id": "99d379ba-406e-11ea-b1c0-0242ac1a0003",
    "asset_id": "e6a04bc0-406d-11ea-99c1-0242ac1a0003",
    "tick": "BTC",
    "balance": 0,
    "holds": 0,
    "created_at": "2020-01-26T19:03:55.938624Z"
  }
]
```

### `POST/accounts/:account_id/deposit`
Endpoint to deposit funds to user accounts.
This endpoint requires an account id and a valid `exchange_userCookie` cookie to be passed in with the request.
The payload should have the following field:
```json
{
  "amount": 1000
}
```
- `amount`: int64 amount to deposit

The response body should look something like:
```json
{
  "balance": 1000
}
```

### `POST/accounts/:account_id/withdraw`
Endpoint to withdraw funds from user accounts.
This endpoint requires an account id and a valid `exchange_userCookie` cookie to be passed in with the request.
The payload should have the following field:
```json
{
  "amount": 1000
}
```
- `amount`: int64 amount to withdraw, must be <= (account balance - holds on account)

The response body should look something like:
```json
{
  "balance": 0
}
```

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
    "order_id": "4b8a7396-33df-4fdc-b921-4bad3a264f0b",
    "amount": 11,
    "price": 1200,
    "side": 0,
    "type": 0,
    "message": "Confirmed",
    "created_at": "2020-01-24T06:34:32.232515945Z"
  }
}
```