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

	// auth and reg
	mux.HandleFunc("/register", users.Register)
	mux.HandleFunc("/login", users.Login)

	// create content folder

	// create workspace
	//mux.HandleFunc("/createw")

	err := http.ListenAndServe(":3333", mux)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Hello world")
}
