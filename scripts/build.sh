#!/bin/bash

echo "Building INSEC..."

# Build agent
echo "Building agent..."
cd agent/insec-agent
cargo build --release
cd ../..

# Build server
echo "Building server..."
cd server
go mod tidy
go build -o insec-server main.go
cd ..

# Build UI
echo "Building UI..."
cd ui
npm install
npm run build
cd ..

echo "Build complete."
