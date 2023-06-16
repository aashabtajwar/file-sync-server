package users

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Login(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			InternalError("Could not ready body", err, writer)
		}
		var bodyData map[string]string
		er := json.Unmarshal(body, &bodyData)
		if er != nil {
			InternalError("Could not unmarshal json body", er, writer)
		}

		email := bodyData["email"]
		password := bodyData["password"]

		fmt.Println(password)

		db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/filesync")
		defer db.Close()
		if err != nil {
			DatabaseError(err, writer)
		}

		query := "SELECT * FROM users WHERE email='" + email + "';"
		fmt.Println(query)
		res, err := db.Query(query)
		defer res.Close()
		if err != nil {
			DatabaseError(err, writer)
		}
		var data RegisterUser
		for res.Next() {
			err := res.Scan(&data.Id, &data.First, &data.Last, &data.Email, &data.Password)
			if err != nil {
				fmt.Println("the error")
				fmt.Println(err)
				InternalError("Error scanning data", err, writer)
			}

		}
		fmt.Println("NOW PRINTING DATA")
		fmt.Println(data.Email, data.Password)

		// fmt.Println(res.Columns())

	}
}
