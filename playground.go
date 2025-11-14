package main

import (
	"projektarbeit-go-generics/a"
	"projektarbeit-go-generics/b"
)

type customInt[T int | int8 | int16 | int32 | int64] int

type C struct {
	a.A
	b.B
} // C has two methods called f

func main() {
	c := C{}
	c.A.F()
	c.B.F()
}
