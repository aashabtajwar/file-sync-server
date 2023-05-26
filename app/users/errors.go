package users

import (
	"fmt"
	"net/http"
)

func DbOpenError(err error, writer http.ResponseWriter) {
	fmt.Printf("Error opening database %s", err)
	writer.WriteHeader(http.StatusInternalServerError)
	writer.Write([]byte("500 Internal Server Error"))
}
