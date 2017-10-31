package generateMap

import (
	"os"
	"log"
	"bufio"
	"fmt"
)

// Delete home people
func deleteNotWorked(gamersMap *map[string]string) {
	notWorked, err := os.Open("dontWork.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer notWorked.Close()
	scanner := bufio.NewScanner(notWorked)
	for scanner.Scan() {
		for _, value := range *gamersMap {
			if scanner.Text() == value {
				delete(*gamersMap, value)
			}
		}
	}
}

//generate new map gamers
func GenerateMap(gamers *[]string) map[string]string {
	gamersMap := map[string]string{}

	for _, name := range *gamers {
		gamersMap[name] = name
	}
	deleteNotWorked(&gamersMap)

	return gamersMap
}

//pause between gamers
func Pause() bool {
	var response string
	fmt.Scanln(&response)
	if response == "end" {
		return false
	} else {
		return true
	}
}