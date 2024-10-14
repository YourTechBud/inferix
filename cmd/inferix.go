package main

import (
	"fmt"
	"net/http"

	"github.com/YourTechBud/inferix/server"
)

func main() {
	// TODO: Get a cli creation package like cobra

	// Create the server
	router, err := server.New()
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting server on port 4386")
	http.ListenAndServe(":4386", router)
}
