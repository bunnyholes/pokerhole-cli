#!/bin/bash
echo "🧹 Cleaning build cache..."
go clean -cache
go clean -modcache
rm -f poker-client

echo "🔨 Building fresh binary..."
go build -o poker-client cmd/poker-client/main.go

echo "✅ Clean build completed!"
echo ""
echo "Run with: ./poker-client"
