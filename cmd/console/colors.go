package console

import "fmt"

// printGray prints the provided message in gray color.
func printGray(message string) {
	fmt.Printf("\033[90m%s\033[0m\n", message)
}

// printGreen prints the provided message in green color.
func printGreen(message string) {
	fmt.Printf("\033[32m%s\033[0m\n", message)
}

// printBlue prints the provided message in blue color.
func printBlue(message string) {
	fmt.Printf("\033[34m%s\033[0m\n", message)
}

// printRed prints the provided message in red color.
func printRed(message string) {
	fmt.Printf("\033[31m%s\033[0m\n", message)
}

// printYellow prints the provided message in yellow color.
func printYellow(message string) {
	fmt.Printf("\033[33m%s\033[0m\n", message)
}
