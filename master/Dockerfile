FROM golang:latest

USER root

RUN mkdir /home/master

WORKDIR /home/master

COPY . . 

RUN chmod 777 clearMaster.sh

RUN chmod 777 executerWorkerContainer.sh

RUN chmod 777 dockerInstall.sh

RUN ./dockerInstall.sh

RUN service docker start

RUN usermod -aG docker root

RUN chown root:docker /var/run/docker.sock

RUN docker load < /home/master/worker.tar

RUN go build -o master

EXPOSE 8080

CMD ["./master"]


