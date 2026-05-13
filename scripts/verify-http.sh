#!/usr/bin/env bash
# HTTP-level acceptance test for the registration app.
# Stack-agnostic — runs against whichever stack (Spring Boot / Express / Go)
# is currently listening on BASE_URL.
#
# Asserts status codes only (302 / 400 / 409 / 200 / 503).
# Status codes are the most portable cross-stack contract; the exact error
# message text is defined in db/validation.md but not asserted here.
#
# Usage:
#   ./scripts/verify-http.sh                          # default http://localhost:8080
#   BASE_URL=http://host:port ./scripts/verify-http.sh

set -u
BASE_URL="${BASE_URL:-http://localhost:8080}"
TMP=$(mktemp)
trap 'rm -f "$TMP"' EXIT

PASS=0
FAIL=0
TOTAL=0

assert_status() {
  local label="$1" method="$2" path="$3" data="$4" expected="$5"
  TOTAL=$((TOTAL + 1))
  local actual
  if [ "$method" = "GET" ]; then
    actual=$(curl -s -o "$TMP" -w "%{http_code}" "$BASE_URL$path")
  else
    actual=$(curl -s -o "$TMP" -w "%{http_code}" -X "$method" -d "$data" "$BASE_URL$path")
  fi
  if [ "$actual" = "$expected" ]; then
    printf "  PASS  %-44s -> %s\n" "$label" "$actual"
    PASS=$((PASS + 1))
  else
    printf "  FAIL  %-44s -> expected %s got %s\n" "$label" "$expected" "$actual"
    FAIL=$((FAIL + 1))
  fi
}

echo "verify-http.sh against $BASE_URL"
echo "---"

# Unique per-run namespace so re-runs without reset-db.sh don't clash
# until the duplicate-email step (which intentionally re-uses the namespace).
TS=$(date +%s)

echo "[smoke]"
assert_status "GET /"                "GET"  "/"                ""  "200"
assert_status "GET /css/app.css"     "GET"  "/css/app.css"     ""  "200"
assert_status "GET /health"          "GET"  "/health"          ""  "200"
assert_status "GET /registrations"   "GET"  "/registrations"   ""  "200"

echo
echo "[valid registrations -> 302]"
assert_status "plain ASCII"          "POST" "/register" \
  "email=alice${TS}@example.com&full_name=Alice%20Doe&phone=08123456789" "302"
assert_status "uppercase email (normalized)" "POST" "/register" \
  "email=BOB${TS}@example.com&full_name=Bob%20Smith&phone=08123456790" "302"
assert_status "unicode name (Jose)"  "POST" "/register" \
  "email=jose${TS}@example.com&full_name=Jos%C3%A9%20Rodr%C3%ADguez&phone=08123456791" "302"
assert_status "apostrophe (O'Brien)" "POST" "/register" \
  "email=obrien${TS}@example.com&full_name=Sean%20O%27Brien&phone=08123456792" "302"
assert_status "hyphen (Mary-Jane)"   "POST" "/register" \
  "email=mj${TS}@example.com&full_name=Mary-Jane%20Watson&phone=08123456793" "302"
assert_status "phone with parens"    "POST" "/register" \
  "email=parens${TS}@example.com&full_name=Test%20User&phone=%2B62%20%28811%29%20123-4567" "302"

echo
echo "[invalid email -> 400]"
assert_status "no @"          "POST" "/register" "email=invalid&full_name=Test%20User&phone=08123456789" "400"
assert_status "no dot"        "POST" "/register" "email=a@b&full_name=Test%20User&phone=08123456789" "400"
assert_status "space inside"  "POST" "/register" "email=a%20b@c.com&full_name=Test%20User&phone=08123456789" "400"
assert_status "empty"         "POST" "/register" "email=&full_name=Test%20User&phone=08123456789" "400"

echo
echo "[invalid full_name -> 400]"
assert_status "one char"      "POST" "/register" "email=n1${TS}@example.com&full_name=A&phone=08123456789" "400"
assert_status "empty"         "POST" "/register" "email=n2${TS}@example.com&full_name=&phone=08123456789" "400"
assert_status "digits"        "POST" "/register" "email=n3${TS}@example.com&full_name=Alice%203&phone=08123456789" "400"
assert_status "@ in name"     "POST" "/register" "email=n4${TS}@example.com&full_name=Alice%40Doe&phone=08123456789" "400"

echo
echo "[invalid phone -> 400]"
assert_status "too short"     "POST" "/register" "email=p1${TS}@example.com&full_name=Test%20User&phone=12345" "400"
assert_status "letters"       "POST" "/register" "email=p2${TS}@example.com&full_name=Test%20User&phone=abc12345" "400"
assert_status "empty"         "POST" "/register" "email=p3${TS}@example.com&full_name=Test%20User&phone=" "400"

echo
echo "[duplicate email -> 409]"
assert_status "exact dup"     "POST" "/register" \
  "email=alice${TS}@example.com&full_name=Someone%20Else&phone=08111111111" "409"
assert_status "case dup"      "POST" "/register" \
  "email=ALICE${TS}@example.com&full_name=Someone%20Else&phone=08111111112" "409"

echo
echo "---"
echo "TOTAL: $TOTAL  PASS: $PASS  FAIL: $FAIL"
[ "$FAIL" = "0" ]
