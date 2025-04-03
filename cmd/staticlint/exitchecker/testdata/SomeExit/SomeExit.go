package Pkg

import "fmt"

type t struct{}

// Метод Exit для стороннего типа
func (*t) Exit() {
	fmt.Println("This is exit")
}

// Функция main с вызовом Exit для стороннего типа
func main() {

	var tp t

	tp.Exit() // want
}
