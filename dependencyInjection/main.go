package main

import "github.com/google/wire"

type User struct {
}

func InitUser() User {
	wire.Build(NewUser, NewUserName)
}

func main() {

}
