#!/bin/bash
set -e

image="gomicro"

docker rm -f $(docker ps -a -q -f "ancestor=$image")
