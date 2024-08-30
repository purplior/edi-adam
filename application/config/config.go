package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/podossaem/root/domain/constants"
)

func Init() error {
	phase := Phase()

	fmt.Printf("# set phase: %s\n", convertPhaseEnumToDisplay((phase)))
	var envFilePath string
	switch phase {
	case constants.Phase_Local:
		envFilePath = ".env.local"
	case constants.Phase_Stage:
		envFilePath = ".env.stage"
	case constants.Phase_Production:
		envFilePath = ".env.prod"
	default:
		envFilePath = ".env"
	}

	if _, err := os.Stat(envFilePath); errors.Is(err, os.ErrNotExist) {
		return nil
	}

	err := godotenv.Load(envFilePath)
	if err == nil {
		isEnvLoaded = true
	}

	return err
}

func Phase() constants.Phase {
	phase := os.Getenv("APP_PHASE")
	if len(phase) == 0 {
		return constants.Phase_Production
	}

	return convertPhaseStringToEnum(phase)
}

func IsEnvLoaded() bool {
	return isEnvLoaded
}

func AppPort() int {
	portStr := os.Getenv("PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 8080
	}

	return port
}

func convertPhaseStringToEnum(phaseString string) constants.Phase {
	enm, is := phaseStringToEnumMap[phaseString]
	if !is {
		return constants.Phase_Invalid
	}

	return enm
}

func convertPhaseEnumToDisplay(phaseEnum constants.Phase) string {
	for k, v := range phaseStringToEnumMap {
		if v == phaseEnum {
			return k
		}
	}

	return ""
}

var (
	isEnvLoaded = false

	phaseStringToEnumMap = map[string]constants.Phase{
		"local": constants.Phase_Local,
		"stage": constants.Phase_Stage,
		"prod":  constants.Phase_Production,
	}
)
