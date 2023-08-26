package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"

	"github.com/aashabtajwar/th-server/app/tokenmanager"
	"github.com/aashabtajwar/th-server/app/users"
	"github.com/aashabtajwar/th-server/app/workspace"
	"github.com/aashabtajwar/th-server/tcpserver"
)

func home(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Home Page\n")
}

func authTest(w http.ResponseWriter, r *http.Request) {
	token := r.Header["Authorization"][0]
	decodedString := tokenmanager.DecodeToken(token)["id"]
	fmt.Println(reflect.TypeOf(decodedString))
}

func main() {

	// consist of both HTTP and TCP services

	buf := new(bytes.Buffer)
	fmt.Println(buf)
	// TCP server
	go tcpserver.Start()

	// http for user services
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	// auth and reg
	mux.HandleFunc("/register", users.Register)
	mux.HandleFunc("/login", users.Login)

	mux.HandleFunc("/auth", authTest)
	// create content folder

	fmt.Println("Turning on server")

	// create workspace
	mux.HandleFunc("/createw", workspace.Create)

	// show workspace files
	mux.HandleFunc("/workspace/", workspace.ShowFilesInWorkspace)

	err := http.ListenAndServe(":3333", mux)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Hello world")
}
