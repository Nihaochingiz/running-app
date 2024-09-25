#!/bin/bash

# Stop and remove all containers defined in the docker-compose.yml
echo "Stopping and removing containers..."
docker-compose down

# Start the containers in detached mode
echo "Starting containers in detached mode..."
docker-compose up -d

echo "Done."