package io

import (
	"fmt"
	"os"
	"syscall"

	"github.com/rs/zerolog/log"
	"golang.org/x/term"
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

func AskForMasterKey() string {
	PrintQuestion("Master Key:")
	passkey, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		os.Exit(1)
	}
	return string(passkey)
}

func AskForPassword() string {
	PrintQuestion("Password:")
	passkey, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		os.Exit(1)
	}
	return string(passkey)
}
