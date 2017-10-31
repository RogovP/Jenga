package main

import (
	"golang-test/crud"
	"golang-test/generateMap"
	"golang-test/confirmNextGame"

	"github.com/bclicn/color"
	"fmt"
	"math/rand"
	"time"

	"log"
	"database/sql"
	_ "github.com/lib/pq"
)

var (
	name string
)

func getListGamers(db sql.DB) []string{

	rows, err := db.Query("SELECT name FROM gamers")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	gamers := []string{}
	for rows.Next() {
		err := rows.Scan(&name)
		if err != nil {
			log.Fatal(err)
		}
		gamers = append(gamers,name)
	}
	return gamers
}


func main() {
	db, err := sql.Open("postgres", "postgres://jenga:123@localhost/jenga")
	if err != nil {
		log.Fatal(err)
	}
	gamers := getListGamers(*db)
	//что случайное число не повторялось (фишка go)
	rand.Seed(time.Now().Unix())
	gamersMap := generateMap.GenerateMap(&gamers)
	destroy := false
	for {
		player := gamers[rand.Intn(len(gamers))]
		if name, ok := gamersMap[player]; ok {
			fmt.Println(color.Green(name))
			delete(gamersMap, player)

			if generateMap.Pause() == false {
				crud.SaveDestroyer(db, name)
				fmt.Println(color.Red("Destroy tower by "),name)
				destroy = true
			}
			// написать запрос на обновление count_time
		}
		if len(gamersMap) == 0 || destroy == true {
			if confirmNextGame.AskForConfirmation() == false {
				break
			} else {
				gamersMap = generateMap.GenerateMap(&gamers)
			}
		}
	}
}
