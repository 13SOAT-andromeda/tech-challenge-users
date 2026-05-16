# tech-challenge-users

Microserviço de gestão de identidade e perfil extraído do sistema legado `tech-challenge-s1`. Expõe todos os recursos como sub-rotas de `/users`, roteados via API Gateway v2 com uma única rota `ANY /users/{proxy+}`.

## Tecnologias

- **Go 1.23**
- **Gin** — HTTP router
- **GORM + pgx** — ORM com PostgreSQL
- **bcrypt** (custo 10) — hash de senhas
- **Docker Compose** — ambiente local

## Arquitetura

Hexagonal (ports & adapters):

```
cmd/api/                          → entrypoint
internal/
  domain/                         → entidades e value objects
  application/
    ports/                        → interfaces de repositório
    services/                     → lógica de negócio (serviços)
    usecases/                     → orquestração cross-aggregate
  adapter/
    config/                       → variáveis de ambiente
    database/
      migrations/                 → SQL embutido via go:embed
      model/                      → models GORM por entidade
      repository/                 → implementações dos ports
    http/
      handlers/                   → handlers Gin por recurso
      middlewares/                → AuthRequired, RoleRequired
pkg/
  converters/                     → domínio ↔ model
  encryption/                     → bcrypt hash/compare
```

## Rodando localmente

**Pré-requisitos:** Docker e Docker Compose instalados.

```bash
docker compose up --build -d
```

A API sobe em `http://localhost:8081`.

Para resetar o banco de dados:

```bash
docker compose down -v && docker compose up --build -d
```

## Variáveis de ambiente

| Variável         | Descrição                          | Padrão     |
|------------------|------------------------------------|------------|
| `DB_HOST`        | Host do PostgreSQL                 | `localhost` |
| `DB_PORT`        | Porta do PostgreSQL                | `5432`     |
| `DB_USER`        | Usuário do banco                   | —          |
| `DB_PASSWORD`    | Senha do banco                     | —          |
| `DB_NAME`        | Nome do banco                      | —          |
| `HTTP_PORT`      | Porta HTTP da aplicação            | `8081`     |
| `ENV`            | Ambiente (`local`, `production`)   | `local`    |
| `ADMIN_EMAIL`    | E-mail do admin inicial            | —          |
| `ADMIN_PASSWORD` | Senha do admin inicial             | —          |
| `ADMIN_DOCUMENT` | CPF do admin inicial               | —          |

## Autenticação

As rotas sob `/users` exigem os headers injetados pelo API Gateway após validação do JWT:

| Header          | Descrição                                      |
|-----------------|------------------------------------------------|
| `X-User-Id`     | ID do usuário autenticado (obrigatório)        |
| `X-User-Email`  | E-mail do usuário autenticado                  |
| `X-User-Role`   | Role: `administrator`, `attendant`, `mechanic`, `customer` |

## Rotas

### Usuários — `/users`

| Método   | Rota          | Role mínimo          | Descrição               |
|----------|---------------|----------------------|-------------------------|
| `POST`   | `/users`      | attendant            | Criar usuário/funcionário |
| `GET`    | `/users`      | attendant            | Listar usuários          |
| `GET`    | `/users/:id`  | attendant            | Buscar usuário por ID    |
| `PUT`    | `/users/:id`  | attendant            | Atualizar usuário        |
| `DELETE` | `/users/:id`  | attendant            | Remover usuário          |

### Clientes — `/users/customers`

| Método   | Rota                                      | Role mínimo | Descrição                        |
|----------|-------------------------------------------|-------------|----------------------------------|
| `POST`   | `/users/customers`                        | attendant   | Criar cliente                    |
| `GET`    | `/users/customers`                        | attendant   | Listar clientes (filtros opcionais) |
| `GET`    | `/users/customers/:id`                    | attendant   | Buscar cliente por ID            |
| `PUT`    | `/users/customers/:id`                    | attendant   | Atualizar cliente                |
| `DELETE` | `/users/customers/:id`                    | attendant   | Remover cliente                  |
| `GET`    | `/users/customers/:id/vehicles`           | attendant   | Listar veículos do cliente       |
| `POST`   | `/users/customers/:id/vehicles/:vehicleId`| attendant   | Associar veículo ao cliente      |
| `DELETE` | `/users/customers/:id/vehicles/:vehicleId`| attendant   | Desassociar veículo do cliente   |

Filtros disponíveis em `GET /users/customers`:

| Query param | Tipo    | Descrição                     |
|-------------|---------|-------------------------------|
| `document`  | string  | Filtrar por CPF/CNPJ          |
| `name`      | string  | Filtrar por nome (ILIKE)      |
| `email`     | string  | Filtrar por e-mail (ILIKE)    |
| `status`    | string  | `active` ou `inactive`        |

### Veículos — `/users/vehicles`

| Método   | Rota              | Role mínimo | Descrição            |
|----------|-------------------|-------------|----------------------|
| `POST`   | `/users/vehicles` | mechanic    | Criar veículo        |
| `GET`    | `/users/vehicles` | mechanic    | Listar veículos      |
| `GET`    | `/users/vehicles/:id` | mechanic | Buscar veículo por ID |
| `PUT`    | `/users/vehicles/:id` | mechanic | Atualizar veículo    |
| `DELETE` | `/users/vehicles/:id` | mechanic | Remover veículo      |

Filtros disponíveis em `GET /users/vehicles`:

| Query param | Tipo   | Descrição                  |
|-------------|--------|----------------------------|
| `plate`     | string | Filtrar por placa (ILIKE)  |
| `name`      | string | Filtrar por nome (ILIKE)   |
| `brand`     | string | Filtrar por marca (ILIKE)  |
| `color`     | string | Filtrar por cor (ILIKE)    |
| `year`      | int    | Filtrar por ano exato      |

Placas aceitas: padrão antigo (`ABC1234`) e Mercosul (`ABC1D23`).

### Empresas — `/users/companies`

| Método   | Rota                  | Role mínimo    | Descrição           |
|----------|-----------------------|----------------|---------------------|
| `POST`   | `/users/companies`    | administrator  | Criar empresa       |
| `GET`    | `/users/companies/:id`| administrator  | Buscar empresa por ID |
| `PUT`    | `/users/companies/:id`| administrator  | Atualizar empresa   |
| `DELETE` | `/users/companies/:id`| administrator  | Remover empresa     |

### Rota interna — `/internal`

Não requer autenticação. Deve ser protegida por NetworkPolicy no Kubernetes (acesso restrito ao CIDR da subnet do Lambda).

| Método | Rota                          | Descrição                                         |
|--------|-------------------------------|---------------------------------------------------|
| `GET`  | `/internal/users/by-document` | Lookup por CPF/CNPJ — usado pelo Lambda de autenticação |

#### Exemplo de uso pelo Lambda

```
GET /internal/users/by-document?document=12345678909
```

Resposta `200 OK`:
```json
{
  "person": { "id": 1, "name": "...", "email": "...", "is_active": true, ... },
  "user":   { "id": 1, "password": "<bcrypt hash>", "role": "customer", "person_id": 1 }
}
```

Fluxo de autenticação:
1. `404` → responder `401 Unauthorized` (usuário não encontrado)
2. `200` e `person.is_active == false` → responder `401 Unauthorized` (usuário inativo)
3. `200` → comparar `user.password` com a senha enviada via `bcrypt.CompareHashAndPassword`
   - match → gerar JWT
   - sem match → `401 Unauthorized`

## Testes

```bash
go test ./...
```

## Roles

| Role            | Permissões                                                   |
|-----------------|--------------------------------------------------------------|
| `administrator` | Acesso total a todos os recursos                             |
| `attendant`     | Usuários, clientes e associação de veículos a clientes       |
| `mechanic`      | Veículos                                                     |
| `customer`      | Sem acesso às rotas administrativas                          |
