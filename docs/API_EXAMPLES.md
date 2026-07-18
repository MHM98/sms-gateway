# SMS Gateway API Examples

Base URL:

```text
http://127.0.0.1:3000/api/v1/sms-gateway
```

Unsafe routes require `X-Idempotency-Key` with a UUID value. Idempotency is handled by Fiber middleware.

## Add Credit to Wallet

```bash
curl -i -X POST \
  'http://127.0.0.1:3000/api/v1/sms-gateway/wallet/add' \
  -H 'Content-Type: application/json' \
  -H 'X-Idempotency-Key: 11111111-1111-4111-8111-111111111111' \
  -d '{
    "user_id": 12,
    "amount": 100
  }'
```

Success: `204 No Content`

Wallet not found: `404 Not Found`

## Get Wallet Balance

```bash
curl -i \
  'http://127.0.0.1:3000/api/v1/sms-gateway/wallet/12'
```

```json
{
  "user_id": 12,
  "balance": 100
}
```

Wallet not found: `404 Not Found`

## Send SMS

```bash
curl -i -X POST \
  'http://127.0.0.1:3000/api/v1/sms-gateway/message/send' \
  -H 'Content-Type: application/json' \
  -H 'X-Idempotency-Key: 33333333-3333-4333-8333-333333333333' \
  -d '{
    "user_id": 12,
    "recipient": "+989121234567",
    "body": "Your verification code is 123456",
    "service_type": "express"
  }'
```

Success: `204 No Content`

`service_type` must be `normal` or `express`. Each accepted message deducts one wallet credit. Insufficient balance returns `409 Conflict`.

## Get User Message Report

```http
GET /report
Content-Type: application/json
```

| Body field | Required | Validation |
| --- | :---: | --- |
| `user_id` | Yes | Positive integer |
| `from` | Yes | `YYYY-MM-DD`; start date is inclusive |
| `to` | Yes | `YYYY-MM-DD`; must be after `from` and is exclusive |

```bash
curl -i -X GET \
  'http://127.0.0.1:3000/api/v1/sms-gateway/report' \
  -H 'Content-Type: application/json' \
  -d '{
    "user_id": 12,
    "from": "2026-07-17",
    "to": "2026-07-19"
  }'
```

Success: `200 OK`

The example returns messages created on July 17 and July 18. Messages are ordered newest first.

Example response:

```json
{
  "user_id": 12,
  "from": "2026-07-17",
  "to": "2026-07-19",
  "messages": [
    {
      "id": 42,
      "recipient": "+989121234567",
      "body": "Your verification code is 123456",
      "service_type": "express",
      "status": "submitted",
      "created_at": "2026-07-18T09:10:11.123456Z",
      "submission_latency_seconds": 3
    }
  ]
}
```

## Common Errors

| Status | Meaning |
| ---: | --- |
| `400` | Missing idempotency key, invalid request body, invalid user ID, or invalid report date range |
| `404` | Wallet not found |
| `409` | Insufficient wallet balance |
| `500` | Internal server error |
