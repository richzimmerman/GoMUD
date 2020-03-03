package utils

import (
	"encoding/json"
	"strings"
)

func IndexOf(value string, slice []string) int {
	// TODO: If i need something like this for other data types, figure out how to make this generic to avoid repetative code
	for i, v := range slice {
		if v == value {
			return i
		}
	}
	return -1
}

// Checks user input for prompts/menus to see if q/quit to exit out of prompt
func CheckIfQuit(input string) bool {
	input = strings.ToLower(input)
	return input == "q" || input == "quit"
}


func JsonMarshalMapToArray(input map[string]*interface{}) (string, error) {
	var tmp []string
	for key, _ := range input {
		tmp = append(tmp, key)
	}
	b, err := json.Marshal(tmp)
	if err != nil {
		return "", Error(err)
	}
	retval := string(b)
	// To parse from the DB easier later, we want JSON formatted strings, not "null"
	if retval == "null" {
		retval = "[]"
	}
	return retval, nil
}