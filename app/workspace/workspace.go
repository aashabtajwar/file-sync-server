package workspace

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/aashabtajwar/th-server/app/tokenmanager"
	"github.com/aashabtajwar/th-server/app/users"
	"github.com/aashabtajwar/th-server/database"
	"github.com/aashabtajwar/th-server/errorhandling"
	"github.com/aashabtajwar/th-server/tcpserver"
)

// var permissionChange = make(map[string]int)

func resolveRequestBody(r *http.Request) map[string]string {
	body, err := io.ReadAll(r.Body)
	errorhandling.RequestBodyReadingError(err)
	requestBodyData := make(map[string]string)
	err = json.Unmarshal(body, &requestBodyData)
	if err != nil {
		fmt.Println("Error Unmarshalling Json\n", err)
	}
	return requestBodyData
}

func ViewWorkspaceFiles(w http.ResponseWriter, r *http.Request) {
	// token := r.Header["Authorization"][0]
	// claims := tokenmanager.DecodeToken(token)
	// user_id := claims["id"]
	body, err := io.ReadAll(r.Body)
	requestBodyData := make(map[string]string)
	err = json.Unmarshal(body, &requestBodyData)
	workspaceID := requestBodyData["workspace_id"]

	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/filesync")
	var fileNames []string
	q := fmt.Sprintf("SELECT filename FROM workspace_files WHERE workspace_id='%s'", workspaceID)
	rows, err := db.Query(q)

	if err != nil {
		fmt.Println("Error fetching database rows\n", err)
	}

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			fmt.Println("DB Scan Error\n", err)
		}
		fileNames = append(fileNames, name)
	}
	payload := make(map[string][]string)
	payload["file_names"] = fileNames
	jsonResponse, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error Marshalling json\n", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)

}

func ViewFileVersions(w http.ResponseWriter, r *http.Request) {
	// version no. - timestamp
	versions := make(map[string]string)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error Reading Body\n", err)
	}
	requestBodyData := make(map[string]string)
	err = json.Unmarshal(body, &requestBodyData)

	match := requestBodyData["workspace_name"] + "_" + requestBodyData["file_name"]
	fmt.Println(match)
	entries, err := os.ReadDir("./storage/")

	if err != nil {
		fmt.Println("Error Reading Directory\n", err)
	}
	versionsString := `{`
	for _, e := range entries {
		if strings.Contains(e.Name(), match) {
			splittedFileName := strings.Split(e.Name(), "_")
			versions[splittedFileName[3]] = splittedFileName[4]
		}
	}

	count := 0
	for i, n := range versions {

		x := fmt.Sprintf(`"%s": "%s"`, i, n)
		fmt.Println(x)
		if count != len(versions)-1 {
			x += ","
		}
		versionsString += x
		count++
	}
	versionsString += `}`

	response := make(map[string]string)
	response["message"] = "Success Fetching File Versions"
	response["versions"] = versionsString
	jsonResponse, err := json.Marshal(response)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)

}

func ViewPersonalWorkspaces(w http.ResponseWriter, r *http.Request) {

	token := r.Header["Authorization"][0]
	claims := tokenmanager.DecodeToken(token)
	user_id := claims["id"]

	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/filesync")

	errorhandling.DbConnectionError(err)

	var workspaceIDs []map[string]string

	q := fmt.Sprintf("SELECT name, workspace_id FROM workspace WHERE user_id='%s'", user_id.(string))
	rows, err := db.Query(q)

	errorhandling.DbQueryError(err)

	for rows.Next() {
		var workspace_id string
		var workspace_name string
		keyValue := make(map[string]string)
		if err := rows.Scan(&workspace_name, &workspace_id); err != nil {
			fmt.Println("DB Row Scan Error\n", err)
		}
		keyValue[workspace_name] = workspace_id
		workspaceIDs = append(workspaceIDs, keyValue)
		fmt.Println("===> ", workspaceIDs)
	}
	fmt.Println(workspaceIDs)
	payload := make(map[string][]map[string]string)
	payload["workspaces"] = workspaceIDs
	fmt.Println("PAYLOAD = ", payload)
	jsonResponse, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error Marshalling Json\n", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

// Shared
func ViewWorkspaces(w http.ResponseWriter, r *http.Request) {
	token := r.Header["Authorization"][0]
	claims := tokenmanager.DecodeToken(token)
	user_id := claims["id"]

	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/filesync")
	errorhandling.DbConnectionError(err)
	var workspaceIDs []map[string]string
	fmt.Println("user id == ", user_id.(string))
	q := fmt.Sprintf("SELECT workspace.name, workspace.workspace_id FROM workspace INNER JOIN shared_workspace ON workspace.workspace_id=shared_workspace.workspace_id WHERE shared_workspace.user_id='%s'", user_id.(string))
	rows, err := db.Query(q)
	fmt.Println("rows = ", rows)

	errorhandling.DbConnectionError(err)

	for rows.Next() {
		var workspace_id string
		var workspace_name string
		keyValue := make(map[string]string)
		if err := rows.Scan(&workspace_name, &workspace_id); err != nil {
			fmt.Println("DB Row Scan Error\n", err)
		}
		keyValue[workspace_name] = workspace_id
		workspaceIDs = append(workspaceIDs, keyValue)
		fmt.Println("===> ", workspaceIDs)
	}
	fmt.Println(workspaceIDs)
	payload := make(map[string][]map[string]string)
	payload["workspaces"] = workspaceIDs
	fmt.Println("PAYLOAD = ", payload)
	jsonResponse, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error Marshalling Json\n", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)

}

// download version 2.0
// just send the files for the intended workspace
func DownloadV2(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request came")
	token := r.Header["Authorization"][0]
	claims := tokenmanager.DecodeToken(token)
	user_id := claims["id"]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading Request Body\n", err)
	}

	data := make(map[string]string)
	er := json.Unmarshal(body, &data)

	if er != nil {
		fmt.Println("Error Unmarshalling Data\n", er)
	}

	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/filesync")
	var workspaceName string
	q := fmt.Sprintf("SELECT name FROM workspace WHERE workspace_id='%s'", data["workspace_id"])
	if err := db.QueryRow(q).Scan(&workspaceName); err != nil {
		fmt.Println("Error Querying Row\n", err)
	}

	// send file
	fmt.Println("sending files now...")
	fmt.Println(data["workspace_id"])
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Okay"))

	tcpserver.SendFiles(workspaceName, data["workspace_id"], user_id.(string))

}

func Download(w http.ResponseWriter, r *http.Request) {
	// token := r.Header["Authorization"][0]
	// userId := tokenmanager.DecodeToken(token)["id"]

	body, err := io.ReadAll(r.Body)

	if err != nil {
		fmt.Println("Error reading request body:\n", err)
	}

	data := make(map[string]string)
	payload := make(map[string]string)
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

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonRes)
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

		fmt.Println("PRINTING USER ID = ", user_id)

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
		if err != nil {
			users.DatabaseError(err, writer)
		}
		defer db.Close()

		insert := "INSERT INTO workspace(name, user_id) VALUES ('" + data["name"] + "', '" + user_id.(string) + "')"

		_, err = db.Query(insert)
		if err != nil {
			users.DatabaseError(err, writer)
		}
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
			fmt.Println("Error marshalling json\n", err)
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

		errorhandling.DbConnectionError(err)

		defer db.Close()
		fmt.Println("workspace function here -->")

		query := fmt.Sprintf("SELECT id FROM users WHERE email='%s'", email)
		if er := db.QueryRow(query).Scan(&user_id); er != nil {
			fmt.Println("Error Fetching Row:\n", er)
		}
		fmt.Println(user_id) // successfully returns user_id
		fmt.Println("Workspace ID --> ", workspace_id)
		// converting workspace_id to int
		workspace_id_int, err := strconv.Atoi(workspace_id)
		if err != nil {
			fmt.Println("erorr converting workspace id to integer\n", err)
		}

		insertQuery := fmt.Sprintf("INSERT INTO shared_workspace (workspace_id, author_id, user_id) VALUES ('%d', '%s', '%s')", workspace_id_int, author_id.(string), strconv.Itoa(user_id))
		response, err := db.Exec(insertQuery)

		errorhandling.DbInsertError(err)

		fmt.Println(response)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User Added to Workspace"))

		// send notification message via tcp
		// msg := fmt.Sp	rintf()
		tcpserver.SendNotificationMessage("Workspace Has Been Shared With You", strconv.Itoa(user_id))

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 - Method Not Allowed"))
	}

}

func ViewAddedUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("came here...")
	requestBodyData := resolveRequestBody(r)
	db := database.DBInit()
	q := fmt.Sprintf(`
		SELECT users.id, users.fname, users.lname
		FROM shared_workspace
		INNER JOIN users ON users.id = shared_workspace.user_id
		WHERE shared_workspace.workspace_id='%s'
	`, requestBodyData["workspace_id"])

	var data []map[string]string

	rows, err := db.Query(q)
	errorhandling.DbQueryError(err)
	for rows.Next() {
		var userID string
		var firstName string
		var lastName string
		keyValue := make(map[string]string)
		if err := rows.Scan(&userID, &firstName, &lastName); err != nil {
			fmt.Println("DB row scan error\n", err)
		}
		keyValue[firstName+" "+lastName] = userID
		data = append(data, keyValue)
	}
	payload := make(map[string][]map[string]string)
	payload["users"] = data
	jsonresponse, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error Marshalling Json\n", err)
	}
	fmt.Println("sending response back...")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonresponse)
}

func sendJsonResponse(payload map[string]string, w http.ResponseWriter) {
	jsonresponse, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error Marshalling Json\n", err)
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonresponse)
}

func SetPermission(w http.ResponseWriter, r *http.Request) {
	fmt.Println("changing permission...")
	requestBodyData := resolveRequestBody(r)
	db := database.DBInit()
	set := 0
	if requestBodyData["permission"] == "write" {
		set = 1
	}
	insert := fmt.Sprintf(`
		UPDATE shared_workspace
		SET permission=%d
		WHERE user_id="%s" AND workspace_id="%s"
	`, set, requestBodyData["user_id"], requestBodyData["workspace_id"])
	_, err := db.Query(insert)
	if err != nil {
		users.DatabaseError(err, w)
	}
	// send message through tcp
	tcpserver.SendPermissionGrant(requestBodyData["workspace_name"], requestBodyData["workspace_id"], requestBodyData["user_id"])
	payload := make(map[string]string)
	sendJsonResponse(payload, w)

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
