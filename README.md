# Nats-Golang

CLI de Chat
1- inicializar el modulo de go
go mod init nats-chat-cli

2- añadir las dependecnias necesarias
go get github.com/nats-io/nats.go

3- Lanzar el servidor de Nats utilzando jetstream en local
docker run --rm -it -p 4222:4222 nats:latest

4- Compilar main.go
go build main.go

5- Ejecutar main.exe 
- paramtro 1: servidor de nats al que nos queremos conectar
- parametro 2: canal al que nos queremos conectar
- paramentro 3: usario con el que nos vamos a identificar
 main.exe --nats nats://{{ip}}:{{puerto}} --channel {{canal}} --name {{nombre}}
