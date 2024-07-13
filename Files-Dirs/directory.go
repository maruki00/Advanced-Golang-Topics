package main

import (
	"os"
)

func main() {
	err := os.Mkdir("dirTmp", 777)
	if err != nil {
		panic("Could Not Create Dir!")
	}
}
