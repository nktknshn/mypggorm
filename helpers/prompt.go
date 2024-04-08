package helpers

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

func AskYesNo(question string) (bool, error) {
	reply, error := AskString(fmt.Sprintf("%v (y/n): ", question))

	if error != nil {
		return false, error
	}

	if reply == "y" || reply == "Y" {
		return true, nil
	}

	if reply == "n" || reply == "N" {
		return false, nil
	}

	return AskYesNo(question)
}

func AskString(question string) (string, error) {
	print(question)
	reader := bufio.NewReader(os.Stdin)
	confirm, err := reader.ReadString('\n')

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(confirm), nil
}

func AskStringMust(question string) string {
	reply, err := AskString(question)

	if err != nil {
		panic(err)
	}

	return reply
}

func AskPassword(question string) (string, error) {
	print(question)
	pwdBytes, err := term.ReadPassword(int(os.Stdin.Fd()))

	if err != nil {
		return "", err
	}

	return string(pwdBytes), nil
}
