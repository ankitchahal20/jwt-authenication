package main

import (
	"fmt"
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

// only accessible with a valid token
func authHandler(w http.ResponseWriter, r *http.Request) {
	// validate the token

	authToken := r.Header.Get("Authorization")
	claims := jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(authToken, claims, func(token *jwt.Token) (interface{}, error) {
		// since we only use one private key to sign the tokens,
		// we also only use its public counter part to verify
		return verifyKey, nil
	})
	if err != nil {
		switch err1 := err.(type) {
		case *jwt.ValidationError: // something was wrong during the validation
			vErr := err1
			switch vErr.Errors {
			case jwt.ValidationErrorExpired:
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintln(w, "Token Expired, get a new one.")
				return
			default:
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(w, "Error while Parsing Token!")
				log.Printf("ValidationError error: %+v\n", vErr.Errors)
				return
			}
		default: // something else went wrong
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Error while Parsing Token!")
			log.Printf("Token parse error: %v\n", err)
			return
		}
	}
	if token.Valid {
		response := Response{"Authorized to the system"}
		jsonResponse(response, w)
	} else {
		response := Response{"Invalid token"}
		jsonResponse(response, w)
	}
}
