# Internal API Contract

Internal routes are served under the `/internal` prefix. They carry **no JWT authentication** and must be protected at the network layer (Kubernetes NetworkPolicy restricting access to Lambda subnet CIDRs only).

---

## GET /internal/users/by-document

Lookup a user by their CPF/CNPJ document. Used by the `tech-challenge-user-authentication` Lambda to authenticate login requests without querying the database directly.

### Query Parameters

| Parameter  | Type   | Required | Description           |
|------------|--------|----------|-----------------------|
| `document` | string | yes      | CPF (11 digits) or CNPJ (14 digits), digits only |

### Response — 200 OK

```json
{
  "person": {
    "id":           1,
    "name":         "string",
    "email":        "string",
    "contact":      "string | null",
    "document":     "string",
    "is_active":    true,
    "address":      "string | null",
    "address_number": "string | null",
    "neighborhood": "string | null",
    "city":         "string | null",
    "country":      "string | null",
    "zip_code":     "string | null",
    "created_at":   "RFC3339 timestamp",
    "updated_at":   "RFC3339 timestamp"
  },
  "user": {
    "id":        1,
    "password":  "bcrypt hash — SENSITIVE",
    "role":      "customer | attendant | mechanic | administrator",
    "person_id": 1
  }
}
```

### Sensitive fields

| Field           | Notes |
|-----------------|-------|
| `user.password` | bcrypt hash (cost 10). The Lambda calls `bcrypt.CompareHashAndPassword` and must **never** log or forward this value. |

### Lambda usage

1. Call `GET <USERS_API_URL>/internal/users/by-document?document=<cpf>`.
2. If **404** → respond with `401 Unauthorized` (invalid credentials).
3. If **200** and `person.is_active == false` → respond with `401 Unauthorized` (inactive user).
4. If **200** → call `bcrypt.CompareHashAndPassword(user.password, submitted_password)`:
   - Match → generate JWT and return `200`.
   - No match → return `401 Unauthorized`.

### Error responses

| Status | Body                              | Condition               |
|--------|-----------------------------------|-------------------------|
| 400    | `{"error": "document is required"}` | `document` query param missing |
| 404    | `{"error": "user not found"}`     | No active user with that document |
