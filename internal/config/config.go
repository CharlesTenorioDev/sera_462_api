package config

import (
	"os"
	"strconv"
	"time"

	"github.com/go-chi/jwtauth"
)

const (
	DEVELOPER    = "developer"
	HOMOLOGATION = "homologation"
	PRODUCTION   = "production"
)

type Config struct {
	PORT          string `json:"port"`
	Mode          string `json:"mode"`
	MongoDBConfig `json:"mongo_config"`
	JWTSecretKey  string `json:"jwt_secret_key"`
	JWTTokenExp   int    `json:"jwt_token_exp"`
	TokenAuth     *jwtauth.JWTAuth
	DataInicial   time.Time
	DataFinal     time.Time
	AsaasConfig   `json:"asaas_config"`

	*AWS_CONFIG
	*RMQConfig
}
type MongoDBConfig struct {
	MDB_URI                string `json:"mdb_uri"`
	MDB_NAME               string `json:"mdb_name"`
	MDB_DEFAULT_COLLECTION string `json:"mdb_default_collection"`
}

type AsaasConfig struct {
	URL_ASAAS       string `json:"url_asaas"`
	ASAAS_API_KEY   string `json:"asas_api_key"`
	ASAAS_WALLET_ID string `json:"asas_wallet_id"`
	ASAAS_TIMEOUT   int    `json:"asas_timeout"`
}

type AWS_CONFIG struct {
	ACCESS_KEY_ID     string `json:"access_key_id"`
	SECRET_ACCESS_KEY string `json:"secret_access_key"`
	REGION            string `json:"region"`
	BUCKET_NAME       string `json:"bucket_name"`
}

type RMQConfig struct {
	RMQ_URI                  string `json:"rmq_uri"`
	RMQ_MAXX_RECONNECT_TIMES int    `json:"rmq_maxx_reconnect_times"`
}

type ConsumerConfig struct {
	ExchangeName  string `json:"exchange_name"`
	ExchangeType  string `json:"exchange_type"`
	RoutingKey    string `json:"routing_key"`
	QueueName     string `json:"queue_name"`
	ConsumerName  string `json:"consumer_name"`
	ConsumerCount int    `json:"consumer_count"`
	PrefetchCount int    `json:"prefetch_count"`
	Reconnect     struct {
		MaxAttempt int `json:"max_attempt"`
		Interval   int `json:"interval"`
	}
}

func NewConfig() *Config {
	conf := defaultConf()

	if port := os.Getenv("SRV_PORT"); port != "" {
		conf.PORT = port
	}

	if mode := os.Getenv("SRV_MODE"); mode != "" {
		conf.Mode = mode
	}

	if uri := os.Getenv("SRV_MDB_URI"); uri != "" {
		conf.MDB_URI = uri
	}

	if name := os.Getenv("SRV_MDB_NAME"); name != "" {
		conf.MDB_NAME = name
	}

	if collection := os.Getenv("SRV_MDB_DEFAULT_COLLECTION"); collection != "" {
		conf.MDB_DEFAULT_COLLECTION = collection
	}

	if secretKey := os.Getenv("SRV_JWT_SECRET_KEY"); secretKey != "" {
		conf.JWTSecretKey = secretKey
	}

	if tokenExp := os.Getenv("SRV_JWT_TOKEN_EXP"); tokenExp != "" {
		conf.JWTTokenExp, _ = strconv.Atoi(tokenExp)
	}

	if urlAsaas := os.Getenv("SRV_ASAAS_URL_ASAAS"); urlAsaas != "" {
		conf.AsaasConfig.URL_ASAAS = urlAsaas
	}

	if apiKey := os.Getenv("SRV_ASAAS_API_KEY"); apiKey != "" {
		conf.AsaasConfig.ASAAS_API_KEY = apiKey
	}

	if walletId := os.Getenv("SRV_ASAAS_WALLET_ID"); walletId != "" {
		conf.AsaasConfig.ASAAS_WALLET_ID = walletId
	}

	if timeout := os.Getenv("SRV_ASAAS_TIMEOUT"); timeout != "" {
		conf.AsaasConfig.ASAAS_TIMEOUT, _ = strconv.Atoi(timeout)
	}

	if SRV_ACCESS_KEY_ID := os.Getenv("AWS_ACCESS_KEY_ID"); SRV_ACCESS_KEY_ID != "" {
		conf.AWS_CONFIG.ACCESS_KEY_ID = SRV_ACCESS_KEY_ID
	}

	if SRV_SECRET_ACCESS_KEY := os.Getenv("AWS_SECRET_ACCESS_KEY"); SRV_SECRET_ACCESS_KEY != "" {
		conf.AWS_CONFIG.SECRET_ACCESS_KEY = SRV_SECRET_ACCESS_KEY
	}

	if SRV_REGION := os.Getenv("AWS_REGION"); SRV_REGION != "" {
		conf.AWS_CONFIG.REGION = SRV_REGION
	}

	if SRV_BUCKET_NAME := os.Getenv("AWS_BUCKET_NAME"); SRV_BUCKET_NAME != "" {
		conf.AWS_CONFIG.BUCKET_NAME = SRV_BUCKET_NAME
	}

	SRV_RMQ_URI := os.Getenv("SRV_RMQ_URI")
	if SRV_RMQ_URI != "" {
		conf.RMQConfig.RMQ_URI = SRV_RMQ_URI

	}

	return conf
}

func defaultConf() *Config {
	defaultConf := &Config{
		PORT:         "8080",
		Mode:         DEVELOPER,
		JWTSecretKey: "RgUkXp2s5v8y/B?EH+KbPeShVmYq3t6",
		JWTTokenExp:  300,
		MongoDBConfig: MongoDBConfig{
			MDB_DEFAULT_COLLECTION: "cfSera",
		},
		AsaasConfig: AsaasConfig{
			URL_ASAAS: "https://sandbox.asaas.com/api/",
		},

		AWS_CONFIG: &AWS_CONFIG{
			BUCKET_NAME: "",
		},

		RMQConfig: &RMQConfig{

			RMQ_MAXX_RECONNECT_TIMES: 3,
		},
	}

	defaultConf.TokenAuth = jwtauth.New("HS256", []byte(defaultConf.JWTSecretKey), nil)

	return defaultConf
}
