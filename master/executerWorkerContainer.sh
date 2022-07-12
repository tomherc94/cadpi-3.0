#!/bin/bash

docker run --cpus=1.0 --net=host --network=mongodb_rede-mongo --hostname $1 tomherc94/worker 


