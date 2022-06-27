FROM golang:latest

USER root

RUN mkdir /home/worker

WORKDIR /home/worker

COPY . . 

RUN chmod 777 worker.sh

RUN ./worker.sh

VOLUME /home/tomas/Documentos/UNISAL/Modulo9/TCC/worker/workerInput /home/worker/workerInput

VOLUME /home/tomas/Documentos/UNISAL/Modulo9/TCC/worker/workerOutput /home/worker/workerOutput

RUN go build -o worker

CMD ["./worker"]


