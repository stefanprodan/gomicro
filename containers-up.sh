#!/bin/bash
set -e

image="gomicro"
network="gomicro-net"

# build image
if [ ! "$(docker images -q  $image)" ];then
    docker build -t $image ./src/gomicro
fi

# create network
if [ ! "$(docker network ls --filter name=$network -q)" ];then
    docker network create $network
fi

# start workers
docker run -d -p 3020:3000 \
--name "${image}-worker" \
--network "$network" \
--restart unless-stopped \
-e "SERVICE_NAME=${image}-worker" \
-e "SERVICE_TAGS=gomicro,production" \
-e ENV="PROD" \
-e ROLE="worker" \
-e PORT="3000" \
$image 

# start proxy
docker run -d -p 3010:3000 \
--name "${image}-proxy" \
--network "$network" \
--restart unless-stopped \
-e "SERVICE_NAME=${image}-proxy" \
-e "SERVICE_TAGS=gomicro,production" \
-e ENV="PROD" \
-e ROLE="proxy" \
-e PORT="3000" \
-e ENDPOINTS="http://${image}-worker:3000" \
$image 

# start monitor
docker run -d -p 3030:3000 \
--name "${image}-monitor" \
--network "$network" \
--restart unless-stopped \
-e "SERVICE_NAME=${image}-monitor" \
-e "SERVICE_TAGS=gomicro,production" \
-e ENV="PROD" \
-e ROLE="monitor" \
-e PORT="3000" \
-e ENDPOINTS="http://${image}-worker:3000,http://${image}-proxy:3000" \
$image 