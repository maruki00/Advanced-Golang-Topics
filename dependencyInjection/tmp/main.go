package main

func main() {
	// Use the wire-generated injector to create a Greeter
	greeter := main.InitializeGreeter()InitializeGreeter().message

	// Call the Greet method
	greeter.Greet()
}
