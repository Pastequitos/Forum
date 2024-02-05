#!/bin/bash

echo "Build the docker image"
docker build -f Dockerfile -t forum-image .
echo "built docker images and proceeding to delete existing container"
result=$( docker ps -q -f name=forum-image )
if [[ $? -eq 0 ]]; then
echo "Container exists"
docker container rm -f forum-image
echo "Deleted the existing docker container"
else
echo "No such container"
fi
echo "Deploying the updated container"
docker run -dp 3003:3003 --name forum --health-cmd "curl --fail http://localhost:3003/ || exit 1" --health-interval=10s forum-image
echo "Deploying the container"
echo "Link is : http://localhost:3003/"