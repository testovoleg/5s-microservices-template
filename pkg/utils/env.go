package utils

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
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

func CheckEnvBool(variable *bool, envName string) {
	var tempString string = os.Getenv(envName)
	if tempString != "" {
		tempBool, err := strconv.ParseBool(tempString)
		if err == nil {
			*variable = tempBool
		}
	}
}

func CheckOTLName(variable *string, part, str string) {
	if variable == nil {
		return
	}

	str = strings.ToLower(str)
	str = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(str, "_")
	suffix := strings.ToLower(part)
	*variable = fmt.Sprintf("%s_%s_service", str, suffix)
}

func CheckKafkaGroup(variable *string, str string) {
	if variable == nil {
		return
	}

	str = strings.ToLower(str)
	str = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(str, "_")
	*variable = fmt.Sprintf("%s_microservice_consumer", str)
}

func CheckHttpTitle(variable *string, str string) {
	if variable == nil {
		return
	}

	str = regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(str, "-")
	*variable = formatServiceName(str)
}
