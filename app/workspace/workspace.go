package workspace

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/aashabtajwar/th-server/app/tokenmanager"
	"github.com/aashabtajwar/th-server/app/users"
	"github.com/aashabtajwar/th-server/errorhandling"
)

func ViewWorkspaces(w http.ResponseWriter, r *http.Request) {
	token := r.Header["Authorization"][0]
	claims := tokenmanager.DecodeToken(token)
	user_id := claims["id"]

	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/filesync")

	var workspaceIDs []string
	// q := fmt.Sprintf("SELECT workspace_id FROM shared_workspace WHERE user_id='%s'", user_id)
	q := fmt.Sprintf("SELECT workspace.name FROM workspace INNER JOIN shared_workspace ON workspace.workspace_id=shared_workspace.workspace_id WHERE shared_workspace.user_id='%s'", user_id)
	rows, err := db.Query(q)

	errorhandling.DbConnectionError(err)

	for rows.Next() {
		var workspace_id string
		if err := rows.Scan(&workspace_id); err != nil {
			fmt.Println("DB Row Scan Error\n", err)
		}
		workspaceIDs = append(workspaceIDs, workspace_id)
	}
	for _, e := range workspaceIDs {
		fmt.Println(e)
	}

}

func Download(w http.ResponseWriter, r *http.Request) {
	// token := r.Header["Authorization"][0]
	// userId := tokenmanager.DecodeToken(token)["id"]

	body, err := io.ReadAll(r.Body)

	if err != nil {
		fmt.Println("Error reading request body:\n", err)
	}

	var data map[string]string
	var payload map[string]string
	er := json.Unmarshal(body, &data)

	if er != nil {
		fmt.Println("Error Unmarshalling Data:\n", er)
	}

	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/filesync")
	defer db.Close()
	if err != nil {
		fmt.Println("DB Opening Error:\n", err)
	}

	q := fmt.Sprintf("SELECT name FROM workspace WHERE workspace_id='%s'", data["workspace_id"])
	var workspaceName string
	if err := db.QueryRow(q).Scan(&workspaceName); err != nil {
		fmt.Println("Error Querying Row:\n", err)
		payload["message"] = "Could Not Fetch Workspace Details"

	}
	// fmt.Println(workspaceName)
	payload["message"] = "Successfully fetched"
	payload["workspace_name"] = workspaceName
	jsonRes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error Marshalling Json\n", jsonRes)
	}

	// now send files from this workspace
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonRes)
	// w.Write([]byte(workspaceName))

}

func Create(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {

		/*
		* get the user id from jwt token
		* create new workspace
		 */
		fmt.Println("received request")
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

		_, err = db.Query(insert)
		if err != nil {
			users.DatabaseError(err, writer)
		}
		// fmt.Println(insert)
		q := "SELECT workspace_id FROM workspace ORDER BY workspace_id DESC LIMIT 1"
		var lastId int
		// v, err := db.Query(q)
		if err := db.QueryRow(q).Scan(&lastId); err != nil {
			fmt.Println("Error Querying Row:\n", err)
		}

		payload := make(map[string]string)

		payload["message"] = "New Workspace Created"
		payload["workspace_id"] = strconv.Itoa(lastId)
		jsonRes, err := json.Marshal(payload)
		if err != nil {
			fmt.Println("Error marshalling json\n", jsonRes)
		}

		// respond with workspace id
		writer.WriteHeader(http.StatusCreated)
		writer.Header().Set("Content-Type", "application/json")
		writer.Write(jsonRes)

	} else {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		writer.Write([]byte("405 - Method Not Allowed"))
	}
}

func AddUserToWorkspace(w http.ResponseWriter, r *http.Request) {
	// author adds the user by giving his email
	// add the user to the workspace
	// do not give all permissions
	// that is, do not give the same permission as the author

	// the user_email is passed to request body
	// the workspace_id is passed to url
	fmt.Println(r.Method)
	if r.Method == "POST" {
		var user_id int

		// workspace_id := r.URL.Query()["workspace_id"][0]

		token := r.Header["Authorization"][0]
		author_id := tokenmanager.DecodeToken(token)["id"]

		body, err := ioutil.ReadAll(r.Body)

		errorhandling.JsonMarshallingError(err)

		var bodyData map[string]string
		er := json.Unmarshal(body, &bodyData)

		errorhandling.JsonMarshallingError(er)
		fmt.Println("workspace function here -->")

		// ideally user_id should be directly inserted
		// but here, the user_id is first queried and then inserted
		// this adds a bit more time
		// find a way to eliminate this
		email := bodyData["email"]
		workspace_id := bodyData["workspace_id"]

		db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/filesync")
		defer db.Close()

		errorhandling.DbConnectionError(er)

		fmt.Println("workspace function here -->")

		query := fmt.Sprintf("SELECT id FROM users WHERE email='%s'", email)
		if er := db.QueryRow(query).Scan(&user_id); er != nil {
			fmt.Println("Error Fetching Row:\n", er)
		}
		fmt.Println(user_id) // successfully returns user_id

		insertQuery := fmt.Sprintf("INSERT INTO shared_workspace (workspace_id, author_id, user_id) VALUES ('%s', '%s', '%s')", workspace_id, author_id.(string), strconv.Itoa(user_id))
		response, err := db.Exec(insertQuery)

		errorhandling.DbInsertError(err)

		fmt.Println(response)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User Added to Workspace"))

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 - Method Not Allowed"))
	}

}

func DeleteWorksapce(w http.ResponseWriter, r *http.Request) {
	/* delete the workspace */

}

func RemoveUserFromWorkspace(w http.ResponseWriter, r *http.Request) {}

func MakeUserAnAuthor(w http.ResponseWriter, r *http.Request) {}

func UserPermissions(w http.ResponseWriter, r *http.Request) {}

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
		queryWorkspaceFiles(db, w, r)

	} else {
		// check if other users have access to this workspace
		sharedWorkspaceQuery := fmt.Sprintf("SELECT user_id FROM shared_workspace WHERE workspace_id='%s'", r.URL.Query()["id"][0])
		if err := db.QueryRow(sharedWorkspaceQuery).Scan(&foreignUserId); err != nil {
			fmt.Println("Error Querying Row:\n", err)
		}
		fmt.Println("Fetched F-User: ", foreignUserId)
		fmt.Println("This Use: ", user_id)
		if user_id == foreignUserId {
			queryWorkspaceFiles(db, w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401 - Unauthorized"))
		}

	}

}

func queryWorkspaceFiles(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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
}
