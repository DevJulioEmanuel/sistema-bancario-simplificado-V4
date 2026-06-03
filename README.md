# 🏦 Sistema Bancário Simplificado e Distribuído — V4

> **Disciplina:** Sistemas Distribuídos
> **Autor:** Arthur Lelis Uchoa e Julio Emanuel Pereira da Silva
> **Arquitetura:** Microserviços / Cliente-Servidor (REST API Poliglota)

> **Link do vídeo:** https://youtu.be/dEfyc4pequA?si=qfTYXjjeT2UMt9UV
---

## 📌 Visão Geral

Este projeto consiste em um **Sistema Bancário Simplificado Distribuído**, desenvolvido como evolução das versões anteriores da disciplina. O sistema foi transformado em uma **API REST centralizada**, permitindo que múltiplos clientes independentes e poliglotas — construídos em diferentes linguagens e plataformas — se comuniquem e realizem operações concorrentes no mesmo núcleo bancário via requisições HTTP/JSON.

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
               |     SERVIDOR CENTRAL      |
               |  (Go API REST - Gin Engine)|
               +---------------------------+
```

---

## 🏗️ Arquitetura do Servidor (Go)

O servidor central foi desenvolvido em **Go**, utilizando o framework **Gin** para otimizar o roteamento, aplicação de middlewares e manipulação ágil de JSON. O projeto adota uma arquitetura rigorosamente dividida em camadas de responsabilidade:

```
Cliente HTTP  ──>  Handler  ──>  Service  ──>  Repository  ──>  Banco (Em Memória)
```

| Camada | Responsabilidade |
|---|---|
| **Handler** | Isola a camada HTTP, recebe requisições, valida dados de entrada, converte DTOs e retorna respostas JSON padronizadas. |
| **Service** | Centraliza todas as regras de negócio do domínio bancário (validações, depósitos, saques, pagamentos, transferências e cálculo de rendimentos). |
| **Repository** | Encapsula, isola e abstrai o acesso aos dados da aplicação. |
| **Banco** | Armazenamento em memória que gerencia as estruturas internas de Clientes e Contas, sem dependência de banco de dados externo. |

### 🛡️ Tratamento de Concorrência Distribuída

Como múltiplas operações podem ocorrer simultaneamente através de diferentes clientes disparando requisições, o servidor utiliza primitivas de sincronização via `sync.RWMutex` para proteger o acesso e a escrita nas contas. Isso evita **condições de corrida** (*race conditions*), como dois depósitos ou saques alterando o mesmo saldo ao mesmo tempo.

---

## 🗂️ Modelagem de Domínio

A estrutura interna do banco de dados em memória agrega o domínio em três entidades principais:

```
Banco
 ├── Clientes
 └── Contas
      └── Histórico de Transações
```

**Cliente** — Representa um usuário único no sistema.
- `Nome`, `CPF` *(identificador único)*, `Senha`

**Conta** — Representa a conta bancária atrelada a um cliente.
- `Número`, `Saldo`, `Titular`, `Tipo`, `Histórico de Transações`

**Transação** — Gerada a cada operação bancária.
- `Tipo`, `Descrição`, `Valor`, `Data`

---

## 📋 Regras de Negócio

### Restrições de Contas por Cliente

- Um cliente pode possuir no máximo **1 Conta Corrente** e **1 Conta Poupança**.
- É expressamente **proibido** que um cliente possua duas contas do mesmo tipo.
- No momento do **login**, o cliente fornece CPF, Senha e o Tipo de Conta desejado para determinar qual ambiente será acessado.

### Tipos de Conta e Suas Faixas

| Tipo | Numeração | Operações Exclusivas |
|---|---|---|
| **Conta Corrente** | A partir de `1000` | Pagamento |
| **Conta Poupança** | A partir de `5000` | Rendimento |

### Matriz de Operações Bancárias

| Operação | Descrição / Fluxo | Validações |
|---|---|---|
| **Depósito** | Adiciona saldo à conta e retorna o valor depositado e o saldo atualizado. | Apenas valores estritamente positivos. |
| **Saque** | Retira dinheiro da conta logada. | Requer saldo disponível suficiente. |
| **Transferência** | Envia recursos entre contas do ecossistema. | Valida existência da conta destino e saldo suficiente na origem. |
| **Pagamento** | Quita boletos com descrição textual e valor. | Exclusiva para Conta Corrente. Requer saldo disponível. |
| **Rendimento** | Projeta rendimentos futuros sob juros compostos. | Exclusiva para Conta Poupança. Taxa: **0,5% a.m.** (`0.005`). |

---

## 🌐 Ecossistema de Clientes (Poliglotismo)

A arquitetura REST baseada em JSON permitiu o desenvolvimento de clientes em ecossistemas completamente distintos:

**Cliente Java** (Spring Boot)
- Injeção de dependências
- Tratamento rígido de tipos de dados
- Interface via console CLI
- Consumo de endpoints através de HTTP templates

**Cliente Python 3**
- Interface Textual de Usuário (TUI) rica
- Controle reativo de sessão em memória local do runtime
- Mapeamento dinâmico de chaves JSON: converte o padrão `CamelCase` nativo do Go para o `snake_case` esperado pelas boas práticas do Python

---

## 🎓 Conclusões e Aprendizados

A terceira evolução do sistema (V3) consolidou os benefícios práticos do **desacoplamento arquitetural**. Ao isolar o estado do banco de dados e as travas de concorrência exclusivamente no servidor Go, eliminou-se a necessidade de lógica pesada nos terminais dos clientes.

Contanto que o **contrato da API REST** e a integridade dos payloads JSON sejam estritamente respeitados, diferentes linguagens de programação com paradigmas e tempos de execução opostos conseguem interoperar de forma **transparente, estável e segura** sobre a rede.

# Rodando o Servidor com Docker

## Build da imagem

```bash
docker build -t banco-api .
```

---

## Executar o container

```bash
docker run -p 8080:8080 banco-api
```

A API ficará disponível em:

```text
http://localhost:8080
```
