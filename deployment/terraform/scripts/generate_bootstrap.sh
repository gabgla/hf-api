#!/bin/bash
set -euo pipefail

# Generate a dummy Lambda bootstrap zip for initial Terraform deployment
# The real Lambda code gets deployed via SAM/CI pipeline

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
OUTPUT_DIR="${SCRIPT_DIR}/../bootstrap"
BOOTSTRAP_FILE="${OUTPUT_DIR}/bootstrap"
ZIP_FILE="${OUTPUT_DIR}/lambda.zip"

mkdir -p "${OUTPUT_DIR}"

# Create a minimal bootstrap that returns a simple response
cat > "${BOOTSTRAP_FILE}" << 'EOF'
#!/bin/sh
set -e

while true; do
  # Get next invocation
  HEADERS=$(mktemp)
  EVENT=$(curl -sS -LD "$HEADERS" "http://${AWS_LAMBDA_RUNTIME_API}/2018-06-01/runtime/invocation/next")
  REQUEST_ID=$(grep -i Lambda-Runtime-Aws-Request-Id "$HEADERS" | tr -d '\r' | cut -d: -f2 | xargs)

  # Return placeholder response
  RESPONSE='{"statusCode":503,"body":"{\"message\":\"Lambda placeholder - deploy real code\"}"}'
  curl -sS -X POST "http://${AWS_LAMBDA_RUNTIME_API}/2018-06-01/runtime/invocation/${REQUEST_ID}/response" -d "$RESPONSE"

  rm -f "$HEADERS"
done
EOF

chmod +x "${BOOTSTRAP_FILE}"

# Create the zip
(cd "${OUTPUT_DIR}" && zip -j lambda.zip bootstrap)

echo "Created: ${ZIP_FILE}"
