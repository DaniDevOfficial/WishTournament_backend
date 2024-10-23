package user

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
	"wishtournament/util/hashing"
	"wishtournament/util/jwt"
)

/** AuthService */
func CreateNewUser(c *gin.Context, db *sql.DB) {

	var newUser RequestNewUser
	if err := c.ShouldBindJSON(&newUser); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error decoding request"})
		return
	}

	userId, err := GetUserIdByName(newUser.Username, db)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking for users"})
		return
	}

	if userId != -1 {
		c.JSON(http.StatusConflict, gin.H{"error": "User Does already exist"})
		return
	}

	hashedPassword, err := hashing.HashPassword(newUser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Hashing error"})
		return
	}

	userInDB := DBNewUser{
		username:      newUser.Username,
		email:         newUser.Email,
		password_hash: hashedPassword,
	}

	id, uuid, err := CreateUserInDB(userInDB, db)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error Creating User"})
		return
	}
	jwtUserData := jwt.JWTUser{
		Username: userInDB.username,
		UserId:   int(id),
		UUID:     uuid,
	}
	jwtToken, err := jwt.CreateToken(jwtUserData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating JWT"})
		return
	}
	response := struct {
		Token string `json:"token"`
	}{
		Token: jwtToken,
	}
	c.JSON(http.StatusInternalServerError, response)
}

func SignIn(c *gin.Context, db *sql.DB) {

	var credentials SignInCredentials
	if err := c.ShouldBindJSON(&credentials); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error decoding request"})
		return
	}

	userData, err := GetUserByName(credentials.Username, db)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Wrong USERNAME Or Password"})
		return
	}

	if !hashing.CheckHashedString(userData.password_hash, credentials.Password) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Wrong Username Or Password"})
		return
	}

	jwtUserData := jwt.JWTUser{
		Username: userData.username,
		UserId:   userData.user_id,
		UUID:     userData.uuid,
	}
	jwtToken, err := jwt.CreateToken(jwtUserData)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error creating JWT"})
		return
	}

	response := struct {
		Token string `json:"token"`
	}{
		Token: jwtToken,
	}
	c.JSON(http.StatusOK, response)
}

/** (C)RUD für user */

func GetUserByUUID(c *gin.Context, db *sql.DB) {

	uuid := c.Param("uuid")
	uuid = strings.Trim(uuid, " ")
	if uuid == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "No uuid attatched"})
		return
	}

	userData, err := GetUserByUUIDFromDB(uuid, db)
	if errors.Is(err, sql.ErrNoRows) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "shit hit the fan"})
		return
	}

	response := struct {
		UserData UserFromDB `json:"userData"`
	}{
		UserData: userData,
	}

	c.JSON(http.StatusOK, response)
}

func DeleteUserWithJWT(c *gin.Context, db *sql.DB) {
	jwtToken := c.Request.Header.Get("bearer")
	isValid, err := jwt.VerifyToken(jwtToken)
	if !isValid {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "JWT Token is not valid"})
	}
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "idk what happened"})
		return
	}
	decodedJWT, err := jwt.DecodeBearer(jwtToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "idk what happened¨jwt decoding"})
		return
	}
	// TODO: Do a email for validation and then handle the delete in another function
	_, err = DeleteUserInDB(decodedJWT.UserId, db)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "user wasnt deleted"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user deleted sucessfuly"})
}
