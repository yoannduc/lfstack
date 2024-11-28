// Copyright 2024 yoannduc. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the
// LICENSE file.

// Package lfstack provides lock free stack of any value type. It uses atomic
// functions to perform its synchronization.

package lfstack

import (
	"sync/atomic"
)

// elem is an item containing a value of type T on a Stack of type T. An elem
// also contains its next item on the stack.
type elem[T any] struct {
	next atomic.Pointer[elem[T]]
	v    T
}

// Stack is a lock free stack of items with value of type T. The zero value is
// an empty stack.
type Stack[T any] struct {
	top atomic.Pointer[elem[T]]
	len atomic.Uint64
}

// Push adds an item on top of the stack.
func (s *Stack[T]) Push(v T) {
	// item will be the new top item on the stack.
	item := elem[T]{v: v}
	// top is the current top item on the stack. It is updated on each iteration
	// of the for loop to represent this iteration's top item and compare it
	// with stack top and item (new top).
	var top *elem[T]

	// Use an endless for to retry until update is successfull. This way, we
	// enforce retries until the item is actually pushed onto the stack.
	for {
		// Load current top item into `top`.
		top = s.top.Load()
		// Store current top as next item for future top.
		item.next.Store(top)
		// Try to atomically update the top of the stack. If top has not changed
		// since it was loaded (i.e., no other goroutines have modified the
		// stack in the meantime), it successfully push the item by changing
		// stack's top to point to `item`. Succes is the only way to exit the
		// loop/function.
		if s.top.CompareAndSwap(top, &item) {
			// Increment stack len by one.
			_ = s.len.Add(1)
			return
		}
	}
}

// Pop removes the top item from the stack and returns its value.
func (s *Stack[T]) Pop() T {
	// top & next items. We need next to put next as new top in stack struct.
	var (
		// top is the current top item on the stack. It is updated on each
		// iteration of the for loop to represent this iteration's top item and
		// compare it with stack top and its next item that will become new top.
		top *elem[T]
		// next is the current iteration's top item's next item. It will then be
		// set as new top for the stack.
		next *elem[T]
	)

	// Use an endless for to retry until update is successfull or no value to
	// send. This way, we enforce retries until the item is actually updated to
	// next value on stack, or nil.
	for {
		// Load current top value and return early if nil, meaning nothing to
		// pop.
		if top = s.top.Load(); top == nil {
			// We use an empty named var because we cannot return direct empty
			// val or nil with generics. This way, go will send the zero value
			// of T.
			var empty T
			return empty
		}

		// Load next item on stack from top. If nil, we will be updating top as
		// nil, no further check.
		next = top.next.Load()
		// Try to atomically update the top of the stack. If top has not changed
		// since it was loaded (i.e., no other goroutines have modified the
		// stack in the meantime), it successfully remove the item from the
		// stack by changing stack's to point to `next`.
		if s.top.CompareAndSwap(top, next) {
			// Decrease the counter by one.
			_ = s.len.Add(^uint64(0))
			// Return top (popped item)'s value.
			return top.v
		}
	}
}

// Len returns the number of items on the stack.
func (s *Stack[T]) Len() uint64 {
	return s.len.Load()
}
