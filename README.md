# 🏦 Sistema Bancário Simplificado e Distribuído — V4

> **Disciplina:** Sistemas Distribuídos — UFC Quixadá

> **Autores:** Arthur Lelis Uchoa e Julio Emanuel Pereira da Silva

> **Arquitetura:** Microserviços com Comunicação Indireta via Filas de Mensagens

---

## 📌 Visão Geral

Esta versão representa a **quarta evolução** do sistema bancário distribuído. A principal mudança arquitetural desta versão é a introdução de **comunicação indireta** via **RabbitMQ**, substituindo o processamento síncrono das operações financeiras por um modelo assíncrono baseado em filas de mensagens.

As operações financeiras (depósito, saque, transferência e pagamento) deixaram de ser processadas diretamente pela API e passaram a ser delegadas a um serviço independente — o **payment-service** — por meio de filas persistentes no RabbitMQ.

```
  +-----------------------+      +-----------------------+
  |      Cliente Java     |      |     Cliente Python    |
  |  (Spring Boot / CLI)  |      |   (Python 3 / TUI)    |
  +-----------------------+      +-----------------------+
              \                              /
               \--- HTTP / JSON (Port 8080) /
                             |
                             v
               +---------------------------+
               |         core-api          |
               |  (Go API REST - Gin Engine)|
               +---------------------------+
                             |
                     publica na fila
                             |
                             v
               +---------------------------+
               |         RabbitMQ          |
               |   (Message Broker)        |
               +---------------------------+
                             |
                    consome da fila
                             |
                             v
               +---------------------------+
               |      payment-service      |
               |  (Go - Consumer Worker)   |
               +---------------------------+
                             |
                             v
               +---------------------------+
               |          SQLite           |
               |   (Banco de Dados         |
               |    Compartilhado)         |
               +---------------------------+
```

---

## 🏗️ Arquitetura

### Serviços

| Serviço | Responsabilidade |
|---|---|
| **core-api** | Recebe requisições HTTP, valida dados, publica eventos no RabbitMQ e responde imediatamente ao cliente |
| **payment-service** | Consome mensagens das filas, processa as operações financeiras e atualiza os saldos no banco de dados |
| **RabbitMQ** | Intermediário de mensagens — garante entrega, persistência e desacoplamento entre os serviços |

### Fluxo de uma Operação Financeira

```
1. Cliente envia POST /contas/1001/depositar
2. core-api valida a requisição
3. core-api publica mensagem na fila "deposito"
4. core-api responde 202 Accepted ao cliente imediatamente
5. payment-service consome a mensagem da fila
6. payment-service atualiza o saldo no SQLite
```

### Operações Síncronas (core-api processa diretamente)

- Cadastro de cliente
- Login
- Consulta de conta
- Extrato
- Cálculo de rendimento

### Operações Assíncronas (via RabbitMQ → payment-service)

- Depósito
- Saque
- Transferência
- Pagamento

---

## 🔑 Propriedades de Comunicação Indireta Demonstradas

### Desacoplamento Espacial
O `core-api` publica mensagens na fila sem conhecer o endereço ou identidade do `payment-service`. Os serviços se comunicam exclusivamente através do RabbitMQ.

### Desacoplamento Temporal
O `core-api` não precisa que o `payment-service` esteja online para publicar mensagens. As mensagens ficam persistidas nas filas até que o consumidor esteja disponível para processá-las.

### Resiliência a Falhas
Se o `payment-service` cair, as mensagens aguardam na fila. Quando o serviço retornar, todas as operações pendentes são processadas automaticamente — sem perda de dados.

---

## 🗂️ Modelagem de Domínio

```
Cliente
 └── Contas
      └── Histórico de Transações
```

**Cliente** — Representa um usuário único no sistema.
- `Nome`, `CPF` *(identificador único)*, `Senha` *(hash bcrypt)*

**Conta** — Representa a conta bancária atrelada a um cliente.
- `Número`, `Saldo`, `Titular`, `Tipo`

**Transação** — Gerada a cada operação bancária.
- `Tipo`, `Descrição`, `Valor`, `Data`

---

## 📋 Regras de Negócio

### Tipos de Conta

| Tipo | Numeração | Operações Exclusivas |
|---|---|---|
| **Conta Corrente** | A partir de `1001` | Pagamento |
| **Conta Poupança** | A partir de `5001` | Rendimento |

### Matriz de Operações

| Operação | Processamento | Validações |
|---|---|---|
| **Depósito** | Assíncrono | Valor positivo |
| **Saque** | Assíncrono | Saldo suficiente |
| **Transferência** | Assíncrono | Conta destino existe, saldo suficiente |
| **Pagamento** | Assíncrono | Exclusivo Conta Corrente, saldo suficiente |
| **Rendimento** | Síncrono | Exclusivo Conta Poupança. Taxa: **0,5% a.m.** |

---

## 🔐 Autenticação JWT

O sistema utiliza **JSON Web Tokens (JWT)** para autenticação. Após o login, o cliente recebe um token com validade de **15 minutos** que deve ser enviado no header de todas as requisições às rotas protegidas:

```
Authorization: Bearer <token>
```

---

## 🛠️ Stack Tecnológica

| Tecnologia | Uso |
|---|---|
| **Go 1.26** | Linguagem principal |
| **Gin** | Framework HTTP |
| **RabbitMQ 3.13** | Message Broker |
| **SQLite** | Banco de dados compartilhado |
| **GORM** | ORM para acesso ao banco |
| **JWT** | Autenticação |
| **bcrypt** | Hash de senhas |
| **Docker Compose** | Orquestração dos serviços |

---

## 📁 Estrutura do Projeto

```
/
├── core-api/               # API principal (Go + Gin)
│   ├── cmd/                # Ponto de entrada
│   ├── internal/
│   │   ├── auth/           # Chave JWT
│   │   ├── db/             # Conexão com banco
│   │   ├── dto/            # Objetos de transferência
│   │   ├── handler/        # Handlers HTTP
│   │   ├── middleware/     # JWT Middleware
│   │   ├── model/          # (legado)
│   │   ├── publisher/      # Publicador RabbitMQ
│   │   ├── repository/     # Repositórios
│   │   ├── routes/         # Configuração de rotas
│   │   └── service/        # Regras de negócio
│   └── Dockerfile
├── payment-service/        # Serviço de pagamento (Go)
│   ├── internal/
│   │   └── processor/      # Processadores de transações
│   ├── main.go
│   └── Dockerfile
├── shared/                 # Código compartilhado
│   └── events.go           # Structs e constantes das filas
├── client-java/            # Cliente Java (Spring Boot)
├── client-python/          # Cliente Python (TUI)
└── docker-compose.yml
```

---

## 🚀 Como Executar

### Pré-requisitos
- Docker
- Docker Compose

### Subir todos os serviços

```bash
docker compose up --build
```

A API ficará disponível em `http://localhost:8080`.
O RabbitMQ Management em `http://localhost:15672` (usuário: `guest`, senha: `guest`).

### Limpar o ambiente

```bash
docker compose down -v
```

---

## 🌐 Endpoints

### Clientes
| Método | Rota | Descrição | Auth |
|---|---|---|---|
| POST | `/clientes` | Cadastrar cliente | ❌ |
| POST | `/clientes/login` | Login | ❌ |

### Contas
| Método | Rota | Descrição | Auth |
|---|---|---|---|
| GET | `/contas/:num/` | Consultar conta | ✅ |
| POST | `/contas/:num/depositar` | Depositar | ✅ |
| POST | `/contas/:num/sacar` | Sacar | ✅ |
| POST | `/contas/:num/transferir` | Transferir | ✅ |
| POST | `/contas/:num/pagamento` | Pagar | ✅ |
| GET | `/contas/:num/extrato` | Extrato | ✅ |
| GET | `/contas/:num/rendimento/:meses` | Calcular rendimento | ✅ |

---

## 🌐 Ecossistema de Clientes

**Cliente Java** (Spring Boot)
- Injeção de dependências
- Interface via console CLI

**Cliente Python 3**
- Interface Textual de Usuário (TUI)
- Mapeamento dinâmico de chaves JSON

---

## 🎓 Conclusões

A quarta evolução do sistema introduziu comunicação indireta via filas de mensagens, demonstrando na prática os conceitos de **desacoplamento espacial e temporal** estudados na disciplina de Sistemas Distribuídos. A adoção do RabbitMQ como intermediário permitiu que o `core-api` respondesse imediatamente ao cliente sem aguardar o processamento das operações financeiras, aumentando a resiliência e a escalabilidade do sistema.
