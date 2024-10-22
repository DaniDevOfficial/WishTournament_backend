package user

import (
	"database/sql"
)

func GetUserIdByName(username string, db *sql.DB) (int, error) {
	sqlStmt := "SELECT id FROM users WHERE username = $1"
	row := db.QueryRow(sqlStmt, username)
	var userId int
	err := row.Scan(&userId)
	if err == sql.ErrNoRows {
		return -1, err
	}
	return userId, err
}

func GetUserByName(username string, db *sql.DB) (UserFromDB, error) {
	sql := `SELECT 
				username,
				email,
				password,
				id,
				uuid
			FROM
				users
			WHERE
				username = $1`
	row := db.QueryRow(sql, username)
	var userData UserFromDB
	err := row.Scan(&userData.username, &userData.email, &userData.password_hash, &userData.user_id, &userData.uuid)
	return userData, err
}

func GetUserPasswordHashByName(username string, db *sql.DB) (string, error) {
	sql := "SELECT password_hash FROM user WHERE username = ?"
	row := db.QueryRow(sql, username)
	var passwordHash string
	err := row.Scan(&passwordHash)
	return passwordHash, err
}

func CreateUserInDB(userData DBNewUser, db *sql.DB) (int64, string, error) {
	sql := "INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id, uuid"

	var id int64
	var uuid string

	err := db.QueryRow(sql, userData.username, userData.email, userData.password_hash).Scan(&id, &uuid)
	if err != nil {
		return -1, "", err
	}

	return id, uuid, nil
}

func GetUserById(id int, db *sql.DB) (UserFromDB, error) {

	sql := `SELECT
    			username,
    			email,
    			password_hash,
    			user_id
    		FROM
    		    user
    		WHERE 
        		user_id = ?`
	row := db.QueryRow(sql, id)

	var userData UserFromDB
	err := row.Scan(&userData.username, &userData.email, &userData.password_hash, &userData.user_id)
	return userData, err
}
