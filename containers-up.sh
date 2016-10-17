#!/bin/bash
set -e

image="gomicro"
network="gomicro-net"

# build image
if [ ! "$(docker images -q  $image)" ];then
    docker build -t $image .
fi

# create network
if [ ! "$(docker network ls --filter name=$network -q)" ];then
    docker network create $network
fi

# start workers
docker run -d -p 3020:3000 \
--name "${image}-worker-alpha" \
--network "$network" \
--restart unless-stopped \
-e ENV="DEBUG" \
-e ROLE="worker" \
-e PORT="3000" \
$image 

docker run -d -p 3030:3000 \
--name "${image}-worker-beta" \
--network "$network" \
--restart unless-stopped \
-e ENV="DEBUG" \
-e ROLE="worker" \
-e PORT="3000" \
$image 

# start proxy
docker run -d -p 3010:3000 \
--name "${image}-proxy" \
--network "$network" \
--restart unless-stopped \
-e ENV="DEBUG" \
-e ROLE="proxy" \
-e PORT="3000" \
-e ENDPOINTS="http://${image}-worker-alpha:3000,http://${image}-worker-beta:3000" \
$image 
