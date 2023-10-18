package loadenv

import (
	"log"

	"github.com/joho/godotenv"
)

func Loadenv() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Panic(err)
	}
}
