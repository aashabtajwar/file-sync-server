package users

import (
	"fmt"
	"net/http"
)

func DatabaseError(err error, writer http.ResponseWriter) {
	fmt.Printf("Database Error\n %s", err)
	writer.WriteHeader(http.StatusInternalServerError)
	writer.Write([]byte("500 Internal Server Error"))
}
