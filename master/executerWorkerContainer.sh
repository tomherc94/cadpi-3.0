#!/bin/bash

docker run --cpu=0.25 --network=mongodb_rede-mongo --hostname $1 worker 


