package main

import (
	"bufio"
	"fmt"
	"os"
)

func write(file *os.File, data string) {
	_, _ = file.WriteString(data)
}

func read(path string) string {
	file, _ := os.Open(path)
	defer file.Close()
	d := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		d += scanner.Text()
	}
	return d
}

func main() {

	//Create the file
	file, _ := os.Open("file.txt")
	defer file.Close()
	_, _ = file.WriteString("hello world\n")
	write(file, "hello world how are you \n")

	fmt.Println(read("file.txt"))

	return
	state, er := os.Stat("renamedFile.txt")
	fmt.Print("result : ", os.IsNotExist(er), state.ModTime())
	return
	//Rename The file name
	os.Rename("file.txt", "renamedFile.txt")

}
