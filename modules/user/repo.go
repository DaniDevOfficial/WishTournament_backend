package user

import (
	"database/sql"
)

func GetUserIdByName(username string, db *sql.DB) (int, error) {
	sqlStmt := "SELECT user_id FROM user WHERE username = ?"
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
				password_hash,
				user_id
			FROM
				user
			WHERE
				username = ?`
	row := db.QueryRow(sql, username)
	var userData UserFromDB
	err := row.Scan(&userData.username, &userData.email, &userData.password_hash, &userData.user_id)
	return userData, err
}

func GetUserPasswordHashByName(username string, db *sql.DB) (string, error) {
	sql := "SELECT password_hash FROM user WHERE username = ?"
	row := db.QueryRow(sql, username)
	var passwordHash string
	err := row.Scan(&passwordHash)
	return passwordHash, err
}

func CreateUserInDB(userData DBNewUser, db *sql.DB) (int, error) {
	sql := "INSERT INTO user (username, email, password_hash) VALUES (?, ?, ?)"
	stmt, err := db.Prepare(sql)

	if err != nil {
		return -1, err
	}
	result, err := stmt.Exec(userData.username, userData.email, userData.password_hash)
	if err != nil {
		return -1, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(id), nil
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
