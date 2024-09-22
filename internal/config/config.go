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

func NewConfig() *Config {
	conf := defaultConf()

	SRV_PORT := os.Getenv("SRV_PORT")
	if SRV_PORT != "" {
		conf.PORT = SRV_PORT
	}

	SRV_MODE := os.Getenv("SRV_MODE")
	if SRV_MODE != "" {
		conf.Mode = SRV_MODE
	}

	SRV_MDB_URI := os.Getenv("SRV_MDB_URI")
	if SRV_MDB_URI != "" {
		conf.MDB_URI = SRV_MDB_URI
	}

	SRV_MDB_NAME := os.Getenv("SRV_MDB_NAME")
	if SRV_MDB_NAME != "" {
		conf.MDB_NAME = SRV_MDB_NAME
	}

	SRV_MDB_DEFAULT_COLLECTION := os.Getenv("SRV_MDB_DEFAULT_COLLECTION")
	if SRV_MDB_DEFAULT_COLLECTION != "" {
		conf.MDB_DEFAULT_COLLECTION = SRV_MDB_DEFAULT_COLLECTION
	}

	SRV_JWT_SECRET_KEY := os.Getenv("SRV_JWT_SECRET_KEY")
	if SRV_JWT_SECRET_KEY != "" {
		conf.JWTSecretKey = SRV_JWT_SECRET_KEY
	}

	SRV_JWT_TOKEN_EXP := os.Getenv("SRV_JWT_TOKEN_EXP")
	if SRV_JWT_SECRET_KEY != "" {
		conf.JWTTokenExp, _ = strconv.Atoi(SRV_JWT_TOKEN_EXP)
	}

	SRV_ASAAS_URL_ASAAS := os.Getenv("SRV_ASAAS_URL_ASAAS")
	if SRV_ASAAS_URL_ASAAS != "" {
		conf.AsaasConfig.URL_ASAAS = SRV_ASAAS_URL_ASAAS
	}

	SRV_ASAAS_API_KEY := os.Getenv("SRV_ASAAS_API_KEY")
	if SRV_ASAAS_API_KEY != "" {
		conf.AsaasConfig.ASAAS_API_KEY = SRV_ASAAS_API_KEY
	}

	SRV_ASAAS_WALLET_ID := os.Getenv("SRV_ASAAS_WALLET_ID")
	if SRV_ASAAS_WALLET_ID != "" {
		conf.AsaasConfig.ASAAS_WALLET_ID = SRV_ASAAS_WALLET_ID
	}

	SRV_ASAAS_TIMEOUT := os.Getenv("SRV_ASAAS_TIMEOUT")
	if SRV_ASAAS_TIMEOUT != "" {
		conf.AsaasConfig.ASAAS_TIMEOUT, _ = strconv.Atoi(SRV_ASAAS_TIMEOUT)
	}

	return conf
}

func defaultConf() *Config {

	default_conf := Config{
		PORT:         "8080",
		Mode:         DEVELOPER,
		JWTSecretKey: "RgUkXp2s5v8y/B?E(H+KbPeShVmYq3t6", // "----your-256-bit-secret-here----" length 32
		JWTTokenExp:  300,
		// 15m
		MongoDBConfig: MongoDBConfig{
			MDB_DEFAULT_COLLECTION: "cfSera",
		},

		AsaasConfig: AsaasConfig{
			URL_ASAAS: "https://sandbox.asaas.com/api/",
		},
	}
	// Adicione as coleções padrão ao mapa MDB_COLLECTIONS

	default_conf.TokenAuth = jwtauth.New("HS256", []byte(default_conf.JWTSecretKey), nil)

	return &default_conf
}
