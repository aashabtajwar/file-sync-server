package users

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aashabtajwar/th-server/app/tokenmanager"
)

func ValidateUser(w http.ResponseWriter, r *http.Request) {

	message := make(map[string]string)
	token := r.Header["Authorization"][0]
	isValid, err := tokenmanager.ValidateToken(token)

	if err != nil {
		fmt.Println(err)
	}

	if !isValid {
		message["message"] = "Token has expired"
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		message["message"] = "Token still valid"
		w.WriteHeader(http.StatusOK)
	}
	jsonResponse, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error Marshalling Json\n", jsonResponse)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)

}
