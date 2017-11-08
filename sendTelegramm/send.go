package sendTelegramm

import (
	"github.com/Syfaro/telegram-bot-api"
	"database/sql"
	"log"
)

func SendMessageToChat(chatId int64, reply string, bot *tgbotapi.BotAPI){

	msg := tgbotapi.NewMessage(chatId, reply)
	bot.Send(msg)
}

func SendNextGamer(db sql.DB, randGamers *map[string]string, lastGamer *string, bot *tgbotapi.BotAPI){

	reply := ""
	for _, name := range *randGamers {
		*lastGamer = name
		reply = name
		delete(*randGamers, name)
		break
	}

	row, err := db.Query("SELECT chat_id FROM gamers")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	var chatId int64
	for row.Next() {
		err := row.Scan(&chatId)
		if err != nil {
			log.Fatal(err)
		}
		SendMessageToChat(chatId, reply, bot)
		break
	}
}
