package Pkg

import "fmt"

type t struct{}

func (*t) Exit() {
	fmt.Println("This is exit")
}
func main() {

	var tp t

	tp.Exit() // want
}
