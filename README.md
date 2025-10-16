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
go run main.go
```

4. Acesse no navegador:

```
http://localhost:5000
```

## Como Usar

### 1. Conectar ao Memcached

- Digite a URL no formato `host:porta` (ex: `localhost:11211`)
- Clique em "Conectar"
- Aguarde confirmação da conexão

### 2. Operações CRUD

Após conectar, você terá acesso a:

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

## Limitações

- Memcached não possui comando nativo para listar todas as chaves
- É necessário conhecer a chave específica para buscar dados
- Recomenda-se usar prefixos organizados (ex: `user:`, `product:`)

## Estrutura do Projeto

```
memcached-management/
├── main.go              # Backend em Go
├── main_test.go         # Testes unitários
├── templates/
│   └── index.html       # Interface web
├── docker-compose.yml   # Configuração Memcached
├── go.mod              # Dependências Go
└── README.md           # Documentação
```

## Testes

Execute os testes unitários:

```bash
go test -v
```
