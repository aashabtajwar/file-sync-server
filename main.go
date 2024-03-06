package main

import (
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

// home test func
func home(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Home Page\n")
}

func authTest(w http.ResponseWriter, r *http.Request) {
	token := r.Header["Authorization"][0]
	decodedString := tokenmanager.DecodeToken(token)["id"]
	fmt.Println(reflect.TypeOf(decodedString))
}

func testURL(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query()
	fmt.Println(url)
}

func main() {

	// start db
	// database.DBInit()

	// consist of both HTTP and TCP services

	// TCP server
	go tcpserver.Start()

	// http for user services
	mux := http.NewServeMux()

	mux.HandleFunc("/", home)

	// validate token
	mux.HandleFunc("/validate-token", users.ValidateUser)

	// test url query
	mux.HandleFunc("/test/", testURL)

	// auth and reg
	mux.HandleFunc("/register", users.Register)
	mux.HandleFunc("/login", users.Login)

	mux.HandleFunc("/auth", authTest)
	// create content folder

	fmt.Println("Turning on server")

	// users that have been shared with
	// mux.HandleFunc("/shared-users)

	// create workspace
	mux.HandleFunc("/createw", workspace.Create)

	// show workspace files
	mux.HandleFunc("/workspace/", workspace.ShowFilesInWorkspace)

	// add user to workspace
	mux.HandleFunc("/add-user", workspace.AddUserToWorkspace)

	// check shared users
	mux.HandleFunc("/shared-users", workspace.ViewAddedUsers)

	// set permission
	mux.HandleFunc("/set-permission", workspace.SetPermission)

	// view shared remote workspaces to user
	mux.HandleFunc("/check", workspace.ViewWorkspaces)

	// personal remote
	mux.HandleFunc("/check-remote", workspace.ViewPersonalWorkspaces)

	// download workspace
	mux.HandleFunc("/download", workspace.DownloadV2)

	// check file versions
	mux.HandleFunc("/versions", workspace.ViewFileVersions)

	// workspace-files
	mux.HandleFunc("/workspace-files", workspace.ViewWorkspaceFiles)

	// exec.Command("xdg-open", "http:/127.0.0.1:3333/").Run()

	err := http.ListenAndServe("192.168.0.103:3333", mux)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println("Hello world")
}
