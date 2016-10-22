#!/bin/bash
set -e

image="gomicro"
network="gomicro-overlay"

# build image
if [ ! "$(docker images -q  $image)" ];then
    docker build -t $image ./src/gomicro
fi

# create network
if [ ! "$(docker network ls --filter name=$network -q)" ];then
    docker network create --driver overlay $network
fi

docker service create -p 3200:3000 \
--name "${image}-worker" \
--network "$network" \
-e ENV="PROD" \
-e ROLE="worker" \
-e PORT="3000" \
$image 

docker service create -p 3100:3000 \
--name "${image}-proxy" \
--network "$network" \
-e ENV="PROD" \
-e ROLE="proxy" \
-e PORT="3000" \
-e ENDPOINTS="http://${image}-worker:3000" \
$image 

