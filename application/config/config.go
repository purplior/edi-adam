package config

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/purplior/podoroot/domain/shared/constant"
)

func Init() error {
	phase := Phase()

	log.Printf("[#] PHASE: %s\n", convertPhaseEnumToDisplay((phase)))
	var envFilePath string
	switch phase {
	case constant.Phase_Local:
		envFilePath = ".env.local"
	case constant.Phase_Alpha:
		envFilePath = ".env.alpha"
	case constant.Phase_Production:
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

func Phase() constant.Phase {
	phase := os.Getenv("APP_PHASE")
	if len(phase) == 0 {
		return constant.Phase_Production
	}

	return convertPhaseStringToEnum(phase)
}

func JwtSecretKey() string {
	return os.Getenv("JWT_SECRET_KEY")
}

func SqlDbDSN() string {
	return os.Getenv("SQL_DB_DSN")
}

func MongoDbURI() string {
	return os.Getenv("MONGO_DB_URI")
}

func MongoDbName() string {
	return os.Getenv("MONGO_DB_NAME")
}

func OpenAiServiceAccountApiKey() string {
	return os.Getenv("OPENAI_SA_API_KEY")
}

func OpenAiOrganizationID() string {
	return os.Getenv("OPENAI_ORGANIZATION_ID")
}

func OpenAiProjectID() string {
	return os.Getenv("OPENAI_PROJECT_ID")
}

func CsEmail() string {
	return os.Getenv("CS_EMAIL")
}

func CsEmailPassword() string {
	return os.Getenv("CS_EMAIL_PASSWORD")
}

func SlackBotToken() string {
	return os.Getenv("SLACK_BOT_TOKEN")
}

func NCloudSMSServiceID() string {
	return os.Getenv("NCLOUD_SMS_SERVICE_ID")
}

func NCloudSMSFrom() string {
	return os.Getenv("NCLOUD_SMS_FROM")
}

func NCloudAccessKey() string {
	return os.Getenv("NCLOUD_ACCESS_KEY")
}

func NCloudSecretKey() string {
	return os.Getenv("NCLOUD_SECRET_KEY")
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

func DebugMode() bool {
	return os.Getenv("DEBUG_MODE") == "true"
}

func convertPhaseStringToEnum(phaseString string) constant.Phase {
	enm, is := phaseStringToEnumMap[phaseString]
	if !is {
		return constant.Phase_Invalid
	}

	return enm
}

func convertPhaseEnumToDisplay(phaseEnum constant.Phase) string {
	for k, v := range phaseStringToEnumMap {
		if v == phaseEnum {
			return k
		}
	}

	return ""
}

var (
	isEnvLoaded = false

	phaseStringToEnumMap = map[string]constant.Phase{
		"local": constant.Phase_Local,
		"alpha": constant.Phase_Alpha,
		"prod":  constant.Phase_Production,
	}
)
