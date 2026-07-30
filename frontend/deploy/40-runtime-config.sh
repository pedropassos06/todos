#!/bin/sh
set -eu

API_BASE_URL_VALUE="${API_BASE_URL:-http://localhost:8081}"

cat > /usr/share/nginx/html/config.js <<EOF
window.__APP_CONFIG__ = {
  API_BASE_URL: "${API_BASE_URL_VALUE}"
};
EOF
