#!/bin/bash
#criar sistema de arquivos
mkdir /home/worker/workerInput
mkdir /home/worker/workerOutput
chmod 777 workerInput
chmod 777 workerOutput
chmod 777 workerApp.jar
chmod 777 workerCopy.jar
chmod 777 executeWorkerApp.sh
chmod 777 executeWorkerCopy.sh
chmod 777 clearWorker.sh

#instalar dependências
apt-get update
apt-get install default-jre -y

#SSHPASS
apt-get install sshpass -y

#configurar autenticação SSH
sed -i 's/prohibit-password/yes/' /etc/ssh/sshd_config
sed -i 's/PasswordAuthentication no/PasswordAuthentication yes/' /etc/ssh/sshd_config

#resetar o SSH
/etc/init.d/ssh restart

#configurar senha dos usuários
usermod -p $(openssl passwd -1 '123') root
