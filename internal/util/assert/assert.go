// Package assert contains simple assertations functions
package assert

import "fmt"

// C is asserted on Condition: panics if condition is true.
func C(condition bool, msg string) {
	if !condition {
		panic(msg)
	}
}

// NoErr panics if err is not nil.
func NoErr(err error, comments ...any) {
	if err != nil {
		if len(comments) > 0 {
			fmt.Println(comments...)
		}
		panic(err)
	}
}
