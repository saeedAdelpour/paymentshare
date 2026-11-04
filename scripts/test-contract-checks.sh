#!/bin/sh
# Exercise the actual validation and generation commands in a throwaway checkout.
set -eu
root=$(pwd)
temporary=$(mktemp -d)
trap 'rm -rf "$temporary"' EXIT HUP INT TERM
mkdir -p "$temporary/api" "$temporary/internal/httpapi" "$temporary/bin"
cp api/openapi.yaml api/oapi-codegen.yaml "$temporary/api/"
cp internal/httpapi/generated.go "$temporary/internal/httpapi/"
ln -s "$root/bin/oapi-codegen" "$temporary/bin/oapi-codegen"

go run ./cmd/contract
cat > "$temporary/invalid.yaml" <<'SPEC'
openapi: 3.0.3
info: {title: Invalid fixture, version: '1'}
paths: {}
components:
  schemas:
    Invalid:
      type: impossible
SPEC
if go run ./cmd/contract -file "$temporary/invalid.yaml" >"$temporary/validation.log" 2>&1; then
    printf '%s\n' 'Contract validation accepted an invalid schema' >&2
    exit 1
fi
case "$(cat "$temporary/validation.log")" in
    *'OpenAPI validation failed:'*) ;;
    *) cat "$temporary/validation.log" >&2; exit 1 ;;
esac
make --no-print-directory -C "$temporary" -f "$root/Makefile" -o bin/oapi-codegen check-generated
printf '\n// Deliberately stale fixture.\n' >> "$temporary/internal/httpapi/generated.go"
if make --no-print-directory -C "$temporary" -f "$root/Makefile" -o bin/oapi-codegen check-generated >"$temporary/generation.log" 2>&1; then
    printf '%s\n' 'Generation check accepted stale code' >&2
    exit 1
fi
printf '%s\n' 'Verified rejection of an invalid contract and stale generated code'
