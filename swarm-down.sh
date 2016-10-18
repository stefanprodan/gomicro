#!/bin/bash
set -e

image="gomicro"

docker service scale ${image}-proxy=0
docker service scale ${image}-worker=0

docker service rm ${image}-proxy ${image}-worker

docker rmi -f $image
