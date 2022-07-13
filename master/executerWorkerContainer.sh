#!/bin/bash

docker run --cpus=1.0 --network=mongodb_rede-mongo --hostname $1 worker 


