<div align="center">

# 🎬 Streaming Platform

**Plataforma de streaming com microserviços em Go e gRPC.**

[![Status](https://img.shields.io/badge/status-em%20desenvolvimento-orange?style=flat-square)](#)
[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev)
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

Monorepo de microserviços construído para explorar **gRPC na prática** com arquitetura limpa, comunicação assíncrona via NATS, transcodificação de vídeos com FFmpeg, autenticação via JWT e infraestrutura com Docker e Kubernetes.

---

## 🛠 Stack

| Tecnologia | Uso |
|---|---|
| **Go** | Linguagem principal |
| **gRPC / protobuf** | Comunicação entre serviços |
| **MongoDB** | Banco de dados por serviço |
| **JWT + Argon2id** | Autenticação e hash de senha |
| **NATS** | Mensageria assíncrona entre serviços |
| **FFmpeg** | Transcodificação de vídeos |
| **Docker Compose** | Ambiente de desenvolvimento local |
| **Kubernetes** | Orquestração em produção |
| **GitHub Actions** | CI/CD |

---

## 📦 Serviços

<!-- SERVICES_START -->
| Serviço | Porta | Descrição |
|---|---|---|
| **user-service** | `:50051` | Cadastro, login e validação de token JWT |
| **upload-service** | `:50052` | Upload de vídeos em chunks via client streaming gRPC |
| **transcode-service** | — | Transcodificação assíncrona de vídeos via NATS + FFmpeg |
<!-- SERVICES_END -->

---

## 🚀 Início Rápido

### Pré-requisitos

| Ferramenta | Uso |
|---|---|
| **Go 1.23+** | Linguagem principal |
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
git clone https://github.com/Turgho/go-streaming-platform.git
cd go-streaming-platform
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
# User Service - MongoDB
MONGO_USER_USER=root
MONGO_USER_PASSWORD=root
MONGO_USER_DB=user_db
MONGO_USER_URI=mongodb://root:root@user-db:27017/user_db?authSource=admin

# Upload Service - MongoDB
MONGO_UPLOAD_USER=root
MONGO_UPLOAD_PASSWORD=root
MONGO_UPLOAD_DB=upload_db
MONGO_UPLOAD_URI=mongodb://root:root@upload-db:27017/upload_db?authSource=admin

# JWT
JWT_SECRET=seu-segredo-aqui-minimo-32-caracteres

# NATS
NATS_URL=nats://nats:4222

# Comunicação interna entre serviços
INTERNAL_SERVICE_KEY=chave-interna-secreta
```

---

## 📦 Comandos disponíveis

<!-- MAKEFILE_START -->
```bash
make proto-user                        # Gera código Go do proto do user-service
make proto-upload                      # Gera código Go do proto do upload-service
make proto-userpb-upload               # Copia user proto para o upload-service
make proto-uploadpb-transcode          # Copia upload proto para o transcode-service
make proto-all                         # Gera todos os protos
make docker-up                         # Sobe todos os serviços e bancos
make docker-down                       # Para os containers
```
<!-- MAKEFILE_END -->

---

## 📁 Estrutura do projeto

<!-- TREE_START -->
```
go-streaming-platform/
├── proto/
│   ├── user/
│   │   └── user.proto
│   └── upload/
│       └── upload.proto
├── services/
│   ├── user/
│   │   ├── cmd/server/
│   │   ├── internal/
│   │   │   ├── domain/
│   │   │   ├── handler/grpc/
│   │   │   ├── infra/
│   │   │   ├── repository/
│   │   │   └── usecase/
│   │   └── pkg/
│   │       ├── hash/
│   │       └── pb/
│   ├── upload/
│   │   ├── cmd/
│   │   │   ├── client/
│   │   │   └── server/
│   │   ├── internal/
│   │   │   ├── domain/
│   │   │   ├── handler/grpc/
│   │   │   ├── infra/
│   │   │   ├── repository/
│   │   │   └── usecase/
│   │   └── pkg/
│   │       ├── pb/
│   │       └── userpb/
│   └── transcode/
│       ├── cmd/worker/
│       ├── internal/
│       │   ├── domain/
│       │   ├── handler/
│       │   ├── infra/
│       │   └── usecase/
│       └── pkg/uploadpb/
├── pkg/
│   └── events/
├── gateway/
├── k8s/
├── docker-compose.yml
└── Makefile
```
<!-- TREE_END -->

---

## 📄 Licença

Distribuído sob a licença **MIT**. Consulte o arquivo [LICENSE](./LICENSE) para mais detalhes.

---

<div align="center">
  Feito com ❤️ por <a href="https://github.com/Turgho">Turgho</a>
</div>