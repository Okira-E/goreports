package utils

import "fmt"

// Log logs a string or more to the console.
func Log(str ...string) {
	for _, s := range str {
		fmt.Print(s)
	}
	fmt.Println()
}
