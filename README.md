
# Currency Exchange API

Esse repositório contém duas aplicações Go: uma API de conversão de moeda (`server.go`) e um cliente que consome essa API e salva o valor do dólar em um arquivo (`client.go`).

## Tecnologias Utilizadas

- **Go**: Linguagem de programação utilizada para construir a API e o cliente.
- **GORM**: ORM (Object-Relational Mapping) para realizar operações com o banco de dados SQLite.
- **SQLite**: Banco de dados leve, utilizado para armazenar a cotação do dólar.
- **Context**: Utilizado para controlar o tempo de execução de operações assíncronas, implementando timeout nas requisições HTTP e operações no banco de dados.

## Estrutura do Projeto

- **server.go**: Cria uma API que consulta a taxa de câmbio do dólar (USD-BRL) em uma API pública e armazena as informações no banco de dados SQLite.
- **client.go**: Realiza uma requisição para a API (`server.go`) para obter o valor atual do dólar e salva o valor em um arquivo de texto (`cotacao.txt`).

## Pré-requisitos

Para executar as aplicações, é necessário ter o Go instalado na máquina. A instalação pode ser feita conforme as instruções no [site oficial do Go](https://golang.org/dl/).

## Instruções de Execução

### Passo 1: Executar o Servidor

1. Clone o repositório e navegue até o diretório do projeto.
   
   ```bash
   git clone https://github.com/brunoofgod/goexpert-lesson-1.git
   cd goexpert-lesson-1
   ```

2. Compile e execute o servidor:
   ```bash
   cd server
   go build -o server server.go 
   ./server
   ```

### Passo 2: Executar o Cliente

1. Em um novo terminal, no mesmo diretório, compile e execute o cliente:
   ```bash
   cd client
   go build -o client client.go
   ./client
   ```
