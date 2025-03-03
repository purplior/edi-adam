package config

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/purplior/edi-adam/domain/shared/constant"
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

func DebugMode() bool {
	return os.Getenv("APP_DEBUG_MODE") == "1"
}

func Phase() constant.Phase {
	phase := os.Getenv("APP_PHASE")
	if len(phase) == 0 {
		return constant.Phase_Production
	}

	return convertPhaseStringToEnum(phase)
}

func Port() int {
	portStr := os.Getenv("APP_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 8080
	}

	return port
}

func SymmetricKey() string {
	return os.Getenv("APP_SYMMETRIC_KEY")
}

func JwtSecretKey() string {
	return os.Getenv("APP_JWT_SECRET_KEY")
}

func PostgreSQLDSN() string {
	host := os.Getenv("DB_POSTGRE_HOST")
	user := os.Getenv("DB_POSTGRE_USER")
	password := os.Getenv("DB_POSTGRE_PASSWORD")
	dbName := os.Getenv("DB_POSTGRE_DBNAME")
	dsn := "host=" + host + " user= " + user + " password=" + password + " dbname=" + dbName + " port=5432 sslmode=require TimeZone=Asia/Seoul"

	return dsn
}

func MongoURI() string {
	return os.Getenv("DB_MONGO_URI")
}

func MongoDatabaseName() string {
	return os.Getenv("DB_MONGO_DATABASE_NAME")
}

func OpenAIServiceAccountApiKey() string {
	return os.Getenv("OPENAI_SA_API_KEY")
}

func OpenAIOrganizationID() string {
	return os.Getenv("OPENAI_ORGANIZATION_ID")
}

func OpenAIProjectID() string {
	return os.Getenv("OPENAI_PROJECT_ID")
}

func CustomerVoiceEmail() string {
	return os.Getenv("CLIENT_CS_EMAIL")
}

func CustomerVoiceEmailPassword() string {
	return os.Getenv("CLIENT_CS_EMAIL_PASSWORD")
}

func SlackBotToken() string {
	return os.Getenv("CLIENT_SLACK_BOT_TOKEN")
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

func AWSBaseRegion() string {
	return os.Getenv("AWS_BASE_REGION")
}

func AWSDynamoDBEndpoint() string {
	return os.Getenv("AWS_DYNAMODB_ENDPOINT")
}

func IsEnvLoaded() bool {
	return isEnvLoaded
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
