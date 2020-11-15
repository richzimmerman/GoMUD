package config

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var alphaStart = regexp.MustCompile("^[a-zA-Z]")
var Configuration map[string]string

func LoadConfiguration(filename string) error {
	Configuration = make(map[string]string)

	if len(filename) == 0 {
		return nil
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// If line doesn't start with text, assume it's comments or invalid and skip the line
		if !alphaStart.Match([]byte(line)) {
			continue
		}
		if equal := strings.Index(line, "="); equal >= 0 {
			if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
				value := ""
				if len(line) > equal {
					value = strings.TrimSpace(line[equal+1:])
				}
				Configuration[key] = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func GetValue(key string) (string, error) {
	if value, found := Configuration[key]; found {
		return value, nil
	}
	return "", fmt.Errorf("config value (%s) not specified", key)
}

func GetStringValue(key string) (string, error) {
	return GetValue(key)
}

func GetIntegerValue(key string) (int, error) {
	s, err := GetValue(key)
	if err != nil {
		return 0, err
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("value (%s) is not an integer type for key (%s)", s, key)
	}
	return v, nil
}
