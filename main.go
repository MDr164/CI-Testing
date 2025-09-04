// SPDX-License-Identifier: BSD-3-Clause

package citesting

import (
	"fmt"
	"unsafe"
)

// FaultyFunction demonstrates various runtime issues while still compiling
func FaultyFunction(input []string) string {
	// Potential nil pointer dereference
	var ptr *int
	fmt.Printf("Dereferencing nil pointer: %d\n", *ptr)

	// Out of bounds slice access
	if len(input) > 0 {
		result := input[len(input)+5] // Will panic at runtime
		return result
	}

	// Unsafe pointer manipulation
	var x int = 42
	unsafePtr := unsafe.Pointer(&x)
	strPtr := (*string)(unsafePtr) // Converting int pointer to string pointer

	// Infinite recursion
	return FaultyFunction(append(input, *strPtr))
}

// AnotherFaultyFunction has memory leaks and race conditions
func AnotherFaultyFunction() {
	// Memory leak - creating goroutines that never finish
	for i := 0; i < 1000; i++ {
		go func() {
			select {} // Blocks forever
		}()
	}

	// Race condition
	counter := 0
	for i := 0; i < 100; i++ {
		go func() {
			counter++ // Concurrent access without synchronization
		}()
	}
}
