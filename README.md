# tech-challenge-users

Microserviço de gestão de identidade e perfil extraído do sistema legado `tech-challenge-s1`. Expõe todos os recursos como sub-rotas de `/users`, roteados via API Gateway v2 com uma única rota `ANY /users/{proxy+}`.

## Tecnologias

- **Go 1.23**
- **Gin** — HTTP router
- **GORM + pgx** — ORM com PostgreSQL
- **bcrypt** (custo 10) — hash de senhas
- **Docker Compose** — ambiente local
- **k6** — testes de stress

## Arquitetura

Hexagonal (ports & adapters):

```
cmd/api/                          → entrypoint
internal/
  domain/                         → entidades e value objects
  application/
    ports/                        → interfaces de repositório
    services/                     → lógica de negócio
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
test/
  e2e/                            → testes de integração end-to-end
  stress/                         → testes de carga com k6
```

**Fluxo de dados:**

```
HTTP Request → Handler → Service → Repository → PostgreSQL
                ↓            ↓
           (validação)  (regra de negócio)
```

## Entidades de domínio

| Entidade | Campos principais |
|----------|------------------|
| `Person` | `id`, `name`, `email`, `document` (CPF/CNPJ), `contact`, `is_active`, `address` |
| `User` | `id`, `role`, `password` (bcrypt), `person_id` |
| `Employee` | `id`, `position`, `person_id` |
| `Customer` | `id`, `type` (pf/pj), `person_id` |
| `Vehicle` | `id`, `plate`, `name`, `year`, `brand`, `color` |
| `Company` | `id`, `name`, `email`, `document` (CNPJ), `contact`, `address` |
| `CustomerVehicle` | `id`, `customer_id`, `vehicle_id` |

> `User`, `Employee` e `Customer` são vinculados a uma `Person` pelo `person_id`. A criação de um `User` com role diferente de `customer` cria automaticamente um `Employee` associado.

## Rodando localmente

**Pré-requisitos:** Docker e Docker Compose instalados.

```bash
docker compose up --build -d
```

A API sobe em `http://localhost:8081`.

**Rodando sem Docker** (requer PostgreSQL local):

```bash
export DB_HOST=localhost DB_PORT=5432 DB_USER=postgres DB_PASSWORD=secret \
       DB_NAME=usersdb DB_SSLMODE=disable DB_TIMEZONE=America/Sao_Paulo \
       HTTP_PORT=8081 ADMIN_EMAIL=admin@garage.com \
       ADMIN_PASSWORD=Admin@123 ADMIN_DOCUMENT=00000000000 ENV=development

go run ./cmd/api/main.go
```

**Resetar o banco de dados:**

```bash
docker compose down -v && docker compose up --build -d
```

## Variáveis de ambiente

| Variável | Descrição | Padrão |
|----------|-----------|--------|
| `DB_HOST` | Host do PostgreSQL | `localhost` |
| `DB_PORT` | Porta do PostgreSQL | `5432` |
| `DB_USER` | Usuário do banco | — |
| `DB_PASSWORD` | Senha do banco | — |
| `DB_NAME` | Nome do banco | — |
| `DB_SSLMODE` | SSL mode (`disable` para local) | `require` |
| `DB_TIMEZONE` | Timezone do banco | `America/Sao_Paulo` |
| `HTTP_PORT` | Porta HTTP da aplicação | `8081` |
| `HTTP_ALLOWED_ORIGINS` | CORS origins permitidas | `*` |
| `ENV` | Ambiente (`development`, `production`) | `development` |
| `ADMIN_EMAIL` | E-mail do admin inicial | — |
| `ADMIN_PASSWORD` | Senha do admin inicial | — |
| `ADMIN_DOCUMENT` | CPF do admin inicial (11 dígitos) | — |

## Autenticação

As rotas sob `/users` exigem os headers injetados pelo API Gateway após validação do JWT:

| Header | Obrigatório | Descrição |
|--------|-------------|-----------|
| `X-User-Id` | Sim | ID do usuário autenticado |
| `X-User-Email` | Não | E-mail do usuário autenticado |
| `X-User-Role` | Sim | Role: `administrator`, `attendant`, `mechanic`, `customer` |

> A rota `/internal/users/by-document` não requer esses headers — deve ser protegida por NetworkPolicy no Kubernetes (acesso restrito ao CIDR da subnet do Lambda).

## Roles e permissões

| Role | Permissões |
|------|-----------|
| `administrator` | Acesso total a todos os recursos |
| `attendant` | Usuários, clientes, funcionários e associação de veículos |
| `mechanic` | Veículos |
| `customer` | Sem acesso às rotas administrativas |

### Matriz de acesso

| Endpoint | administrator | attendant | mechanic |
|----------|:---:|:---:|:---:|
| `GET /health` | ✓ | ✓ | ✓ |
| `POST /users` | ✓ | ✓ | — |
| `GET /users` | ✓ | ✓ | — |
| `GET /users/:id` | ✓ | ✓ | — |
| `PUT /users/:id` | ✓ | ✓ | — |
| `DELETE /users/:id` | ✓ | ✓ | — |
| `POST /users/customers` | ✓ | ✓ | — |
| `GET /users/customers` | ✓ | ✓ | — |
| `GET /users/customers/:id` | ✓ | ✓ | — |
| `PUT /users/customers/:id` | ✓ | ✓ | — |
| `DELETE /users/customers/:id` | ✓ | ✓ | — |
| `GET /users/customers/:id/vehicles` | ✓ | ✓ | — |
| `POST /users/customers/:id/vehicles/:vehicleId` | ✓ | ✓ | — |
| `DELETE /users/customers/:id/vehicles/:vehicleId` | ✓ | ✓ | — |
| `POST /users/vehicles` | ✓ | ✓ | ✓ |
| `GET /users/vehicles` | ✓ | ✓ | ✓ |
| `GET /users/vehicles/:id` | ✓ | ✓ | ✓ |
| `PUT /users/vehicles/:id` | ✓ | ✓ | ✓ |
| `DELETE /users/vehicles/:id` | ✓ | ✓ | ✓ |
| `POST /users/companies` | ✓ | — | — |
| `GET /users/companies/:id` | ✓ | — | — |
| `PUT /users/companies/:id` | ✓ | — | — |
| `DELETE /users/companies/:id` | ✓ | — | — |
| `GET /users/employees/:id` | ✓ | ✓ | — |

## Rotas

### Health check

```
GET /health
```

Resposta `200 OK`:
```json
{ "status": "ok" }
```

---

### Usuários — `/users`

#### Criar usuário
```
POST /users
```

Body:
```json
{
  "name": "João Silva",
  "email": "joao.silva@garage.com",
  "password": "Senha@123",
  "role": "attendant",
  "document": "52998224725",
  "contact": "11999990000",
  "position": "Atendente",
  "address": {
    "street": "Rua das Flores",
    "number": "100",
    "neighborhood": "Centro",
    "city": "São Paulo",
    "country": "Brasil",
    "zip_code": "01310100"
  }
}
```

> Roles válidos: `administrator`, `attendant`, `mechanic`, `customer`.
> Ao criar com role diferente de `customer`, um `Employee` é criado automaticamente com o `position` informado.

Resposta `201 Created`:
```json
{
  "id": 2,
  "role": "attendant",
  "person": {
    "id": 2,
    "name": "João Silva",
    "email": "joao.silva@garage.com",
    "contact": "11999990000",
    "document": "52998224725",
    "is_active": true,
    "address": { "street": "Rua das Flores", "number": "100", ... }
  },
  "employee": {
    "id": 2,
    "position": "Atendente"
  }
}
```

#### Listar usuários
```
GET /users?name=João&email=joao
```

#### Buscar usuário por ID
```
GET /users/:id
```

#### Atualizar usuário
```
PUT /users/:id
```

Body (todos os campos opcionais):
```json
{
  "name": "João Santos",
  "email": "joao.santos@garage.com",
  "contact": "11988880000",
  "position": "Supervisor",
  "address": { "city": "Campinas" }
}
```

#### Remover usuário
```
DELETE /users/:id
```

Resposta `204 No Content`.

---

### Funcionários — `/users/employees`

#### Buscar funcionário por ID
```
GET /users/employees/:id
```

Resposta `200 OK`:
```json
{
  "id": 2,
  "position": "Atendente",
  "person_id": 2
}
```

> Use o `employee.id` retornado na criação do usuário (`POST /users`) para buscar o funcionário diretamente.

---

### Clientes — `/users/customers`

#### Criar cliente
```
POST /users/customers
```

Body:
```json
{
  "name": "Maria Souza",
  "email": "maria@email.com",
  "password": "Maria@123",
  "type": "pf",
  "document": "11144477735",
  "contact": "11977770000"
}
```

> `type`: `pf` (pessoa física — CPF) ou `pj` (pessoa jurídica — CNPJ).

#### Listar clientes
```
GET /users/customers?document=111&name=Maria&email=maria&status=active
```

| Query param | Valores |
|-------------|---------|
| `document` | CPF ou CNPJ (parcial) |
| `name` | Nome (parcial, case-insensitive) |
| `email` | E-mail (parcial, case-insensitive) |
| `status` | `active` ou `inactive` |

#### Buscar cliente por ID
```
GET /users/customers/:id
```

#### Atualizar cliente
```
PUT /users/customers/:id
```

#### Remover cliente
```
DELETE /users/customers/:id
```

Resposta `204 No Content`.

#### Listar veículos do cliente
```
GET /users/customers/:id/vehicles
```

#### Associar veículo ao cliente
```
POST /users/customers/:id/vehicles/:vehicleId
```

#### Desassociar veículo do cliente
```
DELETE /users/customers/:id/vehicles/:vehicleId
```

---

### Veículos — `/users/vehicles`

#### Criar veículo
```
POST /users/vehicles
```

Body:
```json
{
  "plate": "ABC1234",
  "name": "Corolla XEI",
  "year": 2021,
  "brand": "Toyota",
  "color": "Branco"
}
```

> Placas aceitas: padrão antigo (`ABC1234`) e Mercosul (`ABC1D23`).

#### Listar veículos
```
GET /users/vehicles?plate=ABC&brand=Toyota&color=Branco&year=2021
```

#### Buscar veículo por ID
```
GET /users/vehicles/:id
```

#### Atualizar veículo
```
PUT /users/vehicles/:id
```

Body (todos os campos opcionais):
```json
{
  "color": "Preto",
  "name": "Corolla XEI Flex"
}
```

#### Remover veículo
```
DELETE /users/vehicles/:id
```

Resposta `204 No Content`.

---

### Empresas — `/users/companies`

#### Criar empresa
```
POST /users/companies
```

Body:
```json
{
  "name": "Garage LTDA",
  "email": "contato@garage.com",
  "document": "11222333000181",
  "contact": "1133330000",
  "address": {
    "street": "Av. Paulista",
    "number": "1000",
    "city": "São Paulo",
    "country": "Brasil",
    "zip_code": "01310100"
  }
}
```

#### Buscar empresa por ID
```
GET /users/companies/:id
```

#### Atualizar empresa
```
PUT /users/companies/:id
```

#### Remover empresa
```
DELETE /users/companies/:id
```

Resposta `204 No Content`.

---

### Rota interna — `/internal`

Não requer autenticação. Deve ser protegida por NetworkPolicy no Kubernetes.

#### Buscar usuário por documento
```
GET /internal/users/by-document?document=52998224725
```

Resposta `200 OK`:
```json
{
  "person": {
    "id": 1,
    "name": "João Silva",
    "email": "joao@garage.com",
    "is_active": true
  },
  "user": {
    "id": 1,
    "password": "$2a$10$...",
    "role": "attendant",
    "person_id": 1
  }
}
```

**Fluxo de autenticação do Lambda:**
1. `404` → retornar `401 Unauthorized` (usuário não encontrado)
2. `200` com `person.is_active == false` → retornar `401 Unauthorized`
3. `200` → comparar `user.password` com bcrypt → gerar JWT ou retornar `401`

## Códigos de resposta

| Código | Situação |
|--------|----------|
| `200 OK` | Sucesso em leituras/atualizações |
| `201 Created` | Recurso criado com sucesso |
| `204 No Content` | Recurso removido com sucesso |
| `400 Bad Request` | Dados inválidos no body ou parâmetros |
| `401 Unauthorized` | Header `X-User-Id` ausente |
| `403 Forbidden` | Role sem permissão para o recurso |
| `404 Not Found` | Recurso não encontrado |
| `409 Conflict` | E-mail, documento ou placa já em uso |
| `500 Internal Server Error` | Erro interno do servidor |

## Testes

### Unitários

```bash
go test ./internal/... -v
```

### End-to-end

Requer a aplicação rodando em `http://localhost:8081`:

```bash
go test ./test/e2e/... -v
```

Variáveis de ambiente opcionais:

| Variável | Padrão |
|----------|--------|
| `API_URL` | `http://localhost:8081` |
| `ADMIN_EMAIL` | `admin@garage.com` |

### Stress (k6)

Requer [k6](https://k6.io) instalado e aplicação rodando:

```bash
k6 run test/stress/stress-test.js
```

Configuração via variáveis de ambiente:

| Variável | Descrição | Padrão |
|----------|-----------|--------|
| `BASE_URL` | URL base da API | `http://localhost:8081` |
| `USER_ID` | ID do usuário para os headers | `1` |
| `USER_EMAIL` | E-mail para os headers | `admin@admin.com` |
| `USER_ROLE` | Role para os headers | `administrator` |
| `MAX_VUS` | Número máximo de usuários virtuais | `100` |
| `P95_MS` | Threshold de latência P95 (ms) | `800` |
| `ERROR_RATE` | Taxa máxima de erro permitida | `0.01` (1%) |
