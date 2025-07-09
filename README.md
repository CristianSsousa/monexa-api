# My Finance Hub API

API backend para sistema de finanças pessoais desenvolvida em Go com Gin framework.

## Recursos

-   📊 Controle de receitas e despesas
-   📈 Estatísticas básicas de transações
-   🎯 Metas de poupança
-   📝 Categorização de transações
-   🗃️ Banco de dados Supabase

## Configuração

### Pré-requisitos

-   Go 1.21+
-   Banco de dados PostgreSQL (Supabase recomendado)

### Instalação

1. Clone o repositório e navegue até a pasta da API:

```bash
cd my-finance-hub-api
```

2. Instale as dependências:

```bash
go mod tidy
```

3. Configure as variáveis de ambiente criando um arquivo `.env` **(duas opções)**:

-   **Opção A** – usar a variável `DATABASE_URL` (string completa do Supabase/PostgreSQL):

```env
DATABASE_URL=postgresql://user:password@db.<HASH>.supabase.co:5432/postgres?sslmode=require
PORT=8080
GIN_MODE=release
```

-   **Opção B** – declarar variáveis individuais do Supabase (o backend montará o DSN automaticamente):

```env
SUPABASE_HOST=db.<HASH>.supabase.co
SUPABASE_PORT=5432        # opcional (padrão 5432)
SUPABASE_USER=postgres
SUPABASE_PASSWORD=<YOUR_PASSWORD>
SUPABASE_DB=postgres
PORT=8080
GIN_MODE=release
```

### Executar

```bash
go run main.go
```

A API estará disponível em `http://localhost:8080`

## Endpoints

### Transações

-   `GET /api/v1/transactions` - Listar transações
-   `POST /api/v1/transactions` - Criar transação
-   `GET /api/v1/transactions/:id` - Obter transação
-   `PUT /api/v1/transactions/:id` - Atualizar transação
-   `DELETE /api/v1/transactions/:id` - Excluir transação
-   `GET /api/v1/transactions/stats` - Estatísticas

### Categorias

-   `GET /api/v1/categories` - Listar categorias
-   `POST /api/v1/categories` - Criar categoria
-   `PUT /api/v1/categories/:id` - Atualizar categoria
-   `DELETE /api/v1/categories/:id` - Excluir categoria

### Metas de Poupança

-   `GET /api/v1/savings` - Listar metas
-   `POST /api/v1/savings` - Criar meta
-   `GET /api/v1/savings/:id` - Obter meta
-   `PUT /api/v1/savings/:id` - Atualizar meta
-   `DELETE /api/v1/savings/:id` - Excluir meta
-   `POST /api/v1/savings/:id/deposit` - Depositar em meta

### Health Check

-   `GET /health` - Verificar status da API

## Modelos de Dados

### Transação

```json
{
    "id": 1,
    "description": "Salário",
    "amount": 5000.0,
    "type": "income",
    "category_id": 1,
    "date": "2024-01-15T00:00:00Z"
}
```

### Categoria

```json
{
    "id": 1,
    "name": "Salário",
    "color": "#10B981",
    "icon": "💰",
    "type": "income"
}
```

### Meta de Poupança

```json
{
    "id": 1,
    "name": "Viagem",
    "description": "Viagem para o Japão",
    "target_amount": 10000.0,
    "current_amount": 3500.0,
    "target_date": "2024-12-31T00:00:00Z",
    "is_completed": false
}
```
