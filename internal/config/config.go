package config

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	Server     ServerConfig
	Cloudinary CloudinaryConfig
	MongoDB    MongoDBConfig
	Redis      RedisConfig
	RabbitMQ   RabbitMQConfig
	JWT        JWTConfig
	Twilio     TwilioConfig
	Sendgrid   SendgridConfig
	Time       TimeConfig
}

type ServerConfig struct {
	AppName     string
	Environment string
	Port        string
}

type CloudinaryConfig struct {
	CloudName string
	APIKey    string
	APISecret string
}

type MongoDBConfig struct {
	URI         string
	Port        string
	Username    string
	Password    string
	Database    string
	Collections MongoDBCollectionConfig
}

type MongoDBCollectionConfig struct {
	Users    string
	Products string
	Orders   string
	Payments string
}

type RedisConfig struct {
	URI      string
	Port     string
	Password string
}

type RabbitMQConfig struct {
	URI      string
	Port     string
	Username string
	Password string

	// Queue names
	EmailVerificationQueue string
	PhoneVerificationQueue string
}

type JWTConfig struct {
	Secret            string
	Expiration        time.Duration
	RefreshExpiration time.Duration
	UserRole          string
	StoreRole         string
}

type TwilioConfig struct {
	AccountSID        string
	AuthToken         string
	MessageServiceSID string
	FromNumber        string
}

type SendgridConfig struct {
	APIKey                   string
	FromEmail                string
	VerificationTemplateID   string
	ForgotPasswordTemplateID string
}

type TimeConfig struct {
	EmailExpireTime          time.Duration
	PhoneExpireTime          time.Duration
	ForgotPasswordExpireTime time.Duration
}

func LoadConfig() *Config {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.SetConfigFile(".env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	// Set default values
	viper.SetDefault("ENVIRONMENT", "development")
	viper.SetDefault("PORT", "8080")

	return &Config{
		Server: ServerConfig{
			AppName:     viper.GetString("APP_NAME"),
			Environment: viper.GetString("ENVIRONMENT"),
			Port:        viper.GetString("PORT"),
		},
		Cloudinary: CloudinaryConfig{
			CloudName: viper.GetString("CLOUDINARY_CLOUD_NAME"),
			APIKey:    viper.GetString("CLOUDINARY_API_KEY"),
			APISecret: viper.GetString("CLOUDINARY_API_SECRET"),
		},
		MongoDB: MongoDBConfig{
			URI:      viper.GetString("MONGODB_URI"),
			Port:     viper.GetString("MONGODB_PORT"),
			Username: viper.GetString("MONGODB_USERNAME"),
			Password: viper.GetString("MONGODB_PASSWORD"),
			Database: viper.GetString("MONGODB_DATABASE"),
			Collections: MongoDBCollectionConfig{
				Users:    viper.GetString("MONGODB_COLLECTION_USERS"),
				Products: viper.GetString("MONGODB_COLLECTION_PRODUCTS"),
				Orders:   viper.GetString("MONGODB_COLLECTION_ORDERS"),
				Payments: viper.GetString("MONGODB_COLLECTION_PAYMENTS"),
			},
		},
		Redis: RedisConfig{
			URI:      viper.GetString("REDIS_URI"),
			Port:     viper.GetString("REDIS_PORT"),
			Password: viper.GetString("REDIS_PASSWORD"),
		},
		RabbitMQ: RabbitMQConfig{
			URI:      viper.GetString("RABBITMQ_URI"),
			Port:     viper.GetString("RABBITMQ_PORT"),
			Username: viper.GetString("RABBITMQ_USERNAME"),
			Password: viper.GetString("RABBITMQ_PASSWORD"),
			// Queue names
			EmailVerificationQueue: "email_verification",
			PhoneVerificationQueue: "phone_verification",
		},
		JWT: JWTConfig{
			Secret:            viper.GetString("JWT_SECRET"),
			Expiration:        viper.GetDuration("JWT_EXPIRES_IN"),
			RefreshExpiration: viper.GetDuration("JWT_REFRESH_EXPIRES_IN"),
		},
		Twilio: TwilioConfig{
			AccountSID:        viper.GetString("TWILIO_ACCOUNT_SID"),
			AuthToken:         viper.GetString("TWILIO_AUTH_TOKEN"),
			MessageServiceSID: viper.GetString("TWILIO_MESSAGE_SERVICE_SID"),
			FromNumber:        viper.GetString("TWILIO_FROM_NUMBER"),
		},
		Sendgrid: SendgridConfig{
			APIKey:                   viper.GetString("SENDGRID_API_KEY"),
			FromEmail:                viper.GetString("SENDGRID_FROM_EMAIL"),
			VerificationTemplateID:   viper.GetString("SENDGRID_VERIFICATION_EMAIL_TEMPLATE_ID"),
			ForgotPasswordTemplateID: viper.GetString("SENDGRID_FORGOT_PASSWORD_EMAIL_TEMPLATE_ID"),
		},
		Time: TimeConfig{
			EmailExpireTime:          viper.GetDuration("SENDGRID_EMAIL_EXPIRE_TIME"),
			PhoneExpireTime:          viper.GetDuration("SENDGRID_PHONE_EXPIRE_TIME"),
			ForgotPasswordExpireTime: viper.GetDuration("SENDGRID_FORGOT_PASSWORD_EXPIRE_TIME"),
		},
	}
}

var cfg *Config

func GetConfig() Config {
	if cfg == nil {
		cfg = LoadConfig()
	}

	return *cfg
}

func GetServerConfig() ServerConfig {
	return ServerConfig{
		AppName:     GetConfig().Server.AppName,
		Environment: GetConfig().Server.Environment,
		Port:        GetConfig().Server.Port,
	}
}

func GetCloudinaryConfig() CloudinaryConfig {
	return GetConfig().Cloudinary
}

func GetMongoDBConfig() MongoDBConfig {
	return GetConfig().MongoDB
}

func GetRedisConfig() RedisConfig {
	return GetConfig().Redis
}

func GetRabbitMQConfig() RabbitMQConfig {
	return GetConfig().RabbitMQ
}

func GetJWTConfig() JWTConfig {
	return GetConfig().JWT
}

func GetTwilioConfig() TwilioConfig {
	return GetConfig().Twilio
}

func GetSendgridConfig() SendgridConfig {
	return GetConfig().Sendgrid
}

func GetTimeConfig() TimeConfig {
	return GetConfig().Time
}
