#!/bin/bash
# Build the Docker image
docker image build -f Dockerfile -t forumimg .
# Run the Docker container
docker container run -p 8080:8080 --detach --name forumcontainer forumimg