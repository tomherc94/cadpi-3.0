Instalar o Golang e suas dependências:
https://go.dev/doc/install
$go get go.mongodb.org/mongo-driver

Instalar Docker e Docker-compose:
https://docs.docker.com/engine/install/ubuntu/

Instalar JAVA 11:
https://www.edivaldobrito.com.br/oracle-java-11-no-ubuntu/

Desativar MONGODB LOCAL (caso esteja rodando):
$sudo /etc/init.d/mongodb stop

Executar o container do MongoDB:
$cd mongodb
$sudo docker-compose -f mongodb.yml up

Verificar execução do MongoDB via browser:
http://localhost:8081

Colocar imagens manualmente no diretório masterInput

Upload de imagens para o BD (Master):
$cd master
$go run main.go up

Download de imagens do BD (Worker):
$cd worker
$go run main.go down

Processamento de imagens:
$./executeWorkerApp.sh

Upload de imagens processadas para o BD (Worker):
$go run main.go up

Download de imagens processadas do BD (Master):
$cd master
$go run main.go down

Verificar imagens convertidas no diretório masterOutput

