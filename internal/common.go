package internal

import (
	"go.uber.org/zap"
)

func check(err error) {
	if err != nil {
		zap.S().Errorf("Something went wrong: %q", err)
	}
}
