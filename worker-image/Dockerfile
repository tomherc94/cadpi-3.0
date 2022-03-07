FROM ubuntu:20.04

USER root

RUN mkdir /home/worker

WORKDIR /home/worker

COPY workerApp.jar ./
COPY workerCopy.jar ./
COPY executeWorkerApp.sh ./
COPY executeWorkerCopy.sh ./
COPY clearWorker.sh ./
COPY worker.sh ./

RUN chmod 777 worker.sh

RUN ./worker.sh

VOLUME ./workerInput /home/worker/workerInput

CMD ./executeWorkerApp.sh


