package lfstack

import (
	"reflect"
	"slices"
	"sync"
	"testing"
)

type foobaz struct {
	foo string
	baz int
}

func TestStack(t *testing.T) {
	t.Run("Check basic functionnalities", func(t *testing.T) {
		stack := &Stack[any]{}

		var l uint64
		checklen := func() {
			if stack.Len() != l {
				t.Fatalf("stack length did not match.\n\t     got: %v\n\texpected: %v", stack.Len(), l)
			}
		}
		checklen()

		items := []any{1, "Hello world", nil}
		for _, v := range items {
			stack.Push(v)
			l += 1
			checklen()
		}

		for j := len(items) - 1; j >= 0; j-- {
			e := stack.Pop()
			i := items[j]
			if e != i {
				t.Fatalf("items did not match.\n\t     got: %v\n\texpected: %v", e, i)
			}
			l -= 1
			checklen()
		}
	})

	t.Run("Pop on empty does not error and return zero value for T", func(t *testing.T) {
		s1 := &Stack[any]{}
		if s1.Pop() != nil {
			t.Fatalf("items did not match.\n\t     got: %v\n\texpected: %v", s1.Pop(), nil)
		}

		s2 := &Stack[string]{}
		if s2.Pop() != "" {
			t.Fatalf("items did not match.\n\t     got: %v\n\texpected: %v", s2.Pop(), "")
		}

		s3 := &Stack[int]{}
		if s3.Pop() != 0 {
			t.Fatalf("items did not match.\n\t     got: %v\n\texpected: %v", s3.Pop(), 0)
		}

		s4 := &Stack[foobaz]{}
		zero := foobaz{}
		if s4.Pop() != zero {
			t.Fatalf("items did not match.\n\t     got: %v\n\texpected: %v", s4.Pop(), zero)
		}
	})

	t.Run("Concurrency all items", func(t *testing.T) {
		stack := &Stack[int]{}
		src := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		for _, v := range src {
			v := v
			go func() {
				stack.Push(v)
			}()
		}

		c := make(chan int)

		var wg sync.WaitGroup
		wg.Add(len(src))
		go func() {
			defer close(c)
			wg.Wait()
		}()

		for i := 0; i < len(src); i++ {
			go func() {
				for stack.Len() == 0 {
				}
				c <- stack.Pop()
				wg.Done()
			}()
		}

		dst := make([]int, 0, len(src))
		for v := range c {
			dst = append(dst, v)
		}

		slices.Sort(src)
		slices.Sort(dst)
		if !reflect.DeepEqual(src, dst) {
			t.Fatalf("retrieved items did not match.\n\t     got: %v\n\texpected: %v", dst, src)
		}
	})
}
