package main

import "fmt"

func main() {
	fmt.Printf("hello, world\n")
	a := []int{1, 2, 3}
	b := make([]int, 3)
	b = append(b, 1)
	fmt.Printf("%T,%v\n", a, a)
}
