#!/bin/bash
set -e

echo "Applying migrations..."
# Run migrations
./migrate up

echo "Starting the api"
# Start the app
exec ./athylps