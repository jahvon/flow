package io

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

func Ask(question string) string {
	PrintQuestion(question)
	var answer string
	_, err := fmt.Scanln(&answer)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to scan input")
	}
	return answer
}

func AskYesNo(question string) bool {
	answer := Ask(question + " (y/n) ")
	return answer == "y" || answer == "Y"
}
