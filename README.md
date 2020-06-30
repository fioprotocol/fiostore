# fiostore

Very simple http API for sending FIO requests.

- Listens on port 8080
- Endpoint is `/v1/send_request`
- Options are set via ENV vars:
  - PRIV_KEY wif key for sending requests
  - NODEOS what api to use for sending requests
  - TOKENS comma seperated list of authentication tokens
  - FIO_ADDRESS the sender's FIO address

## example:

```
 $ curl -s localhost:8080/v1/send_request -d '{
     "fio_address": "bp0-east@dapixbp",
     "chain_code": "BTC",
     "token_code": "BTC",
     "amount": 0.001,
     "memo": "invoice for cool t-shirt",
     "access_token": "abc123"
   }' | jq .
   {
     "sent": true,
     "code": 200,
     "message": "success",
     "txid": "b3388eb1d1494be32b32462155a9c97241d73a0935abb00d3ffc6dd7d3bc3afa"
   }
```

*Note: the `amount` field expects a float*

