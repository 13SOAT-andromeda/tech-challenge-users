#!/usr/bin/env bash
# Migrates data from garagedb (tech-challenge-s1) to usersdb (tech-challenge-users).
# Run ONCE after usersdb schema is applied and before cutting over traffic.
#
# Usage:
#   GARAGEDB_URL=postgres://user:pass@host/garagedb \
#   USERSDB_URL=postgres://user:pass@host/usersdb \
#   ./scripts/migrate_garagedb_to_usersdb.sh

set -euo pipefail

: "${GARAGEDB_URL:?GARAGEDB_URL is required}"
: "${USERSDB_URL:?USERSDB_URL is required}"

TABLES=(Person "User" Employee Customer Vehicle Customer_Vehicle Company)

echo "==> Dumping data from garagedb..."
DUMP_FILE=$(mktemp /tmp/garagedb_dump_XXXXXX.sql)
trap 'rm -f "$DUMP_FILE"' EXIT

TABLE_ARGS=()
for t in "${TABLES[@]}"; do
  TABLE_ARGS+=("-t" "\"$t\"")
done

pg_dump \
  --column-inserts \
  --data-only \
  "${TABLE_ARGS[@]}" \
  "$GARAGEDB_URL" \
  > "$DUMP_FILE"

echo "==> Importing into usersdb..."
psql "$USERSDB_URL" < "$DUMP_FILE"

echo "==> Resetting sequences..."
for t in "${TABLES[@]}"; do
  psql "$USERSDB_URL" -c "SELECT setval('\"${t}_id_seq\"', COALESCE((SELECT MAX(id) FROM \"${t}\"), 1));"
done

echo "==> Migration complete."
