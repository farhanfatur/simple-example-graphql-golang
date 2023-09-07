package users

import (
	"database/sql"
	"log"
	"math/rand"

	database "github.com/farhanfatur/simple-example-graphql-golang/internal/pkg/db/migrations/postgres"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"name"`
	Password string `json:"password"`
}

func (user *User) Create() {
	query, err := database.DB.Prepare("INSERT INTO users (id, username, password) VALUES($1, $2, $3)")
	// fmt.Println("query =>", query)
	if err != nil {
		log.Fatal(err)
	}
	ID := rand.Int63n(6)
	hashPassword, err := HashPassword(user.Password)
	_, err = query.Exec(ID, user.Username, hashPassword)
	if err != nil {
		log.Fatal(err)
	}

}

func (user *User) Authenticate() bool {
	query, err := database.DB.Prepare("SELECT password FROM users WHERE username = $1")
	if err != nil {
		log.Fatal(err)
	}

	row := query.QueryRow(user.Username)

	var hashedPassword string
	err = row.Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			log.Fatal(err)
		}
	}

	return CheckHashPassword(user.Password, hashedPassword)
}

func HashPassword(pass string) (string, error) {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(pass), 14)
	return string(bytes), nil
}

func CheckHashPassword(pass, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	return err == nil
}

func GetUserIdByUsername(username string) (int, error) {
	statement, err := database.DB.Prepare("SELECT id FROM users WHERE username = $1")
	if err != nil {
		log.Fatal(err)
	}

	row := statement.QueryRow(username)

	var Id int
	err = row.Scan(&Id)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
		}
		return 0, nil
	}

	return Id, nil
}
