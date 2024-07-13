package main

import (
	"fmt"
	"os"
)

func main() {

	//Create the file
	emptyFile, err := os.Create("file.txt")
	defer emptyFile.Close()
	if err != nil {
		//panic("Could Not Create file.txt.")
	}

	state, er := os.Stat("renamedFile.txt")
	fmt.Print("result : ", os.IsNotExist(er), state.ModTime())
	return
	//Rename The file name
	os.Rename("file.txt", "renamedFile.txt")
	fmt.Println(emptyFile)

}
