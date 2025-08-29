#!/bin/bash

echo "Starting INSEC System..."

# Start the Go server in background
echo "Starting server..."
cd server && go run main.go &
SERVER_PID=$!

# Wait a moment for server to start
sleep 2

# Start the Rust agent
echo "Starting agent..."
cd agent/insec-agent && cargo run --release &
AGENT_PID=$!

echo "INSEC system started!"
echo "Server PID: $SERVER_PID"
echo "Agent PID: $AGENT_PID"
echo ""
echo "Press Ctrl+C to stop all services"

# Wait for user to stop
trap "echo 'Stopping services...'; kill $SERVER_PID $AGENT_PID 2>/dev/null; exit" INT
wait
