package crud

import (
	"database/sql"
	"log"
)

var (
	destroyer_lvl int
)

func SaveDestroyer(db *sql.DB, name string){
	row, err := db.Query("SELECT destroyer_lvl FROM gamers WHERE name=$1",name)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	lvl := 0
	for row.Next() {
		err := row.Scan(&destroyer_lvl)
		if err != nil {
			log.Fatal(err)
		}
		lvl = destroyer_lvl + 1
	}

	_, err = db.Exec("UPDATE gamers SET destroyer_lvl = $1 WHERE name=$2", lvl, name)
	if err != nil {
		log.Fatal(err)
	}
}
