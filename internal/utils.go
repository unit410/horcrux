package internal

import (
	"bufio"
	"io"
	"log"
	"strings"
)

func Assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

// AskForConfirmation asks the user for confirmation. A user must type "y" or "yes"
// Any other input will be considered as a No.
// Return true if the user confirmed with a yes, false otherwise.
func AskForConfirmation(source io.Reader, prompt string) bool {
	reader := bufio.NewReader(source)

	log.Printf("%s [y/N]: ", prompt)

	response, err := reader.ReadString('\n')
	Assert(err)

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}
