#!/bin/bash

docker run --cpus=0.75 --memory 1024m --network=mongodb_rede-mongo --hostname $1 worker 


