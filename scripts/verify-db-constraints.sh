#!/usr/bin/env bash
# Verify DB-level CHECK constraints by bypassing the app and running
# INSERTs directly. Each insert is expected to fail with a named
# CHECK constraint. Confirms defense-in-depth holds.
#
# Stack-agnostic — only touches the database, no app required.
#
# Usage:
#   ./scripts/verify-db-constraints.sh
#   CONTAINER=other-name ./scripts/verify-db-constraints.sh

set -u
CONTAINER="${CONTAINER:-tech-stack-2026-db}"
PASS=0
FAIL=0

assert_violates() {
  local label="$1" expected_constraint="$2" sql="$3"
  local output
  output=$(docker exec "$CONTAINER" psql -U registration -d registration -c "$sql" 2>&1 || true)
  if echo "$output" | grep -q "violates check constraint \"$expected_constraint\""; then
    printf "  PASS  %-30s blocked by %s\n" "$label" "$expected_constraint"
    PASS=$((PASS + 1))
  else
    printf "  FAIL  %-30s did not hit %s\n" "$label" "$expected_constraint"
    echo "        output: $output"
    FAIL=$((FAIL + 1))
  fi
}

echo "verify-db-constraints.sh against container $CONTAINER"
echo "---"

assert_violates "uppercase email" "email_lower" \
  "INSERT INTO registration (id, email, full_name, phone, created_at) VALUES (gen_random_uuid(), 'BAD@EXAMPLE.COM', 'Test User', '08123456789', NOW())"

assert_violates "no @ in email" "email_format" \
  "INSERT INTO registration (id, email, full_name, phone, created_at) VALUES (gen_random_uuid(), 'noatsign', 'Test User', '08123456789', NOW())"

assert_violates "too-long email" "email_length" \
  "INSERT INTO registration (id, email, full_name, phone, created_at) VALUES (gen_random_uuid(), repeat('a', 245) || '@example.com', 'Test User', '08123456789', NOW())"

assert_violates "too-short full_name" "full_name_length" \
  "INSERT INTO registration (id, email, full_name, phone, created_at) VALUES (gen_random_uuid(), 'len@example.com', 'X', '08123456789', NOW())"

assert_violates "too-short phone" "phone_length" \
  "INSERT INTO registration (id, email, full_name, phone, created_at) VALUES (gen_random_uuid(), 'short@example.com', 'Test User', '123', NOW())"

assert_violates "letters in phone" "phone_format" \
  "INSERT INTO registration (id, email, full_name, phone, created_at) VALUES (gen_random_uuid(), 'letters@example.com', 'Test User', 'abcdefgh', NOW())"

echo
echo "---"
echo "PASS: $PASS  FAIL: $FAIL"
[ "$FAIL" = "0" ]
