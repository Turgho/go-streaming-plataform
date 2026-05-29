<div align="center">

# 🎬 Streaming Platform

**Plataforma de streaming com microserviços em Go e gRPC.**

[![Status](https://img.shields.io/badge/status-em%20desenvolvimento-orange?style=flat-square)](#)
[![Go](https://img.shields.io/badge/Go-1.26.3+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev)
[![gRPC](https://img.shields.io/badge/gRPC-protobuf-244c5a?style=flat-square&logo=google&logoColor=white)](https://grpc.io)
[![MongoDB](https://img.shields.io/badge/MongoDB-Docker-47a248?style=flat-square&logo=mongodb&logoColor=white)](https://www.mongodb.com)
[![License](https://img.shields.io/badge/license-MIT-green?style=flat-square)](./LICENSE)

> Projeto de aprendizagem focado em gRPC, Clean Architecture, Kubernetes e CI/CD com Go.

</div>

---

## 📋 Sumário

- [Sobre](#-sobre)
- [Stack](#-stack)
- [Serviços](#-serviços)
- [Início Rápido](#-início-rápido)
- [Configuração](#-configuração)
- [Comandos disponíveis](#-comandos-disponíveis)
- [Estrutura do projeto](#-estrutura-do-projeto)
- [Licença](#-licença)

---

## 📌 Sobre

Monorepo de microserviços construído para explorar **gRPC na prática** com arquitetura limpa, comunicação entre serviços, autenticação via JWT e infraestrutura com Docker e Kubernetes.

---

## 🛠 Stack

| Tecnologia | Uso |
|---|---|
| **Go** | Linguagem principal |
| **gRPC / protobuf** | Comunicação entre serviços |
| **MongoDB** | Banco de dados por serviço |
| **JWT + Argon2id** | Autenticação e hash de senha |
| **Docker Compose** | Ambiente de desenvolvimento local |
| **Kubernetes** | Orquestração em produção |
| **GitHub Actions** | CI/CD |

---

## 📦 Serviços

| Serviço | Porta | Descrição |
|---|---|---|
| **user-service** | `:50051` | Cadastro, login e validação de token |
| **upload-service** | `:50052` | Upload de arquivos em chunks via client streaming |

---

## 🚀 Início Rápido

### Pré-requisitos

| Ferramenta | Uso |
|---|---|
| **Go 1.26.3+** | Linguagem principal |
| **Docker + Docker Compose** | Containers |
| **protoc** | Compilador de `.proto` |
| **protoc-gen-go** | Plugin Go para o protoc |
| **protoc-gen-go-grpc** | Plugin gRPC para o protoc |
| **grpcurl** | Teste dos serviços via terminal |

### Instalando os plugins do protoc

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### 1. Clone o repositório

```bash
git clone https://github.com/Turgho/streaming-platform.git
cd streaming-platform
```

### 2. Configure as variáveis de ambiente

```bash
cp .env.example .env
```

### 3. Suba o ambiente

```bash
make docker-up
```

---

## ⚙️ Configuração

Crie o `.env` na raiz com o seguinte conteúdo:

```env
# MongoDB
MONGO_USER=root
MONGO_PASSWORD=root
MONGO_DB=user_db
MONGO_URI=mongodb://root:root@user-db:27017/user_db?authSource=admin

# JWT
JWT_SECRET=seu-segredo-aqui-minimo-32-caracteres
```

---

## 📦 Comandos disponíveis

```bash
make docker-up        # Sobe todos os serviços e bancos
make docker-down      # Para os containers

make proto-user       # Gera código Go a partir do proto do user-service
make proto-upload     # Gera código Go a partir do proto do upload-service
```

---

## 📁 Estrutura do projeto

```
streaming-platform/
├── proto/
│   ├── user/
│   │   └── user.proto
│   └── upload/
│       └── upload.proto
│
├── services/
│   ├── user/
│   │   ├── cmd/server/
│   │   ├── internal/
│   │   │   ├── domain/
│   │   │   ├── infra/
│   │   │   ├── repository/
│   │   │   └── server/
│   │   ├── pkg/
│   │   │   ├── hash/
│   │   │   └── pb/
│   │   └── Dockerfile
│   │
│   └── upload/
│       └── ...
│
├── gateway/
├── k8s/
├── docker-compose.yml
└── Makefile
```

---

## 📄 Licença

Distribuído sob a licença **MIT**. Consulte o arquivo [LICENSE](./LICENSE) para mais detalhes.

---

<div align="center">
  Feito com ❤️ por <a href="https://github.com/Turgho">Turgho</a>
</div>