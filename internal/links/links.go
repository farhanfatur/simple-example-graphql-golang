package links

import (
	"log"

	database "github.com/farhanfatur/simple-example-graphql-golang/internal/pkg/db/migrations/postgres"
	"github.com/farhanfatur/simple-example-graphql-golang/internal/users"
)

type Link struct {
	ID      int64
	Title   string
	Address string
	User    *users.User
}

func (link Link) Save() int64 {
	// idLink := rand.Intn(9999)

	stmt, err := database.DB.Prepare("INSERT INTO links(id, title, address, userid) VALUES($1, $2, $3, $4)")
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(link.ID, link.Title, link.Address, link.User.ID)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Row inserted!!!")

	return link.ID
}

func (link Link) GetAll() []Link {
	var links []Link
	stmt, err := database.DB.Prepare("select id, title, address from links")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var link Link
		err := rows.Scan(&link.ID, &link.Title, &link.Address)
		if err != nil {
			log.Fatal(err)
		}

		links = append(links, link)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return links
}
