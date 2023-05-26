package users

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
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

	// read body
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

	// fmt.Println((data["body"]))

}
