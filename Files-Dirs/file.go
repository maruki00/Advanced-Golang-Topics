package main

import (
	"fmt"
	"os"
	"time"
)

func main() {

	//Create the file
	emptyFile, err := os.Create("file.txt")
	defer emptyFile.Close()
	if err != nil {
		//panic("Could Not Create file.txt.")
	}

	state, _ := os.Stat("renamedFile.txt")
	fmt.Print(state.ModTime().Add(time.Hour))
	return
	//Rename The file name
	os.Rename("file.txt", "renamedFile.txt")
	fmt.Println(emptyFile)

}
