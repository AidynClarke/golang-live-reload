package main

import (
	"fmt"
	"net/http"
)

// Handler function for the root URL
func handler(w http.ResponseWriter, r *http.Request) {
	// Respond with "Hello World!"
	fmt.Fprintln(w, "Hello World!")
}

func main() {
	// Register the handler for the root URL
	http.HandleFunc("/", handler)

	// Start the server on port 3333
	fmt.Println("Example app listening on port 3333!")
	err := http.ListenAndServe(":3333", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}