
docker run -d --name nats-server -p 4222:4222 nats:latest

mkdir Golang-Event-Driven-Arch
cd Golang-Event-Driven-Arch
go mod init Golang-Event-Driven-Arch
go get github.com/nats-io/nats.go
