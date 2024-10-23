package user

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"wishtournament/util/hashing"
	"wishtournament/util/jwt"
)

func CreateNewUser(c *gin.Context, db *sql.DB) {
	log.Println("Users called")

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
