package utils

import (
	"os"
	"strconv"
)

func CheckEnvStr(variable *string, envName string) {
	var tempString string = os.Getenv(envName)
	if tempString != "" {
		*variable = tempString
	}
}

func CheckEnvArrStr(variable *[]string, envName string) {
	var tempString string = os.Getenv(envName)
	if tempString != "" {
		*variable = []string{tempString}
	}
}

func CheckEnvInt(variable *int, envName string) {
	var tempString string = os.Getenv(envName)
	if tempString != "" {
		tempNum, err := strconv.Atoi(tempString)
		if err == nil {
			*variable = tempNum
		}
	}
}
