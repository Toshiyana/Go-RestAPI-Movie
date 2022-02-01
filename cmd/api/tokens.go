package main

import (
	"backend/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/pascaldekloe/jwt"
	"golang.org/x/crypto/bcrypt"
)

var validUser = models.User{
	ID:       10,
	Email:    "me@here.com",
	Password: "$2a$12$Ks4i6N1YTmSzfoSa6VvAOOj.PeYO/TFrexKzwD0pCJPTMXQZnUSdm",
}

type Credentials struct {
	Username string `json:"email"`
	Password string `json:"password"`
}

func (app *application) Signin(w http.ResponseWriter, r *http.Request) {
	var creds Credentials

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		app.errorJSON(w, errors.New("unauthorized"))
		return
	}

	hashedPassword := validUser.Password

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(creds.Password))
	if err != nil {
		app.errorJSON(w, errors.New("unauthorized"))
		return
	}

	//----------------------------------------------------------------
	// When Using "" library for JWT
	//----------------------------------------------------------------
	var claims jwt.Claims
	claims.Subject = fmt.Sprint(validUser.ID)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add(15 * time.Minute))
	claims.Issuer = "mydomain.com"
	claims.Audiences = []string{"mydomain.com"}

	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(app.config.jwt.secret))
	if err != nil {
		app.errorJSON(w, errors.New("error signing"))
		return
	}
	app.writeJSON(w, http.StatusOK, string(jwtBytes), "response")

	// token, err := CreateToken(app.config.jwt.secret, validUser.ID)
	// if err != nil {
	// 	app.errorJSON(w, errors.New("error signing"))
	// 	return
	// }
	// app.writeJSON(w, http.StatusOK, token, "response")
}

// func CreateToken(secret string, userid uint64) (string, error) {
// 	atClaims := jwt.MapClaims{}
// 	atClaims["authorized"] = true
// 	atClaims["user_id"] = userid
// 	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
// 	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
// 	jwtToken, err := at.SignedString([]byte(secret))
// 	if err != nil {
// 		return "", err
// 	}
// 	return jwtToken, nil
// }
