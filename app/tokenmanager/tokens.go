package tokenmanager

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// change secret location
// to an .env file
var jwtKey = []byte("secret")

type JWTClaim struct {
	// Username string `json:"username"`
	Email string `json:"email"`
	Id    string `json:"id"`
	jwt.StandardClaims
}

func GenerateJWT(email string, id string) (tokenString string, err error) {
	// setting the expiration time to be one hour
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &JWTClaim{
		Email: email,
		Id:    id,

		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)
	return
}

func ValidateToken(signedToken string) (isValid bool, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)

	if err != nil {
		return false, err
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("Could no parse claims")
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("Token has expired. Login again")
	}
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func DecodeToken(token string) jwt.MapClaims {
	claims := jwt.MapClaims{}
	data, err := jwt.ParseWithClaims(token, claims, func(data *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		fmt.Println(err)
	}
	// return data
	if claims, ok := data.Claims.(jwt.MapClaims); ok && data.Valid {
		return claims
	} else {
		return nil
	}
}

//func extractClaims()
