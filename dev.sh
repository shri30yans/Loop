#!/bin/bash

# Function to stop all processes on exit
cleanup() {
    echo "Stopping all services..."
    kill $(jobs -p) 2>/dev/null
    exit
}

# Set up trap for cleanup
trap cleanup EXIT INT TERM

# Start Go backend
echo "Starting Go backend..."
cd Loop_backend
go run main.go &
cd ..

# Wait for Go backend to start
sleep 2

# Start Python embedding service
echo "Starting embedding service..."
cd embedding_service
python3 app.py &
cd ..

# Wait for embedding service to start
sleep 2

# Start Next.js frontend
echo "Starting frontend..."
cd Loop_frontend
npm run dev &
cd ..

# Keep script running
echo "All services started. Press Ctrl+C to stop."
wait
