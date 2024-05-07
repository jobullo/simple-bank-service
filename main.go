package main

import "zeroslope/cmd/http"

// This is for running locally, so we
// create a simple go module here to execute
func main() {
	http.Execute()
}
