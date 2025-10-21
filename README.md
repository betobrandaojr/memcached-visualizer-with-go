# Memcached Visualizer With Go

Interface web para gerenciar instâncias do Memcached.

## Sobre

Esta aplicação fornece uma interface visual intuitiva para:

- Conectar a instâncias do Memcached
- Criar e atualizar dados (chave-valor)
- Buscar dados por chave específica
- Buscar múltiplas chaves simultaneamente
- Deletar dados do cache
- Validação de entrada e tratamento de erros

## Tecnologias Utilizadas

**Backend:**

- Go 1.21+
- Biblioteca `gomemcache` para conexão com Memcached
- Servidor HTTP nativo do Go

**Frontend:**

- HTML5 + CSS3 + JavaScript
- Design responsivo com dark mode
- Interface glassmorphism

## Como Executar

### Pré-requisitos

- Go 1.21 ou superior instalado
- Docker e Docker Compose (para executar Memcached)

### Executando o Memcached

1. Inicie o Memcached usando Docker Compose:

```bash
docker-compose up -d
```

2. Para parar o Memcached:

```bash
docker-compose down
```

3. Para verificar se está rodando:

```bash
docker ps
```

### Executando a Aplicação

1. Clone o repositório:

```bash
git clone <url-do-repositorio>
cd memcached-management
```

2. Instale as dependências:

```bash
go mod tidy
```

3. Execute a aplicação:

```bash
go run cmd/main.go
```

4. Acesse no navegador:

```
http://localhost:5000
```

## Como Usar

### 1. Conectar ao Memcached

- Digite a URL no formato `host:porta`
- **Exemplos válidos:**
  - `localhost:11211` (recomendado)
  - `memcached:11211` (será convertido automaticamente)
  - `127.0.0.1:11211`
- Clique em "Conectar"
- Aguarde confirmação da conexão

### 2. Operações CRUD

Após conectar, você terá acesso a:

**Listar Chaves:**

- Visualize todas as chaves armazenadas no cache
- Mostra o total de chaves encontradas

**Criar/Atualizar:**

- Chave: identificador único (ex: `user:123`)
- Valor: dados a serem armazenados (ex: `{"nome":"João"}`)

**Buscar:**

- Digite a chave específica para recuperar o valor
- Use vírgulas para buscar múltiplas chaves (ex: `user:1,user:2,config:timeout`)

**Deletar:**

- Digite a chave para remover do cache

## Exemplos de Uso

```
Chave: user:1
Valor: {"nome":"João","email":"joao@email.com"}

Chave: config:timeout
Valor: 30

Chave: product:123
Valor: {"id":123,"nome":"Notebook","preco":2500}
```

## Funcionalidades

- **Listar Chaves**: Utiliza comandos `stats items` e `stats cachedump` para extrair todas as chaves
- **CRUD Completo**: Criar, ler, atualizar e deletar dados
- **Busca Múltipla**: Buscar várias chaves simultaneamente
- **Limpeza Total**: Remover todos os dados do cache
- **Interface Responsiva**: Funciona em desktop e mobile

## Limitações

- A listagem de chaves usa comandos internos do Memcached (pode ser lenta em caches grandes)
- Recomenda-se usar prefixos organizados (ex: `user:`, `product:`)
- A funcionalidade de listar chaves pode não funcionar em algumas versões antigas do Memcached

## Estrutura do Projeto

```
memcached-management/
├── cmd/
│   └── main.go          # Ponto de entrada da aplicação
├── handlers/
│   └── handlers.go      # Handlers HTTP
├── models/
│   ├── types.go         # Estruturas de dados
│   └── types_test.go    # Testes dos modelos
├── services/
│   ├── memcached.go     # Lógica de negócio
│   └── memcached_test.go# Testes do serviço
├── tests/
│   ├── integration/     # Testes de integração
│   │   └── api_test.go
│   └── unit/           # Testes unitários específicos
├── web/
│   ├── index.html       # Interface web
│   └── assets/
│       ├── css/         # Estilos CSS
│       └── js/          # Scripts JavaScript
├── Makefile            # Comandos de build e teste
├── docker-compose.yml   # Configuração Memcached
├── go.mod              # Dependências Go
└── README.md           # Documentação
```

## Arquitetura

- **models/**: Estruturas de dados e tipos
- **services/**: Lógica de negócio e integração com Memcached
- **handlers/**: Manipuladores HTTP e validação de entrada
- **cmd/**: Ponto de entrada da aplicação
- **web/**: Interface web e assets estáticos (HTML, CSS, JS)

## Testes

A aplicação possui uma estrutura completa de testes:

### Executar todos os testes:
```bash
make test
# ou
go test -v ./...
```

### Executar apenas testes unitários:
```bash
make test-unit
# ou
go test -v ./services ./models
```

### Executar apenas testes de integração:
```bash
make test-integration
# ou
go test -v ./tests/integration
```

### Executar testes com cobertura:
```bash
make test-coverage
```

### Outros comandos úteis:
```bash
make build          # Compilar aplicação
make run            # Executar aplicação
make docker-up      # Iniciar Memcached
make docker-down    # Parar Memcached
make clean          # Limpar arquivos de build
```
