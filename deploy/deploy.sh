#!/bin/bash
# Deploy CriBot to Yandex Cloud Functions
# Prerequisites:
#   - yc CLI installed and configured
#   - Go 1.21+ installed
#   - TELEGRAM_BOT_TOKEN and TELEGRAM_CHAT_ID environment variables set

set -e

FUNCTION_NAME="${FUNCTION_NAME:-cribot}"
RUNTIME="golang121"
MEMORY="128m"
TIMEOUT="10s"
ENTRYPOINT="cmd/function.Handler"

echo "=== Building CriBot ==="
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o function ./cmd/function

echo "=== Creating deployment archive ==="
zip -r function.zip function config/tickers.csv

echo "=== Deploying to Yandex Cloud ==="
yc serverless function version create \
    --function-name="${FUNCTION_NAME}" \
    --runtime="${RUNTIME}" \
    --entrypoint="${ENTRYPOINT}" \
    --memory="${MEMORY}" \
    --execution-timeout="${TIMEOUT}" \
    --source-path="function.zip" \
    --environment="TELEGRAM_BOT_TOKEN=${TELEGRAM_BOT_TOKEN},TELEGRAM_CHAT_ID=${TELEGRAM_CHAT_ID}"

echo "=== Cleaning up ==="
rm -f function function.zip

echo "=== Deploy complete ==="
echo "To set up a timer trigger, run:"
echo "  yc serverless trigger create timer --name=cribot-timer --cron-expression='0/5 * * * ? *' --invoke-function-name=${FUNCTION_NAME}"
