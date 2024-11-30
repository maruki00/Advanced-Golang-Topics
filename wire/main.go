package main

import "fmt"

func main() {

	fmt.Println("hello world ...")

	app := InintApp(make(map[int]any))


	fmt.Println(app.srv.GetName())



}
