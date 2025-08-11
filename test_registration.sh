#!/bin/bash

echo "Testing Blog API Registration Endpoint"
echo "======================================"

# Check if server is running
echo "1. Checking if server is running..."
if curl -s http://localhost:8080/health > /dev/null; then
    echo "✅ Server is running"
else
    echo "❌ Server is not running. Please start the server first."
    exit 1
fi

echo ""
echo "2. Testing registration endpoint..."
echo "Request payload:"
cat << EOF
{
  "username": "testuser",
  "email": "test@example.com",
  "password": "testpassword123"
}
EOF

echo ""
echo "Response:"
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "testpassword123"
  }' \
  -w "\nHTTP Status: %{http_code}\nTotal Time: %{time_total}s\n"

echo ""
echo "3. Check server logs for any error messages"
echo "==========================================" 