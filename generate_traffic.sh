#!/usr/bin/env bash

# --------------------------------------------------------
# Generates sample HTTP traffic against the Go Ecommerce API
# to produce trace data for OpenTelemetry and SigNoz.
#
# This script hits endpoints like products, orders, inventory,
# CPU test, and concurrency test to simulate real usage patterns.
# --------------------------------------------------------

API_URL="http://localhost:8080"
ITERATIONS=50
SLEEP_INTERVAL=0.5

echo "🚀 Starting telemetry traffic generation..."
echo "Target API: $API_URL"
echo "Iterations: $ITERATIONS"
echo ""

for ((i=1; i<=ITERATIONS; i++)); do
  echo ">>> Request cycle $i"

  # Fetch all products
  curl -s "$API_URL/products" > /dev/null

  # Create a new order
  curl -s -X POST "$API_URL/orders" \
    -H "Content-Type: application/json" \
    -d "{\"product_name\":\"Book\",\"quantity\":$((i % 5 + 1)),\"user_id\":\"user$i\"}" > /dev/null

  # Check inventory
  curl -s "$API_URL/checkInventory" > /dev/null

  # Simulate CPU load
  curl -s "$API_URL/cpuTest" > /dev/null

  # Simulate concurrency
  curl -s "$API_URL/concurrencyTest" > /dev/null

  echo "✔ Completed request cycle $i"
  sleep "$SLEEP_INTERVAL"
done

echo ""
echo "🎉 Completed all $ITERATIONS traffic iterations."
