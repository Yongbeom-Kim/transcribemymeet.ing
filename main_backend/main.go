package main

import (
	"fmt"
	"net/http"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func main() {
	http.HandleFunc("/", getRoot)

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println(err)
	}

}
