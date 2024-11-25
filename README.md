# 🎮 **Sera_462_API**

Backend do game **Sera 462** desenvolvido em Go **v1.22.5**

## 🚀 Tecnologias utilizadas
- **Go** `v1.22.5` ![Go](https://img.shields.io/badge/Go-1.22.5-blue)
- **MongoDB** `v6.0` 🍃
- **Redis Pub/Sub** ![Redis](https://img.shields.io/badge/Redis-6.0-FF0000) 🔄
- **RabbitMQ** ![RabbitMQ](https://img.shields.io/badge/RabbitMQ-3.8.9-FF6600) 🐇
- **Docker** 🐳

---

## 🛠️ **Como rodar localmente a aplicação**

### 1️⃣ Passo:  
Abra o terminal no Linux, acesse a pasta do projeto (**sera_462_api**) e execute:  
```bash
docker compose up -d

### 2️⃣ Passo:  
Acesse o diretório do código-fonte da API:  
```bash
cd cmd/api
export SRV_PORT=8080
export SRV_MODE=developer
export SRV_MDB_URI=mongodb://admin:supersenha@localhost:27017/
export SRV_MDB_NAME=sera_db
export SRV_MDB_DEFAULT_COLLECTION=cfSera
export SRV_JWT_SECRET_KEY=RgUkXp2s5v8yJavaLinux
export SRV_JWT_TOKEN_EXP=300
export PASSWRD_DB=LinuxDB@162
export SRV_ASAAS_API_KEY=KeyASSAS
export SRV_ASAAS_WALLET_ID=ID_ASSAS_WALLET
export SRV_ASAAS_URL_ASAAS=https://sandbox.asaas.com
export SRV_ASAAS_TIMEOUT=35
export AWS_ACCESS_KEY_ID=ID686868
export AWS_SECRET_ACCESS_KEY=SECREY
export AWS_REGION=us-east-1
export AWS_BUCKET_NAME=s3name
export SRV_RMQ_URI=amqp://admin:supersenha@localhost:5672/

### 3️⃣ Passo:  
Execute o seguinte comando para iniciar a aplicação:  
```bash
go run main.go
