#!/bin/bash

docker run --cpus=0.5 --network=mongodb_rede-mongo --hostname $1 worker 


