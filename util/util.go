package util

import (
	"fmt"
	"log"
)

func LogFatal(msg string) {
	log.Fatal(msg)
}

func LogWarn(msg string) {
	log.Println(msg)
}

func Log(msg string) {
	fmt.Println(msg)
}
