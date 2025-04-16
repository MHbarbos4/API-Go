# API Go - CRUD de Itens com SQLite

Este projeto implementa uma API simples em Go que permite realizar operações CRUD (Criar, Ler, Atualizar e Deletar) em um banco de dados SQLite. Ele gerencia itens com informações de nome e preço. A API também está integrada com Swagger para facilitar a documentação e testes das rotas.

## Tecnologias Utilizadas

- Go
- SQLite
- Gorilla Mux (Roteamento)
- Swagger para documentação da API

## Funcionalidades

- **GET /items**: Recupera todos os itens no banco de dados.
- **GET /items/{id}**: Recupera um item específico pelo ID.
- **POST /items**: Cria um novo item no banco de dados.
- **PUT /items/{id}**: Atualiza um item existente pelo ID.
- **DELETE /items/{id}**: Deleta um item pelo ID.

## Instalação

### Requisitos

- Go 1.16 ou superior
- SQLite
- Swagger (gerado automaticamente pelo `swaggo`)

### Passos

1. **Clone o repositório**

   ```bash
   git clone https://github.com/MHbarbos4/API-Go.git
   cd API-Go

2. **Instalar dependências**

  O projeto utiliza o pacote swaggo para gerar a documentação do Swagger, além do gorilla/mux para roteamento.

```bash
   go mod tidy
```
3. **Compilar e Rodar o servidor**

  Para rodar o servidor, execute o seguinte comando:
```bash
chmod +x init.sh
source init.sh
```

4. **Acessar o Swagger**

  Após rodar o servidor, você pode acessar a interface do Swagger para testar as rotas da API. Abra o navegador e vá para:

```bash
http://localhost:8080/swagger/index.html
```
