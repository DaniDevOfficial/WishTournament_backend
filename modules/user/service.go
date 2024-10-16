package user

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"wishtournament/util/error"
	"wishtournament/util/hashing"
	"wishtournament/util/jwt"
	"wishtournament/util/responses"
)

func CreateNewUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	log.Println("Users called")

	var newUser RequestNewUser
	err := json.NewDecoder(r.Body).Decode(&newUser)

	if err != nil {
		error.HttpResponse(w, "Error Decoding Request", http.StatusBadRequest)
		return
	}

	userId, err := GetUserIdByName(newUser.Username, db)
	if err != nil && err != sql.ErrNoRows {
		error.HttpResponse(w, "Error checking for users", http.StatusInternalServerError)
		return
	}

	if userId != -1 {
		error.HttpResponse(w, "User Does already exist", http.StatusConflict)
		return
	}

	hashedPassword, err := hashing.HashPassword(newUser.Password)
	if err != nil {
		error.HttpResponse(w, "Hashing error", http.StatusInternalServerError)
		return
	}

	userInDB := DBNewUser{
		username:      newUser.Username,
		email:         newUser.Email,
		password_hash: hashedPassword,
	}

	id, err := CreateUserInDB(userInDB, db)
	if err != nil {
		error.HttpResponse(w, "Error Creating User", http.StatusInternalServerError)
		return
	}
	jwtUserData := jwt.JWTUser{
		Username: userInDB.username,
		UserId:   id,
	}
	jwtToken, err := jwt.CreateToken(jwtUserData)
	if err != nil {
		error.HttpResponse(w, "Error creating JWT", http.StatusInternalServerError)
		return
	}
	response := struct {
		Token string `json:"token"`
	}{
		Token: jwtToken,
	}
	responses.ResponseWithJSON(w, response, http.StatusCreated)
}

func SignIn(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	var credentials SignInCredentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		error.HttpResponse(w, "Error Decoding Request", http.StatusBadRequest)
		return
	}

	userData, err := GetUserByName(credentials.Username, db)
	if err != nil {
		error.HttpResponse(w, "Wrong Password or Username", http.StatusBadRequest)
		return
	}

	if !hashing.CheckHashedString(userData.password_hash, credentials.Password) {
		error.HttpResponse(w, "Wrong Password or Username", http.StatusBadRequest)
		return
	}

	jwtUserData := jwt.JWTUser{
		Username: userData.username,
		UserId:   userData.user_id,
	}
	jwtToken, err := jwt.CreateToken(jwtUserData)
	if err != nil {
		error.HttpResponse(w, "Error creating JWT", http.StatusInternalServerError)
		return
	}

	response := struct {
		Token string `json:"token"`
	}{
		Token: jwtToken,
	}

	responses.ResponseWithJSON(w, response, http.StatusOK)
}
