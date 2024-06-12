package client

import (
	"database/sql"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

/*
VerifieNameOrEmail(Input string, db *sql.DB) (bool, error)

This function checks if the given input (either username or email) already exists in the database.
It queries the database to search for a matching username or email.

Input:

Input : string, the username or email to be checked.

db : *sql.DB, a pointer to the database connection.

Output:

bool : true if the username or email already exists, false otherwise.

error : an error if there is an issue with the database query.
*/
func VerifieNameOrEmail(Input string, db *sql.DB) (bool, error) {
	var NameOrEmail string
	// Query the database to see if the input matches either a username or an email.
	err := db.QueryRow("SELECT username FROM user WHERE email=? OR username=?", Input, Input).Scan(&NameOrEmail)
	if err != nil {
		// If no rows are returned, it means the input does not exist in the database.
		if err == sql.ErrNoRows {
			return false, nil
		}
		// If there is any other error, return false and the error.
		return false, err
	}
	// If the query returns a row, it means the input exists in the database.
	return true, nil
}

/*
IsUserCreate(Input string, db *sql.DB) (bool, error)

This function checks if a user with the given UUID exists in the database.
It queries the database to search for a matching UUID.

Input:

Input : string, the UUID to be checked.

db : *sql.DB, a pointer to the database connection.

Output:

bool : true if the user with the given UUID exists, false otherwise.

error : an error if there is an issue with the database query.
*/
func IsUserCreate(Input string, db *sql.DB) (bool, error) {
	var uuid string
	err := db.QueryRow("SELECT username FROM user WHERE uuid=?", Input).Scan(&uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // Aucun utilisateur trouvé avec cet UUID
		}
		return false, err // Erreur lors de la requête SQL
	}
	return true, nil // Utilisateur trouvé
}

/*
VerifiePwd(Input string, Password string, db *sql.DB) (bool, error)

This function verifies if the provided password matches the hashed password stored in the database for the given username or email.
It queries the database to retrieve the hashed password for the provided username or email and then compares it with the provided password.

Input:

Input : string, the username or email for which the password needs to be verified.

Password : string, the password to be verified.

db : *sql.DB, a pointer to the database connection.

Output:

bool : true if the provided password matches the hashed password, false otherwise.

error : an error if there is an issue with the database query.
*/
func VerifiePwd(Input string, Password string, db *sql.DB) (bool, error) {
	var hashedPwd string
	err := db.QueryRow("SELECT pwd FROM user WHERE email=? OR username=?", Input, Input).Scan(&hashedPwd)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // Aucun utilisateur trouvé avec cet email ou nom d'utilisateur
		}
		return false, err // Erreur lors de la requête SQL
	}
	return CheckPasswordHash(Password, hashedPwd), nil // Vérifie le hash du mot de passe
}

/*
HashPassword(password string) (string, error)

This function hashes the provided password using bcrypt with a cost factor of 14.

Input:

password : string, the password to be hashed.

Output:

string : the hashed password.

error : an error if there is an issue during the hashing process.
*/

/*
isCorrectPassword(password string) bool

This function checks if the provided password meets the criteria for being considered secure.
It verifies if the password contains at least one uppercase letter, one lowercase letter, one digit, and one special character.

Input:

password : string, the password to be checked.

Output:

bool : true if the password meets the criteria, false otherwise.
*/
func isCorrectPassword(password string) bool {
	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasUpper && hasLower && hasDigit && hasSpecial // Vérifie si toutes les conditions sont remplies
}

/*
The function take two Post struct and compare there ID and return a boolean.

input : post1 Post, post2 Post

output : boolean
*/
func comparePost(post1 Post, post2 Post) bool {
	// Return true if the IDs of the two posts are equal, indicating that they are the same post.
	return post1.ID == post2.ID
}

/*
HashPassword(password)

this function hash the passe word put by the user when he signs up

input : password string

output:

string

error
*/
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err // Retourne le mot de passe hashé
}

/*
CheckPasswordHash(password, hash string) bool

This function compares a password with its hashed version to verify if they match.

Input:

password : string, the password to be checked.

hash : string, the hashed password to compare against.

Output:

bool : true if the password matches the hashed version, false otherwise.
*/
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil // Retourne true si les mots de passe correspondent
}
