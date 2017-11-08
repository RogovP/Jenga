package main

import (
	"golang-test/crud"
	"golang-test/generateMap"
	"golang-test/sendTelegramm"

	"math/rand"
	"time"
	"log"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/Syfaro/telegram-bot-api"
)

var (
	botGlobal  *tgbotapi.BotAPI
	lastGamer  string
	randGamers map[string]string
)

func testForGroup(chatId int64) bool {
	if chatId == -268507686 {
		return true
	} else {
		return false
	}
}

func main() {
	bot, err := tgbotapi.NewBotAPI("")
	botGlobal = bot
	if err != nil {
		log.Panic(err)
	}

	db, err := sql.Open("postgres", "postgres://jenga:123@localhost/jenga")
	if err != nil {
		log.Fatal(err)
	}

	botGlobal.Debug = true

	log.Printf("Authorized on account %s", botGlobal.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := botGlobal.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		reply := ""
		switch update.Message.Command() {
		case "chatAdd":

		case "new":
			if testForGroup(update.Message.Chat.ID) == true {
				chat := update.Message.Chat.ID
				name := update.Message.From.FirstName
				row, err := db.Query("SELECT name FROM gamers WHERE name=$1", name)
				if err != nil {
					log.Fatal(err)
				}
				defer row.Close()

				reply = "Вы уже играете в Jenga"
				if row.Next() == false {
					reply = "Вы были добавлены в игру"
					_, err = db.Exec("INSERT INTO gamers (name, count_time, destroyer_lvl, chat_id) VALUES ($1,0,0,$2)", name, chat)
					if err != nil {
						log.Fatal(err)
					}
				}
			} else {
				reply = "Этого бота можно использовать только в тестовой группе"
			}
			sendTelegramm.SendMessageToChat(update.Message.Chat.ID, reply, botGlobal)
		case "next":
			if testForGroup(update.Message.Chat.ID) == true {
				if len(randGamers) == 0 {
					sendTelegramm.SendMessageToChat(update.Message.Chat.ID, "Начинается новый круг", botGlobal)
					gamers := crud.GetListGamers(*db)
					//что случайное число не повторялось (фишка go)
					rand.Seed(time.Now().Unix())
					randGamers = generateMap.GenerateMap(&gamers)
					sendTelegramm.SendNextGamer(*db, &randGamers, &lastGamer, botGlobal)
				} else {
					sendTelegramm.SendNextGamer(*db, &randGamers, &lastGamer, botGlobal)
				}
			} else {
				reply = "Этого бота можно использовать только в тестовой группе"
			}
		case "destroy":
			if testForGroup(update.Message.Chat.ID) == true {
				//crud.SaveDestroyer(db, lastGamer)
				reply = "На этот раз башню сломал " + lastGamer
				lastGamer = ""
				row, err := db.Query("SELECT chat_id FROM gamers")
				if err != nil {
					log.Fatal(err)
				}
				defer row.Close()
				var chat_id int64
				for row.Next() {
					err := row.Scan(&chat_id)
					if err != nil {
						log.Fatal(err)
					}
					sendTelegramm.SendMessageToChat(chat_id, reply, botGlobal)
					break
				}
				gamers := crud.GetListGamers(*db)
				//что случайное число не повторялось (фишка go)
				rand.Seed(time.Now().Unix())
				randGamers = generateMap.GenerateMap(&gamers)
			} else {
				reply = "Этого бота можно использовать только в тестовой группе"
			}
		case "start":
			reply = "Бот создан для игры в Jenga. Чтобы играть с нами необходимо отправить команду /new." + "\n" +
				"Для начала игры, а так же если Вы завершили ход, необходимо отправить /next." + "\n" +
				"Если кто - то разрушил башню - отправить команду /destroy"
			sendTelegramm.SendMessageToChat(update.Message.Chat.ID, reply, botGlobal)
		}
	}
}
