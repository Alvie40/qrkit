#!/bin/bash

# QRKit End-to-End Testing Script
# This script tests the complete video calling workflow

echo "🎯 QRKit Video Calling - End-to-End Test"
echo "========================================"
echo "Time: $(date)"
echo ""

# Function to test HTTP endpoint
test_endpoint() {
    local url=$1
    local description=$2
    local expected_status=${3:-200}
    
    echo -n "Testing $description... "
    local response=$(curl -s -o /dev/null -w "%{http_code}" "$url" 2>/dev/null)
    
    if [ "$response" = "$expected_status" ]; then
        echo "✅ OK ($response)"
        return 0
    else
        echo "❌ FAILED ($response)"
        return 1
    fi
}

# Function to test API with JSON response
test_api() {
    local url=$1
    local description=$2
    local method=${3:-GET}
    
    echo -n "Testing $description... "
    local response=$(curl -s -X "$method" "$url" 2>/dev/null)
    
    if [ -n "$response" ] && [[ "$response" != *"error"* ]]; then
        echo "✅ OK"
        echo "   Response: $response"
        return 0
    else
        echo "❌ FAILED"
        echo "   Response: $response"
        return 1
    fi
}

# Test basic endpoints
echo "📡 Testing Basic Endpoints"
echo "-------------------------"
test_endpoint "http://localhost:8080/" "Main page"
test_endpoint "http://localhost:8080/empregado" "Employee page"
test_endpoint "http://localhost:8080/debug" "Debug page"
test_endpoint "http://localhost:8080/test" "Test page"
test_endpoint "http://localhost:8080/comprehensive-test" "Comprehensive test page"
test_endpoint "http://localhost:8080/status" "Status dashboard"
echo ""

# Test API endpoints
echo "🔌 Testing API Endpoints"
echo "------------------------"
test_api "http://localhost:8080/api/create-session" "Session creation" "POST"

# Get session for further testing
echo ""
echo "🎬 Testing Video Session Workflow"
echo "--------------------------------"
SESSION_DATA=$(curl -s -X POST http://localhost:8080/api/create-session 2>/dev/null)
if [[ $SESSION_DATA == *"sessionId"* ]]; then
    SESSION_ID=$(echo "$SESSION_DATA" | grep -o '"sessionId":"[^"]*"' | cut -d'"' -f4)
    ROOM_NAME=$(echo "$SESSION_DATA" | grep -o '"roomName":"[^"]*"' | cut -d'"' -f4)
    
    echo "✅ Session created successfully"
    echo "   Session ID: $SESSION_ID"
    echo "   Room Name: $ROOM_NAME"
    
    # Test QR code generation
    test_endpoint "http://localhost:8080/api/qr/$SESSION_ID" "QR code generation"
    
    # Test client page with session
    test_endpoint "http://localhost:8080/cliente/$SESSION_ID" "Client page with session"
    
    # Test token generation for employee
    test_api "http://localhost:8080/api/token?room=$ROOM_NAME&identity=employee&isEmployee=true" "Employee token generation"
    
    # Test token generation for client
    test_api "http://localhost:8080/api/token?room=$ROOM_NAME&identity=client&isEmployee=false" "Client token generation"
    
else
    echo "❌ Failed to create session"
    echo "   Response: $SESSION_DATA"
fi

echo ""
echo "🧪 Testing Static Assets"
echo "------------------------"
test_endpoint "http://localhost:8080/static/comprehensive-test.html" "Comprehensive test static page"
test_endpoint "http://localhost:8080/static/status-dashboard.html" "Status dashboard static page"

echo ""
echo "🏥 Testing LiveKit Integration"
echo "-----------------------------"
# Check if LiveKit server is accessible
if curl -s --connect-timeout 5 http://localhost:7880 >/dev/null 2>&1; then
    echo "✅ LiveKit server: ACCESSIBLE"
else
    echo "❌ LiveKit server: NOT ACCESSIBLE"
    echo "   Make sure LiveKit server is running on port 7880"
fi

echo ""
echo "📋 Test Summary"
echo "==============="
echo "Test completed at: $(date)"
echo ""
echo "🔗 Quick Access URLs:"
echo "   Main: http://localhost:8080/"
echo "   Employee: http://localhost:8080/empregado"
echo "   Debug: http://localhost:8080/debug"
echo "   Comprehensive Test: http://localhost:8080/comprehensive-test"
echo "   Status Dashboard: http://localhost:8080/status"
echo ""
echo "🎯 Next Steps:"
echo "   1. Open http://localhost:8080/comprehensive-test for automated testing"
echo "   2. Open http://localhost:8080/status for real-time monitoring"
echo "   3. Test video calling workflow manually"
echo ""
