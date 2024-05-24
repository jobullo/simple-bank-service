package main

//import runtime "github.com/jobullo/go-api-example/cmd/console"

import runtime "github.com/jobullo/go-api-example/cmd/http"

// This is for running locally, so we
// create a simple go module here to execute
// just to give you a shorthand.
func main() {
	runtime.Execute()
}
