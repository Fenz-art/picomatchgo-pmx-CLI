#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."

log() {
  printf '\n===== %s =====\n' "$1"
}

has_python3() {
  command -v python3 >/dev/null 2>&1
}

validate_json_keys() {
  local output="$1"
  shift
  if has_python3; then
    echo "$output" | python3 - "$@" <<'PY'
import json
import sys

try:
    payload = json.load(sys.stdin)
except Exception as exc:
    print(f'invalid JSON: {exc}')
    sys.exit(1)

missing = []
for key in sys.argv[1:]:
    if key not in payload:
        missing.append(key)

if missing:
    print('missing keys:', ', '.join(missing))
    sys.exit(1)
PY
  else
    for key in "$@"; do
      if ! printf '%s' "$output" | grep -q '"'$key'"'; then
        echo "missing key: $key"
        return 1
      fi
    done
  fi
}

log "Preparing output directory"
mkdir -p bin

log "Running gofmt"
fmt_out=$(gofmt -l .)
if [ -n "$fmt_out" ]; then
  echo "ERROR: gofmt found unformatted files:" >&2
  echo "$fmt_out" >&2
  exit 1
fi

git diff --exit-code >/dev/null 2>&1 || true

log "Running go vet"
go vet ./...

log "Running unit tests"
go test ./... -count=1

log "Running focused regression proof"
go test ./... -count=1 -run 'Regression|Compat|PMX'

log "Running race tests"
go test ./... -count=1 -race

log "Building repository packages"
go build ./...

log "Building PMX CLI"
go build -o ./bin/pmx ./cmd/pmx

PMX=./bin/pmx

log "Exercising PMX subcommands"
$PMX help
$PMX match "*.js" "src/app.js"
$PMX scan "**/*.js"
$PMX parse "**/*.js"
$PMX explain "**/*.js" --input "src/app.js"
$PMX validate "**/*.js" --input "src/app.js"
$PMX compat --suite basic
$PMX bench
$PMX fuzz --target FuzzScan --time 15s
$PMX fuzz --target FuzzParse --time 15s
$PMX fuzz --target FuzzIsMatch --time 15s
$PMX doctor
$PMX doctor --json
$PMX doctor --ci

log "Validating JSON contracts"
doctor_json_out=$($PMX doctor --json)
validate_json_keys "$doctor_json_out" version project diagnostics summary

agent_inspect_out=$($PMX agent inspect --json)
validate_json_keys "$agent_inspect_out" version project diagnostics

agent_check_out=$($PMX agent check --json)
validate_json_keys "$agent_check_out" version result diagnostics checks next_actions

ci_out=$($PMX ci --json)
validate_json_keys "$ci_out" result checks

log "Running direct benchmark verification"
go test ./... -run '^$' -bench=. -benchmem -count=1

log "Running failure-path tests"
if $PMX match "[" "src/app.js"; then
  echo "ERROR: invalid pattern unexpectedly passed in match" >&2
  exit 1
fi

if $PMX validate "[" --input "src/app.js"; then
  echo "ERROR: invalid pattern unexpectedly passed in validate" >&2
  exit 1
fi

log "Checking broken fixture diagnostics"
fixture_out=$($PMX doctor --json fixtures/js-ts-fail)
if ! printf '%s' "$fixture_out" | grep -q '"ecosystem": "javascript/typescript"'; then
  echo "ERROR: fixture doctor output missing expected ecosystem" >&2
  exit 1
fi
if ! printf '%s' "$fixture_out" | grep -q 'PMX-TS-001'; then
  echo "ERROR: fixture doctor output missing expected diagnostic PMX-TS-001" >&2
  exit 1
fi

$PMX doctor --ci fixtures/js-ts-fail

log "Hardcore validation complete"

echo "SUCCESS: all hardcore validation gates passed"
