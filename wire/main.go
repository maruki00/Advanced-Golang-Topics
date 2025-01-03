package main

import "fmt"



type Hello struct {

	Name string
}

func main() {

	fmt.Println("hello world ...")

	app := InintApp(make(map[int]any))


	fmt.Println(app.srv.GetName())



}
