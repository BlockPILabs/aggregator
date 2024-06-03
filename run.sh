#!/bin/bash

# Define variables
IMAGE_NAME="my-aggregator-image"
CONTAINER_NAME="aggregator-container"
HOST_PORT=8012
CONTAINER_PORT=8012

# Navigate to the project root directory
cd "$(dirname "$0")"

# Stop and remove any existing container with the same name
if [ "$(docker ps -q -f name=$CONTAINER_NAME)" ]; then
    echo "Stopping existing container..."
    docker stop $CONTAINER_NAME
fi

if [ "$(docker ps -aq -f name=$CONTAINER_NAME)" ]; then
    echo "Removing existing container..."
    docker rm $CONTAINER_NAME
fi

# Remove the old image
if [ "$(docker images -q $IMAGE_NAME)" ]; then
    echo "Removing old Docker image..."
    docker rmi $IMAGE_NAME
fi

# Build the Docker image
echo "Building Docker image..."
docker build -t $IMAGE_NAME .

# Run the Docker container
echo "Running Docker container..."
docker run -d -p $HOST_PORT:$CONTAINER_PORT --name $CONTAINER_NAME $IMAGE_NAME

# Display running containers
docker ps

echo "Aggregator is now running on port $HOST_PORT"
