package workspace

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

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
		//insert := "INSERT INTO workspace(name, user_id) VALUES ('" + data["Name"] + "', '" + user_id + "')"

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
