# 🎮 **Sera_462_API**

Backend do game **Sera 462** desenvolvido em Go **v1.22.5**

## 🚀 Tecnologias utilizadas
- **Go** `v1.22.5` ![Go](https://img.shields.io/badge/Go-1.22.5-blue)
- **MongoDB** `v6.0` 🍃
- **I.A. Llama** `v3.2` 🦙
- **AWS SNS** ☁️
- **AWS SQS** 📨
- **Docker** 🐳

---

## 🛠️ **Como rodar Llama localmente no Linux (Ubuntu)**

Para configurar e rodar o modelo de inteligência artificial **Llama** na sua máquina Linux, siga os passos abaixo:

1. **Liberar a porta 11434**:
   ```bash
   sudo ufw allow 11434/tcp
   curl -fsSL https://ollama.com/install.sh | sh
   ollama pull llama3.2:1b
   ollama run llama3.2
   roda localsta
    docker run -d -v /var/run/docker.sock:/var/run/docker.sock localstack/localstack
    serviso de SQS
    docker run --rm -it -e SERVICES=sqs -p 4566:4566 localstack/localstack
docker start localstack
