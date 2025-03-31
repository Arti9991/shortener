package SimpleFunc

// Простая функция без Exit
func SimpleFunc(a int, b int) (int, bool) { // want
	if b > a {
		return 0, false
	} else {
		return b + a, true
	}
}
