FROM ubuntu:20.04

USER root

RUN mkdir /home/master

WORKDIR /home/master

COPY master.go ./
COPY master.sh ./

RUN chmod 777 master.sh

RUN ./master.sh

VOLUME ./masterInput /home/worker/masterInput

RUN ls -l /home/master
