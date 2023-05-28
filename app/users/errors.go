package users

import (
	"fmt"
	"io"
	"net/http"
)

func InternalError(message string, err error, writer http.ResponseWriter) {
	fmt.Printf(message)
	io.WriteString(writer, message)
}
func DatabaseError(err error, writer http.ResponseWriter) {
	fmt.Printf("Database Error\n %s", err)
	writer.WriteHeader(http.StatusInternalServerError)
	writer.Write([]byte("500 Internal Server Error"))
}
