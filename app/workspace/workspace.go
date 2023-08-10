package workspace

import (
	"database/sql"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/aashabtajwar/th-server/app/users"
)

func Create(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {

		// extract user data from token
		//		token :=

		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			io.WriteString(writer, "Problem occured while dealing with data")
		}
		var data map[string]string
		er := json.Unmarshal(body, &data)
		if er != nil {
			users.InternalError("Could Not Unmarshal Json Data", er, writer)
		}

		db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/filesync")
		defer db.Close()
		if err != nil {
			users.DatabaseError(err, writer)
		}
		// insert := "INSERT INTO workspace(name, )"

	} else {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		writer.Write([]byte("405 - Method Not Allowed"))
	}
}
