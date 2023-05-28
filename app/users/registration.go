package users

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func InternalError(message string, err error, writer http.ResponseWriter) {
	fmt.Printf(message)
	io.WriteString(writer, message)
}

func Register(writer http.ResponseWriter, request *http.Request) {
	/*
	* ideally, there should be a validation
	* and other forms of security checks
	 */

	// check method
	if request.Method == "POST" {

		body, err := ioutil.ReadAll(request.Body)

		var data map[string]string
		if err != nil {
			InternalError("Could not read body", err, writer)
		}
		fmt.Println(body)
		er := json.Unmarshal(body, &data)
		if er != nil {
			InternalError("Could not unmarshal json body", er, writer)
		}

		// get the field values
		firstName := data["first_name"]
		lastName := data["last_name"]
		email := data["email"]
		password := data["password"]

		fmt.Println(firstName, lastName, email, password)

		// open database
		db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/filesync")
		defer db.Close()
		if err != nil {
			DatabaseError(err, writer)
		}

		insert := "INSERT INTO users(first, last, email, password) VALUES ('" + firstName + "', '" + lastName + "', '" + email + "', '" + password + "')"
		res, err := db.Query(insert)
		if err != nil {
			DatabaseError(err, writer)
		}

		fmt.Println(res)

		// success response (failure response are used at those points)
		writer.WriteHeader(http.StatusCreated)
		writer.Write([]byte("Registered"))

	} else {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		writer.Write([]byte("405 Method Not Allowed"))
	}
	// read body

	// fmt.Println((data["body"]))

}
