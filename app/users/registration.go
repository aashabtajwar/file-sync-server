/*
IMPORTANT NOTICE: ERORRS NEED TO DEALT WITH PROPERLY
					RIGHT NOW THEY ARE NOT PROPERLY IMPLEMENTED
*/

package users

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

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
		username := data["username"]
		email := data["email"]
		password := data["password"]

		fmt.Println(firstName, lastName, email, password)

		// open database
		db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/filesync")
		if err != nil {
			fmt.Println("Error connecting to database\n", err)
		}
		defer db.Close()
		if err != nil {
			DatabaseError(err, writer)
		}

		hashedPassword, err := HashPassword(password)
		if err != nil {
			InternalError("", err, writer)
		}
		fmt.Println("Hashed password " + string(hashedPassword))

		q := fmt.Sprintf(`INSERT INTO users (fname, lname, username, email, password) VALUES ("%s", "%s", "%s", "%s", "%s")`, firstName, lastName, username, email, string(hashedPassword))
		// insert := "INSERT INTO users(fname, lname, username, email, password) VALUES ('" + firstName + "', '" + lastName + "', '" + email + "', '" + string(hashedPassword) + "')"
		res, err := db.Query(q)
		if err != nil {
			DatabaseError(err, writer)
		}

		fmt.Println(res)

		// success response (failure response are used at those points)
		writer.WriteHeader(http.StatusCreated)
		writer.Write([]byte("Registered"))

	} else {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		writer.Write([]byte("405 - Method Not Allowed"))
	}
}
