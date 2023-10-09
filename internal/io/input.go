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
		log.Panic().Err(err).Msg("unable to scan input")
	}
	return answer
}

func AskYesNo(question string) bool {
	answer := Ask(question + " (y/n) ")
	return answer == "y" || answer == "Y"
}

func AskForEncryptionKey() string {
	PrintQuestion("Enter vault encryption key:")
	passkey, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		os.Exit(1)
	}
	return string(passkey)
}

type StdInReader struct{}

func (r StdInReader) Read(p []byte) (n int, err error) {
	return os.Stdin.Read(p)
}
