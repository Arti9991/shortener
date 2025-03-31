package main

import "os"

// Функция main с вызовом os.Exit
func main() {

	os.Exit(1) // want "call to os.Exit in main func in package main"
}
