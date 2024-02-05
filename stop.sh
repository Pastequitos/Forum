#!/bin/bash

echo "deleting the container"
result=$( docker ps -q -f name=forum-image )
if [[ $? -eq 0 ]]; then
echo "Container exists"
docker container rm -f forum-image
echo "Deleted the existing docker container"
else
echo "No such container"
fi
echo "Stopped the container:"
docker stop forum