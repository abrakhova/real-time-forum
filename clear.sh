#!/bin/bash

# Removes stopped containers, unused networks, dangling imgs, unused build cache
docker system prune -f 

# Remove -a all images without at least one container associated
docker image prune -a -f
