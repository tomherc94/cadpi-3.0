#!/bin/bash
#criar sistema de arquivos
mkdir ./workerInput
mkdir ./workerOutput
chmod 777 ./workerInput
chmod 777 ./workerOutput
chmod 777 ./workerApp.jar
chmod 777 ./executeWorkerApp.sh

#instalar dependências
apt-get update
apt-get install default-jre -y

#configurar senha dos usuários
#usermod -p $(openssl passwd -1 '123') root
