#!/bin/bash

# Simple LiveKit Test Script
# Tests if LiveKit server can be started and connected to

echo "🎯 LiveKit Server Test"
echo "====================="
echo "Time: $(date)"
echo ""

# Function to check if port is in use
check_port() {
    local port=$1
    if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null ; then
        echo "✅ Port $port is in use"
        return 0
    else
        echo "❌ Port $port is available"
        return 1
    fi
}

# Function to test HTTP endpoint
test_http() {
    local url=$1
    local timeout=${2:-5}
    
    echo -n "Testing $url... "
    if curl -s --connect-timeout $timeout "$url" >/dev/null 2>&1; then
        echo "✅ ACCESSIBLE"
        return 0
    else
        echo "❌ NOT ACCESSIBLE"
        return 1
    fi
}

echo "🔍 Checking current port status..."
check_port 7880
check_port 8080

echo ""
echo "🐳 Checking Docker..."
if command -v docker >/dev/null 2>&1; then
    echo "✅ Docker is installed"
    
    # Check if Docker daemon is running
    if docker info >/dev/null 2>&1; then
        echo "✅ Docker daemon is running"
        
        # Try to start LiveKit container
        echo ""
        echo "🚀 Starting LiveKit server..."
        
        # Stop any existing LiveKit container
        docker stop qrkit-livekit 2>/dev/null || true
        docker rm qrkit-livekit 2>/dev/null || true
        
        # Start new LiveKit container
        docker run -d \
            --name qrkit-livekit \
            -p 7880:7880 \
            -p 7881:7881/udp \
            -p 50000-50100:50000-50100/udp \
            -v "$(pwd)/livekit-local.yaml:/livekit.yaml" \
            livekit/livekit-server \
            --config /livekit.yaml
        
        if [ $? -eq 0 ]; then
            echo "✅ LiveKit container started"
            
            # Wait for startup
            echo "⏳ Waiting for LiveKit to initialize..."
            sleep 10
            
            # Test connectivity
            echo ""
            echo "🔗 Testing connectivity..."
            test_http "http://localhost:7880"
            check_port 7880
            
            # Show container status
            echo ""
            echo "📊 Container status:"
            docker ps | grep qrkit-livekit || echo "❌ Container not running"
            
            # Show logs
            echo ""
            echo "📋 Recent logs:"
            docker logs --tail 10 qrkit-livekit 2>/dev/null || echo "❌ No logs available"
            
        else
            echo "❌ Failed to start LiveKit container"
        fi
        
    else
        echo "❌ Docker daemon is not running"
    fi
    
else
    echo "❌ Docker is not installed"
fi

echo ""
echo "🎯 Testing QRKit Server..."
test_http "http://localhost:8080"
test_http "http://localhost:8080/api/create-session"

echo ""
echo "📋 Summary"
echo "=========="
echo "- QRKit Server: $(test_http "http://localhost:8080" 1 && echo "RUNNING" || echo "NOT RUNNING")"
echo "- LiveKit Server: $(test_http "http://localhost:7880" 1 && echo "RUNNING" || echo "NOT RUNNING")"
echo ""
echo "🔗 URLs to test:"
echo "  - Main: http://localhost:8080/"
echo "  - Live Video Test: http://localhost:8080/live-video-test"
echo "  - Employee: http://localhost:8080/empregado"
echo ""

if test_http "http://localhost:7880" 1 >/dev/null 2>&1; then
    echo "✅ Ready for video calling tests!"
else
    echo "⚠️  LiveKit server not accessible. Video calling may not work."
    echo "   Try running: docker run -d -p 7880:7880 livekit/livekit-server"
fi

echo ""
