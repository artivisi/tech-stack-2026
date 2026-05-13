# Validation rules

Canonical rules. All three stacks must implement identically — comparing stacks only makes sense if they accept and reject the same inputs.

## Fields

| Field (HTTP / form) | DB column | Required | Min | Max | Format (regex) | Normalization |
|---------------------|-----------|----------|-----|-----|----------------|---------------|
| `email` | `email` | yes | 3 | 254 | `^[^\s@]+@[^\s@]+\.[^\s@]+$` | trim, lowercase |
| `fullName` | `full_name` | yes | 2 | 100 | `^[\p{L}\p{M}\s.'\-]+$` (unicode letter, mark, space, `.`, `'`, `-`) | trim |
| `phone` | `phone` | yes | 7 | 20 | `^[+0-9 ()\-]+$` | trim |

HTTP form field names use camelCase (`fullName`); DB column names use snake_case (`full_name`). Each stack maps between the two in its repository layer (Spring `RowMapper`, Express SQL alias, Go `db` struct tag).

Lengths are **character count** (not byte count) after normalization.

## Error semantics

- Invalid input → HTTP **400** with field-level error messages. Form re-renders with the user's **submitted** values (not the normalized values, so the user sees what they typed).
- Duplicate email → HTTP **409** with `email: "email is already registered"`.
- All other DB errors (including CHECK violations, which indicate app/DB drift) → HTTP **500**. CHECK violations must be loud — they signal that validation has drifted.

## Defense in depth

Validation runs at two layers. They are intentionally redundant:

1. **App layer** — enforces format, length, normalization. Returns 400 with friendly messages.
2. **DB layer** — `CHECK` constraints enforce length bounds + email format + email lowercase + phone format. A CHECK violation in production means the app and DB have drifted and is a programming error, not a user error.

The `full_name` unicode regex (`\p{L}\p{M}`) is enforced **only at the app layer** because PostgreSQL's POSIX regex flavor does not match JavaScript/Java/Go unicode property escapes portably. The DB caps full_name length only.

## Stack implementations

| Stack | Tool |
|-------|------|
| Express | `zod` schema (`src/validation.js`) |
| Spring Boot | Bean Validation annotations (`@NotBlank`, `@Size`, `@Pattern`) — planned |
| Go | `go-playground/validator` struct tags — planned |

Each stack must produce identical accept/reject behavior for the same payload.
