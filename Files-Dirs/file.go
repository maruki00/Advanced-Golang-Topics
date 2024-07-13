package main

import (
	"fmt"
	"os"
)

func main() {

	emptyFile, err := os.Create("file.txt")
	defer emptyFile.Close()
	if err != nil {
		panic("Could Not Create file.txt.")
	}
	fmt.Println(emptyFile)

}
