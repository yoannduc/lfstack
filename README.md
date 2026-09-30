# lfstack

Go lock-free concurrent safe stack structure.

[![LICENSE](https://img.shields.io/badge/License-MIT-turquise.svg)](LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/yoannduc/lfstack.svg)](https://pkg.go.dev/github.com/yoannduc/lfstack)

## Usage

```go
package main

import (
	"fmt"

	"github.com/yoannduc/lfstack"
)

type myType struct {
	ID    int
	Label string
}

func main() {
	s := new(lfstack.Stack[myType])

	s.Push(myType{
		ID:    1,
		Label: "l1",
	})
	s.Push(myType{
		ID:    2,
		Label: "l2",
	})

	fmt.Printf("size: %d\n", s.Len()) // size: 2

	el := s.Pop()

	fmt.Printf("popped: %v\n", el)    // popped: {2 l2}
	fmt.Printf("size: %d\n", s.Len()) // size: 1
}

```
