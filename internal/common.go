package internal

import (
	"go.uber.org/zap"
)

func check(err error) {
	if err != nil {
		zap.S().Panicf("Something went wrong: %q", err)
	}
}
