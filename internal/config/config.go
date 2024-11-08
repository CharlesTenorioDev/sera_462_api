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
	LlamaConfig   *LlamaConfig  `json:"llama_config"`
	GptConfig     *GptConfig    `json:"gpt_config"`
	GiminiConfig  *GiminiConfig `json:"gimini_config"`
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

type LlamaConfig struct {
	SRV_Llama_URL     string `json:"srv_llama_url"`
	SRV_Llama_API_KEY string `json:"srv_llama_api_key"`
}

type GptConfig struct {
	SRV_GPT_URL     string `json:"srv_gpt_url"`
	SRV_GPT_API_KEY string `json:"srv_gpt_api_key"`
	SRV_GPT_MODEL   string `json:"srv_gpt_model"`
}

type GiminiConfig struct {
	API_KEY string `json:"api_key"`
	URL     string `json:"url"`
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

	if llamaUrl := os.Getenv("SRV_Llama_URL"); llamaUrl != "" {
		conf.LlamaConfig.SRV_Llama_URL = llamaUrl
	}

	if llamaApiKey := os.Getenv("SRV_Llama_API_KEY"); llamaApiKey != "" {
		conf.LlamaConfig.SRV_Llama_API_KEY = llamaApiKey
	}

	if gptUrl := os.Getenv("SRV_GPT_URL"); gptUrl != "" {
		conf.GptConfig.SRV_GPT_URL = gptUrl
	}

	if gptApiKey := os.Getenv("SRV_GPT_API_KEY"); gptApiKey != "" {
		conf.GptConfig.SRV_GPT_API_KEY = gptApiKey
	}

	if gptModel := os.Getenv("SRV_GPT_MODEL"); gptModel != "" {
		conf.GptConfig.SRV_GPT_MODEL = gptModel
	}

	if giminiApiKey := os.Getenv("GEMINI_API_KEY"); giminiApiKey != "" {
		conf.GiminiConfig.API_KEY = giminiApiKey
	}

	if giminiUrl := os.Getenv("GIMINI_URL"); giminiUrl != "" {
		conf.GiminiConfig.URL = giminiUrl
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
		LlamaConfig: &LlamaConfig{
			SRV_Llama_URL: "https://api.llama-api.com",
		},
		GptConfig: &GptConfig{
			SRV_GPT_URL: "https://api.openai.com/v1/chat/completions",
		},
		GiminiConfig: &GiminiConfig{
			URL: "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=",
		},
	}

	defaultConf.TokenAuth = jwtauth.New("HS256", []byte(defaultConf.JWTSecretKey), nil)

	return defaultConf
}
