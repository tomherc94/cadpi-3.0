FROM golang:latest

USER root

RUN mkdir /home/worker

WORKDIR /home/worker

COPY . . 

RUN chmod 777 worker.sh

RUN ./worker.sh

VOLUME ./workerInput /home/worker/workerInput

EXPOSE 27017:27017

RUN go build -o worker

CMD ["./worker"]


