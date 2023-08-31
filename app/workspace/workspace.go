package workspace

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/aashabtajwar/th-server/app/tokenmanager"
	"github.com/aashabtajwar/th-server/app/users"
)

func Create(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {

		/*
		* get the user id from jwt token
		* create new workspace
		 */
		token := request.Header["Authorization"][0]
		claims := tokenmanager.DecodeToken(token)
		user_id := claims["id"]

		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			io.WriteString(writer, "Problem occured while dealing with data")
		}
		fmt.Println("Printing request body -> ")
		var data map[string]string
		er := json.Unmarshal(body, &data)
		fmt.Println("printing workspace name -> " + data["name"])
		if er != nil {
			users.InternalError("Could Not Unmarshal Json Data", er, writer)
		}

		db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/filesync")
		defer db.Close()
		if err != nil {
			users.DatabaseError(err, writer)
		}

		insert := "INSERT INTO workspace(name, user_id) VALUES ('" + data["name"] + "', '" + user_id.(string) + "')"

		res, err := db.Query(insert)
		if err != nil {
			users.DatabaseError(err, writer)
		}
		fmt.Println(res)
		writer.WriteHeader(http.StatusCreated)
		writer.Write([]byte("New Workspace Created"))

	} else {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		writer.Write([]byte("405 - Method Not Allowed"))
	}
}

func AddUserToWorkspace(w http.ResponseWriter, r *http.Request) {
	/*
	* add the user to the workspace
	* do not give all permissions
	* that is, do not give the same permission as the author
	 */
}

func DeleteWorksapce(w http.ResponseWriter, r *http.Request) {
	/* delete the workspace */

}

func RemoveUserFromWorkspace(w http.ResponseWriter, r *http.Request) {}

func MakeUserAnAuthor(w http.ResponseWriter, r *http.Request) {}

func ShowFilesInWorkspace(w http.ResponseWriter, r *http.Request) {
	// read id url
	// then query from db

	// authorize
	// get userId from token
	// query the workspace userId foreign key using workspace id from url query
	// if these two match, authorize
	token := r.Header["Authorization"][0]
	user_id := tokenmanager.DecodeToken(token)["id"]
	fmt.Println("Token User: ", user_id)

	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/filesync")
	defer db.Close()
	if err != nil {
		fmt.Println("Error Opening Database:\n", err)
	}
	q := fmt.Sprintf("SELECT user_id FROM workspace WHERE workspace_id='%s'", r.URL.Query()["id"][0])
	var foreignUserId string
	if err := db.QueryRow(q).Scan(&foreignUserId); err != nil {
		fmt.Println("Error Querying Row:\n", err)
	}
	fmt.Println("Fetched User: ", foreignUserId)
	if err != nil {
		fmt.Println("Error Fetching Row\n", err)
	}
	if user_id == foreignUserId {
		// workspace belongs to this user

		file_id := r.URL.Query()["id"][0]
		var files []string

		queryString := fmt.Sprintf("SELECT filename FROM workspace_files WHERE workspace_id='%s'", file_id)
		rows, err := db.Query(queryString)
		if err != nil {
			fmt.Println("Error Making Query to the Database:\n", err)
		}
		for rows.Next() {
			var fileDir string
			if err := rows.Scan(&fileDir); err != nil {
				fmt.Println("Error Scanning through the queried rows\n", err)
			}
			// format file_name string to get the proper file name
			fileName := strings.Split(fileDir, "/")
			files = append(files, fileName[len(fileName)-1])
		}
		f := WorkspaceFiles{Files: files}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Println(f)
		json.NewEncoder(w).Encode(f)

	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 - Unauthorized"))
	}

}
