Instalar o Golang:
https://go.dev/doc/install

Baixar dependencias
$go get go.mongodb.org/mongo-driver
$go get github.com/gorilla/mux

Instalar Docker e Docker-compose:
https://docs.docker.com/engine/install/ubuntu/

Desativar MONGODB LOCAL (caso esteja rodando):
$sudo /etc/init.d/mongodb stop

Executar o container do MongoDB:
$cd mongodb
$sudo docker-compose -f mongodb.yml up

Verificar execução do MongoDB via browser:
http://localhost:8081

Executar servidor (dentro do diretório master):
$go run *.go

Acessar servidor:
http://localhost:8080/upload

Selecionar banco de imagens (compactadas .zip)

Clicar em Upload

Aguardar processamento ...

Fazer Download do banco de imagens convertido
