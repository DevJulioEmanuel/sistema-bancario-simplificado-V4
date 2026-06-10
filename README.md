# 🏦 Sistema Bancário Simplificado e Distribuído — V4

> **Disciplina:** Sistemas Distribuídos — UFC Quixadá

> **Autores:** Arthur Lelis Uchoa e Julio Emanuel Pereira da Silva

> **Arquitetura:** Comunicação Indireta via Filas de Mensagens (RabbitMQ)

> **Link do vídeo da apresentação** https://www.youtube.com/watch?v=jaooVOkIBq8 

---

## 📌 Visão Geral

Esta versão representa a **quarta evolução** do sistema bancário distribuído. A principal mudança arquitetural é a introdução de **comunicação indireta** via **RabbitMQ**, substituindo o processamento síncrono das operações financeiras por um modelo assíncrono baseado em filas de mensagens persistentes.

As operações financeiras (depósito, saque, transferência e pagamento) deixaram de ser processadas diretamente pela API e passaram a ser delegadas a um serviço independente — o **payment-service** — que se comunica com a **core-api** exclusivamente via filas, sem nenhum acoplamento direto entre os dois.


<div align="center">
  <img width="40%" alt="Diagrama de Arquitetura" src="https://github.com/user-attachments/assets/60384997-dc8d-4058-a541-b60e058ee899" />
</div>

---

## 🏗️ Arquitetura

### Serviços

| Serviço | Responsabilidade |
|---|---|
| **core-api** | Recebe requisições HTTP, valida dados, publica operações na fila e responde imediatamente ao cliente |
| **payment-service** | Consome mensagens da fila, processa as operações financeiras, atualiza saldos e publica o resultado de volta |
| **RabbitMQ** | Intermediário de mensagens — garante entrega, persistência e desacoplamento entre os serviços |

### Fluxo de uma Operação Financeira

```
1. Cliente envia POST /contas/1001/depositar
2. core-api valida a requisição
3. core-api publica mensagem na fila "transacoes_conta"
4. core-api responde 202 Accepted ao cliente imediatamente
5. payment-service consome a mensagem da fila (FIFO garantido)
6. payment-service processa e atualiza o saldo no banco
7. payment-service publica resultado na fila "notificacao"
8. core-api consome a notificação e notifica o cliente via WebSocket
```

### Decisão Arquitetural — Fila Única

Uma decisão importante nesta versão foi o uso de **uma única fila** (`transacoes_conta`) para todas as operações financeiras, em vez de uma fila por tipo de operação.

O motivo é a garantia de **ordenação FIFO**. O RabbitMQ garante FIFO dentro de uma fila, mas não entre filas distintas. Com filas separadas, um saque poderia ser processado antes de um depósito feito anteriormente — o que em um sistema bancário poderia gerar "saldo insuficiente" indevido. Com uma única fila, a ordem de chegada é sempre respeitada.

O tipo da operação é identificado pelo campo `tipoOperacao` dentro da mensagem, e o `payment-service` faz o dispatch correto via `switch`.

### Operações Síncronas (core-api processa diretamente)

- Cadastro de cliente
- Login
- Consulta de conta
- Extrato
- Cálculo de rendimento
- Consulta de notificações

### Operações Assíncronas (via RabbitMQ → payment-service)

- Depósito
- Saque
- Transferência
- Pagamento

---

## 🔔 Sistema de Notificações e WebSocket

Uma das decisões arquiteturais desta versão foi implementar um **canal de feedback assíncrono** entre o `payment-service` e o `core-api`, com entrega em tempo real ao cliente via **WebSocket**.

O problema: como o cliente recebe `202 Accepted` imediatamente, ele não sabe se a operação foi bem-sucedida ou falhou (por exemplo, por saldo insuficiente).

A solução adotada combina duas tecnologias:

```
payment-service processa → publica resultado na fila "notificacao"
                                          ↓
                     core-api consome via consumidor dedicado
                                          ↓
                     core-api envia via WebSocket para o cliente
```

**Fluxo do WebSocket:**
1. Cliente abre uma conexão WebSocket com a `core-api`
2. A conexão é registrada no `ClientesHub` associada ao número da conta
3. Quando o `payment-service` publica o resultado na fila `fila_resposta_pagamentos`, o consumidor do `core-api` recebe
4. O `core-api` localiza a conexão WebSocket do cliente no Hub e envia o evento `saldo_atualizado` com o novo saldo

Cada notificação contém:
- Tipo do evento (`saldo_atualizado`)
- Novo saldo após o processamento
- Status da operação (sucesso ou falha com motivo)

Isso mantém o desacoplamento temporal — o `payment-service` não conhece o cliente, apenas publica o resultado na fila.

---

## 🔑 Propriedades de Comunicação Indireta Demonstradas

### Desacoplamento Espacial
O `core-api` publica mensagens na fila sem conhecer o endereço, porta ou identidade do `payment-service`. Os dois serviços se comunicam exclusivamente através do RabbitMQ, sem nenhuma referência direta um ao outro.

### Desacoplamento Temporal
O `core-api` não precisa que o `payment-service` esteja online para publicar mensagens. As mensagens ficam persistidas nas filas até que o consumidor esteja disponível para processá-las — podendo ser minutos, horas ou dias depois.

### Resiliência a Falhas
Se o `payment-service` cair, as mensagens aguardam na fila. Quando o serviço retornar, todas as operações pendentes são processadas automaticamente — sem perda de dados, sem intervenção manual.

---

## 📄 Relatório Técnico

### Justificativa da Escolha — Filas de Mensagens (Opção C)

O sistema bancário simplificado já possuía uma API REST funcional com operações financeiras processadas de forma síncrona. O principal problema dessa arquitetura era o **acoplamento temporal**: a API ficava bloqueada aguardando o processamento de cada operação antes de responder ao cliente, e qualquer falha no processamento propagava erro diretamente para o usuário.

As **Filas de Mensagens** foram a abordagem escolhida por três razões centrais:

**1. Adequação ao domínio bancário.** Operações financeiras como depósitos, saques e transferências não precisam de resposta imediata sobre o resultado — o cliente precisa apenas saber que a operação foi *recebida*. Bancos reais funcionam assim: um Pix entra em "processando" antes de ser confirmado.

**2. Resiliência natural.** Com filas persistentes e duráveis, uma falha no `payment-service` não perde nenhuma operação — as mensagens aguardam na fila e são processadas quando o serviço retorna. Isso seria impossível com comunicação direta.

**3. Desacoplamento completo.** A `core-api` e o `payment-service` não se conhecem — comunicam-se exclusivamente pelo contrato das mensagens. Isso permite evoluir, substituir ou escalar cada serviço independentemente.

---

### Análise de Overhead e Complexidade

A introdução do RabbitMQ como intermediário trouxe ganhos arquiteturais significativos, mas também custos que precisam ser reconhecidos.

**Overhead de desempenho introduzido:**

- Cada operação financeira agora envolve serialização JSON, publicação na fila, persistência em disco pelo broker e deserialização pelo consumidor. O que antes era uma chamada de função direta passou a ter latência adicional de rede e I/O.
- A conexão com o RabbitMQ adiciona um ponto de falha de infraestrutura: se o broker cair, nenhuma operação financeira pode ser publicada.

**Como o overhead afetou o sistema:**

Na prática, a latência adicional foi imperceptível para o usuário — a `core-api` responde `202 Accepted` em menos de 1ms, enquanto o processamento pelo `payment-service` ocorre em paralelo em poucos milissegundos. O modelo assíncrono transformou uma operação bloqueante em não-bloqueante, melhorando a percepção de desempenho do cliente. O WebSocket eliminou a necessidade de polling — o cliente é notificado instantaneamente quando o processamento termina.

**Estratégias de mitigação:**

| Problema | Mitigação Adotada |
|---|---|
| Falha do broker | Retry automático na inicialização dos serviços |
| Perda de mensagens | Filas `durable: true` + mensagens `Persistent` |
| Processamento duplicado | ACK manual — mensagem só é removida após confirmação |
| Erros de negócio em loop | NACK sem requeue para erros permanentes (saldo insuficiente) |
| Ordenação entre operações | Fila única `transacoes_conta` com QoS=1 garantindo FIFO estrito |
| Latência de feedback | WebSocket para notificação em tempo real ao cliente |

Em produção, o overhead do broker seria mitigado com um **cluster RabbitMQ** com replicação de filas, eliminando o ponto único de falha. Para escala horizontal, múltiplas instâncias do `payment-service` poderiam consumir a mesma fila em paralelo, distribuindo a carga automaticamente — algo impossível na arquitetura síncrona original.

---

## 🗂️ Modelagem de Domínio

**Cliente** — Representa um usuário único no sistema.
- `Nome`, `CPF` *(identificador único)*, `Senha` *(hash bcrypt)*

**Conta** — Representa a conta bancária atrelada a um cliente.
- `Número`, `Saldo`, `Titular`, `Tipo`

**Transação** — Gerada a cada operação bancária.
- `Tipo`, `Descrição`, `Valor`, `Data`

**Notificação** — Resultado do processamento assíncrono.
- `Tipo`, `Status`, `Motivo`, `Data`

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
| **Rendimento** | Síncrono | Exclusivo Conta Poupança. Taxa: 0,5% a.m. |

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
| **SQLite** | Banco de dados |
| **GORM** | ORM para acesso ao banco |
| **WebSocket** | Notificações em tempo real |
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
│   │   ├── consumer/       # Consumidor da fila de resposta + WebSocket hub
│   │   ├── db/             # Conexão com banco
│   │   ├── dto/            # Objetos de transferência
│   │   ├── handler/        # Handlers HTTP e WebSocket
│   │   ├── middleware/     # JWT Middleware
│   │   ├── publisher/      # Publicador RabbitMQ
│   │   ├── repository/     # Repositórios
│   │   ├── routes/         # Configuração de rotas
│   │   └── service/        # Regras de negócio
│   └── Dockerfile
├── payment-service/        # Serviço de pagamento (Go)
│   ├── internal/
│   │   ├── processor/      # Processadores de transações
│   │   └── publisher/      # Publicador de respostas e notificações
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

### Demonstração de resiliência

```bash
# Derrubar o payment-service
docker compose stop payment-service

# Fazer operações — mensagens ficam na fila
# Verificar em http://localhost:15672 → Queues

# Subir novamente — processa tudo automaticamente
docker compose start payment-service

# Acompanhar logs em tempo real
docker compose logs -f payment-service
```

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
| GET | `/contas/:num/notificacoes` | Consultar notificações | ✅ |
| WS | `/ws/:num` | WebSocket — notificações em tempo real | ✅ |

---

## 🌐 Ecossistema de Clientes

**Cliente Java** (Spring Boot)
- Injeção de dependências
- Interface via console CLI

**Cliente Python 3**
- Interface Textual de Usuário (TUI)
- Mapeamento dinâmico de chaves JSON

---

## 💻 Como Executar os Clientes

### Cliente Python

**Pré-requisitos:** Python 3.x instalado

```bash
cd client-python
pip install -r requirements.txt
python main.py
```

---

### Cliente Java

**Pré-requisitos:** Java 21 e IntelliJ IDEA instalados

1. Abra a pasta `client-java` no IntelliJ IDEA
2. Aguarde o Maven baixar as dependências automaticamente
3. Localize e execute o arquivo `BancoCliente.java`

Ou via terminal com Maven:

```bash
cd client-java
./mvnw spring-boot:run
```

> Certifique-se de que o servidor (`docker compose up`) está rodando antes de iniciar os clientes.
