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





BIZU!

Método 1 :

Suponha que o nome do seu projeto sejaMyProject
Vá para o seu caminho, corrago build
Ele criará um arquivo executável como o nome do seu projeto ("MyProject")
Em seguida, execute o executável usando./MyProject
Você pode fazer as duas etapas ao mesmo tempo digitando go build && ./MyProject. Os arquivos Go do package mainsão compilados em um executável.

Método 2 :

Basta correr go run *.go. Ele não criará nenhum executável, mas será executado.

