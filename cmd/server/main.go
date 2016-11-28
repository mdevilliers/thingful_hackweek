package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("hello world...")

	http.Handle("/", http.FileServer(http.Dir("./web")))
	http.ListenAndServe(":3000", nil)
}
