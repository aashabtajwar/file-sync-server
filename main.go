package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/aashabtajwar/th-server/app/users"
)

func home(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Home Page\n")
}

func main() {
	// consist of both HTTP and TCP services
	// http for user services
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/register", users.Register)
	mux.HandleFunc("/login", users.Login)
	err := http.ListenAndServe(":3333", mux)
	if err != nil {
		log.Fatal(err)
	}

	// if anything goes wrong when listening, the app will close
	// find a way to restart the server and other benefial ways
	// to handle this (eg. gracefull shutdowns)

	// go func() {

	// 	err := http.ListenAndServe(":3333", mux)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }()

	fmt.Println("Hello world")
}
