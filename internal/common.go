package internal

import (
	"go.uber.org/zap"
	"log"
)

func check(err error) {
	log.Println("Check", err)
	if err != nil {
		zap.S().Errorf("Something went wrong: %q", err)
		panic(err)
	}
}
