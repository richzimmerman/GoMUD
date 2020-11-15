package directions

import (
	"fmt"
	"strings"
)

const (
	directionFirstPerson = "You go %s"
	directionThirdPerson = "<A.NAME> leaves to the %s"
)

func firstPersonMsg(dirName string) string {
	return fmt.Sprintf(directionFirstPerson, strings.ToLower(dirName))
}

func thirdPersonMsg(dirName string) string {
	return fmt.Sprintf(directionThirdPerson, strings.ToLower(dirName))
}
