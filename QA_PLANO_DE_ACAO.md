# PLANO DE AÇÃO - ACCESS MANAGER
## Baseado na Análise QA Completa

**Data de Criação:** 2025-11-12
**Priorização:** MoSCoW (Must, Should, Could, Won't)
**Estimativa:** Story Points (Fibonacci: 1, 2, 3, 5, 8, 13, 21)

---

## 📋 ÍNDICE

1. [Sprint 0 - Preparação e Setup](#sprint-0---preparação-e-setup)
2. [Sprint 1 - MVP Core (Autenticação e Usuários)](#sprint-1---mvp-core)
3. [Sprint 2 - Funcionalidades Principais](#sprint-2---funcionalidades-principais)
4. [Sprint 3 - Pagamentos e Prestadores](#sprint-3---pagamentos-e-prestadores)
5. [Sprint 4 - Melhorias e Segurança](#sprint-4---melhorias-e-segurança)
6. [Backlog Futuro](#backlog-futuro)
7. [Matriz de Riscos](#matriz-de-riscos)

---

## SPRINT 0 - Preparação e Setup
**Duração:** 1 semana
**Objetivo:** Preparar ambiente e documentação base

### 🔴 MUST HAVE

#### TASK-001: Completar README.md
- **Problema:** BUG-015 (README vazio)
- **Prioridade:** P2
- **Story Points:** 2
- **Responsável:** Tech Lead
- **Descrição:**
  - Adicionar descrição do projeto
  - Instruções de setup
  - Arquitetura do sistema
  - Como rodar testes
  - Como contribuir
- **Critérios de Aceite:**
  - [ ] README com mínimo 500 palavras
  - [ ] Diagramas de arquitetura
  - [ ] Instruções testadas por novo dev

#### TASK-002: Criar Migrations Iniciais
- **Problema:** BUG-014 (Migrations não implementadas)
- **Prioridade:** P0
- **Story Points:** 5
- **Responsável:** Backend Dev
- **Descrição:**
  - Instalar golang-migrate
  - Criar estrutura /migrations/mysql/
  - Migration 001: Tabela users
  - Migration 002: Tabela professionals
  - Migration 003: Tabela laboratories
  - Migration 004: Tabela providers
  - Scripts up/down
- **Critérios de Aceite:**
  - [ ] Migrations executam sem erro
  - [ ] Rollback funciona corretamente
  - [ ] Índices criados
  - [ ] Constraints (FK, UK) configurados

```sql
-- Exemplo: 001_create_users_table.up.sql
CREATE TABLE users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    type ENUM('professional', 'laboratory', 'provider') NOT NULL,
    email_verified_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_email (email),
    INDEX idx_type (type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

#### TASK-003: Configurar Middleware Básico
- **Problema:** BUG-008 (CORS ausente), BUG-006 (Error handling)
- **Prioridade:** P1
- **Story Points:** 3
- **Responsável:** Backend Dev
- **Descrição:**
  - Criar /internal/middleware/
  - CORS middleware
  - Error handler middleware
  - Request logger middleware
  - Recovery middleware (panic)
- **Critérios de Aceite:**
  - [ ] CORS configurado para frontend
  - [ ] Erros retornam JSON padronizado
  - [ ] Logs estruturados (JSON)
  - [ ] Panics recuperados

```go
// internal/middleware/cors.go
func CORS() gin.HandlerFunc {
    return cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    })
}
```

#### TASK-004: Configurar Logger Estruturado
- **Problema:** BUG-007 (Logs insuficientes)
- **Prioridade:** P1
- **Story Points:** 2
- **Responsável:** Backend Dev
- **Descrição:**
  - Instalar zap ou logrus
  - Criar /pkg/logger/
  - Logger com níveis (debug, info, warn, error)
  - Formato JSON
  - Contexto (trace_id, user_id)
- **Critérios de Aceite:**
  - [ ] Logs em JSON
  - [ ] Níveis configuráveis via env
  - [ ] TraceID em cada request

---

## SPRINT 1 - MVP Core
**Duração:** 2 semanas
**Objetivo:** Sistema de autenticação funcional

### 🔴 MUST HAVE

#### TASK-005: Implementar Módulo de Autenticação - Parte 1 (Models)
- **Problema:** BUG-002 (Ausência de autenticação)
- **Prioridade:** P0
- **Story Points:** 5
- **Responsável:** Backend Dev
- **Descrição:**
  - Criar /app/auth/
  - entity.go (User, Session, VerificationToken, PasswordResetToken)
  - errors.go (ErrInvalidCredentials, ErrEmailTaken, etc)
  - dto.go (RegisterRequest, LoginRequest, LoginResponse)
- **Critérios de Aceite:**
  - [ ] Entidades com tags GORM
  - [ ] Validações com go-playground/validator
  - [ ] Testes unitários de entities

```go
// app/auth/entity.go
type User struct {
    ID              uint       `gorm:"primaryKey"`
    Name            string     `gorm:"not null"`
    Email           string     `gorm:"uniqueIndex;not null"`
    PasswordHash    string     `gorm:"not null"`
    Type            UserType   `gorm:"type:enum('professional','laboratory','provider');not null"`
    EmailVerifiedAt *time.Time
    CreatedAt       time.Time
    UpdatedAt       time.Time
    DeletedAt       *time.Time `gorm:"index"`
}

type UserType string

const (
    UserTypeProfessional UserType = "professional"
    UserTypeLaboratory   UserType = "laboratory"
    UserTypeProvider     UserType = "provider"
)
```

#### TASK-006: Implementar Autenticação - Parte 2 (Services)
- **Problema:** BUG-002
- **Prioridade:** P0
- **Story Points:** 8
- **Responsável:** Backend Dev
- **Descrição:**
  - /app/auth/service/password_service.go (bcrypt)
  - /app/auth/service/jwt_service.go (JWT generation/validation)
  - /app/auth/service/token_service.go (refresh tokens)
  - /pkg/crypto/ (utils)
- **Critérios de Aceite:**
  - [ ] Senhas com bcrypt (cost 12)
  - [ ] JWT com claims customizados
  - [ ] Tokens expiram corretamente
  - [ ] Testes unitários 100%

```go
// app/auth/service/jwt_service.go
type JWTService interface {
    GenerateToken(user *auth.User) (string, error)
    ValidateToken(tokenString string) (*Claims, error)
}

type Claims struct {
    UserID uint   `json:"user_id"`
    Email  string `json:"email"`
    Type   string `json:"type"`
    jwt.StandardClaims
}
```

#### TASK-007: Implementar Autenticação - Parte 3 (Repository)
- **Problema:** BUG-002
- **Prioridade:** P0
- **Story Points:** 3
- **Responsável:** Backend Dev
- **Descrição:**
  - /app/auth/repository/sql/user_repository.go
  - Métodos: Create, FindByEmail, FindByID, Update, Delete
  - Soft delete
- **Critérios de Aceite:**
  - [ ] Queries otimizadas
  - [ ] Testes com mock DB

#### TASK-008: Implementar Autenticação - Parte 4 (Use Cases)
- **Problema:** BUG-002
- **Prioridade:** P0
- **Story Points:** 8
- **Responsável:** Backend Dev
- **Descrição:**
  - /app/auth/usecase/register.go
  - /app/auth/usecase/login.go
  - /app/auth/usecase/verify_email.go
  - /app/auth/usecase/forgot_password.go
  - /app/auth/usecase/reset_password.go
- **Regras de Negócio:**
  - Senha mínimo 8 caracteres, 1 maiúscula, 1 número, 1 símbolo
  - Email único
  - Enviar email de verificação
  - Bloqueio após 5 tentativas de login
- **Critérios de Aceite:**
  - [ ] Todas as regras implementadas
  - [ ] Testes unitários de cada use case

#### TASK-009: Implementar Autenticação - Parte 5 (Controllers)
- **Problema:** BUG-004 (API REST inexistente)
- **Prioridade:** P0
- **Story Points:** 5
- **Responsável:** Backend Dev
- **Descrição:**
  - /app/auth/entrypoint/rest/auth_controller.go
  - POST /api/v1/auth/register
  - POST /api/v1/auth/login
  - POST /api/v1/auth/logout
  - POST /api/v1/auth/refresh
  - POST /api/v1/auth/forgot-password
  - POST /api/v1/auth/reset-password
  - GET /api/v1/auth/verify-email
- **Critérios de Aceite:**
  - [ ] Todos os endpoints funcionais
  - [ ] Validação de input
  - [ ] Responses padronizados
  - [ ] Status codes corretos

```go
// app/auth/entrypoint/rest/auth_controller.go
func (c *AuthController) Register(ctx *gin.Context) {
    var req dto.RegisterRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(400, ErrorResponse{Error: "VALIDATION_ERROR", Message: err.Error()})
        return
    }

    user, err := c.registerUseCase.Execute(ctx, req)
    if err != nil {
        // handle error
        return
    }

    ctx.JSON(201, user)
}
```

#### TASK-010: Middleware de Autenticação JWT
- **Problema:** BUG-002
- **Prioridade:** P0
- **Story Points:** 3
- **Responsável:** Backend Dev
- **Descrição:**
  - /internal/middleware/auth.go
  - Validar JWT em header Authorization
  - Injetar user no context
  - Verificar se email está verificado
- **Critérios de Aceite:**
  - [ ] Bloqueia requests sem token
  - [ ] Valida token corretamente
  - [ ] User disponível no context

```go
// internal/middleware/auth.go
func AuthRequired() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := extractToken(c.GetHeader("Authorization"))
        if tokenString == "" {
            c.AbortWithStatusJSON(401, gin.H{"error": "TOKEN_REQUIRED"})
            return
        }

        claims, err := jwtService.ValidateToken(tokenString)
        if err != nil {
            c.AbortWithStatusJSON(401, gin.H{"error": "INVALID_TOKEN"})
            return
        }

        c.Set("user_id", claims.UserID)
        c.Next()
    }
}
```

#### TASK-011: Testes de Integração - Autenticação
- **Prioridade:** P1
- **Story Points:** 5
- **Responsável:** Backend Dev
- **Descrição:**
  - Criar /tests/integration/auth_test.go
  - Testar fluxo completo de registro
  - Testar fluxo completo de login
  - Testar recuperação de senha
- **Critérios de Aceite:**
  - [ ] Testes com DB real (testcontainers)
  - [ ] Cobertura > 80%
  - [ ] Executam em < 30s

---

### 🟡 SHOULD HAVE

#### TASK-012: Implementar Rate Limiting
- **Problema:** BUG-009 (Rate limiting inexistente)
- **Prioridade:** P1
- **Story Points:** 3
- **Responsável:** Backend Dev
- **Descrição:**
  - Instalar go-redis/redis
  - /internal/middleware/rate_limit.go
  - Limites por endpoint
  - Limites globais por IP
- **Limites Sugeridos:**
  - /auth/login: 5/minuto
  - /auth/register: 3/hora
  - Geral: 100/minuto
- **Critérios de Aceite:**
  - [ ] 429 Too Many Requests retornado
  - [ ] Header X-RateLimit-*

#### TASK-013: Sistema de Email (SMTP)
- **Prioridade:** P1
- **Story Points:** 3
- **Responsável:** Backend Dev
- **Descrição:**
  - /pkg/email/smtp.go
  - Templates de email (HTML)
  - Envio assíncrono via Kafka (opcional)
- **Emails:**
  - Verificação de email
  - Recuperação de senha
  - Boas-vindas
- **Critérios de Aceite:**
  - [ ] Emails enviados corretamente
  - [ ] Templates responsivos

---

## SPRINT 2 - Funcionalidades Principais
**Duração:** 2 semanas
**Objetivo:** Profissionais e Agendamentos

### 🔴 MUST HAVE

#### TASK-014: Módulo de Profissionais - Models e Repository
- **Problema:** BUG-001 (Sistema não funcional)
- **Prioridade:** P0
- **Story Points:** 5
- **Responsável:** Backend Dev
- **Descrição:**
  - /app/professional/entity.go
  - /app/professional/dto.go
  - /app/professional/repository/sql/
  - Migration para tabela professionals
- **Entidades:**
  - Professional (extends User)
  - Specialties
  - WorkSchedule
- **Critérios de Aceite:**
  - [ ] Migration testada
  - [ ] Repository com CRUD completo

```go
// app/professional/entity.go
type Professional struct {
    ID           uint   `gorm:"primaryKey"`
    UserID       uint   `gorm:"uniqueIndex;not null"`
    User         auth.User `gorm:"foreignKey:UserID"`
    Document     string `gorm:"uniqueIndex;not null"` // CPF
    Phone        string
    Specialties  []Specialty `gorm:"many2many:professional_specialties;"`
    Address      Address     `gorm:"embedded"`
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

type Specialty struct {
    ID   uint   `gorm:"primaryKey"`
    Name string `gorm:"uniqueIndex"`
}
```

#### TASK-015: Módulo de Profissionais - Use Cases e Controllers
- **Prioridade:** P0
- **Story Points:** 5
- **Responsável:** Backend Dev
- **Descrição:**
  - CRUD completo de profissionais
  - Filtros (especialidade, cidade)
  - Paginação
- **Endpoints:**
  - GET /api/v1/professionals
  - GET /api/v1/professionals/:id
  - PUT /api/v1/professionals/:id
  - DELETE /api/v1/professionals/:id
- **Critérios de Aceite:**
  - [ ] Paginação funcionando
  - [ ] Filtros testados
  - [ ] Autorização (só pode editar próprio perfil)

#### TASK-016: Módulo de Agendamentos - Models e Repository
- **Prioridade:** P0
- **Story Points:** 8
- **Responsável:** Backend Dev
- **Descrição:**
  - /app/appointment/entity.go
  - /app/appointment/repository/sql/
  - Migration
- **Entidades:**
  - Appointment
  - Patient (simplificado)
- **Status:**
  - pending, confirmed, completed, cancelled, no_show
- **Critérios de Aceite:**
  - [ ] Índices em scheduled_at
  - [ ] Constraint para evitar overlap

```go
// app/appointment/entity.go
type Appointment struct {
    ID              uint      `gorm:"primaryKey"`
    ProfessionalID  uint      `gorm:"not null;index"`
    Professional    professional.Professional `gorm:"foreignKey:ProfessionalID"`
    PatientID       uint      `gorm:"not null;index"`
    Patient         Patient   `gorm:"foreignKey:PatientID"`
    ServiceType     string    `gorm:"not null"`
    ScheduledAt     time.Time `gorm:"not null;index"`
    Duration        int       `gorm:"default:30"` // minutos
    Status          AppointmentStatus `gorm:"default:'pending'"`
    Notes           string
    CheckedInAt     *time.Time
    CompletedAt     *time.Time
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
```

#### TASK-017: Módulo de Agendamentos - Use Cases
- **Prioridade:** P0
- **Story Points:** 13
- **Responsável:** Backend Dev
- **Descrição:**
  - CreateAppointment (validar conflitos)
  - UpdateAppointment
  - CancelAppointment
  - ConfirmAppointment
  - CheckInAppointment
  - CompleteAppointment
  - ListAppointments (filtros)
- **Regras de Negócio Críticas:**
  - Não permitir overlap de horários
  - Validar horário de trabalho
  - Cancelamento com menos de 24h (aviso)
- **Critérios de Aceite:**
  - [ ] Validação de conflitos funciona
  - [ ] Status transitions corretas
  - [ ] Testes unitários completos

```go
// app/appointment/usecase/create_appointment.go
func (uc *CreateAppointmentUseCase) Execute(ctx context.Context, req dto.CreateAppointmentRequest) (*Appointment, error) {
    // 1. Validar horário de trabalho
    // 2. Checar conflitos
    conflicts, err := uc.repo.FindConflicts(req.ProfessionalID, req.ScheduledAt, req.Duration)
    if len(conflicts) > 0 {
        return nil, ErrSlotUnavailable
    }

    // 3. Criar agendamento
    // 4. Enviar notificações
    // 5. Retornar
}
```

#### TASK-018: Módulo de Agendamentos - Controllers
- **Prioridade:** P0
- **Story Points:** 5
- **Responsável:** Backend Dev
- **Endpoints:**
  - POST /api/v1/appointments
  - GET /api/v1/appointments
  - GET /api/v1/appointments/:id
  - PUT /api/v1/appointments/:id
  - DELETE /api/v1/appointments/:id
  - PATCH /api/v1/appointments/:id/confirm
  - PATCH /api/v1/appointments/:id/check-in
  - PATCH /api/v1/appointments/:id/complete
- **Critérios de Aceite:**
  - [ ] Validação de input
  - [ ] Filtros (data, status, profissional)
  - [ ] Paginação

#### TASK-019: Dashboard do Profissional
- **Prioridade:** P1
- **Story Points:** 5
- **Responsável:** Backend Dev
- **Descrição:**
  - GET /api/v1/dashboard/metrics
  - Métricas agregadas
  - Próximos agendamentos
- **Métricas:**
  - Total de agendamentos (mês)
  - Total de pacientes (mês)
  - Taxa de cancelamento
  - Receita estimada
- **Critérios de Aceite:**
  - [ ] Queries otimizadas (agregações)
  - [ ] Cache de 5 minutos

---

### 🟡 SHOULD HAVE

#### TASK-020: Sistema de Notificações
- **Prioridade:** P2
- **Story Points:** 8
- **Responsável:** Backend Dev
- **Descrição:**
  - /app/notification/
  - Email notifications
  - SMS (Twilio)
  - Push notifications (FCM)
- **Eventos:**
  - Agendamento criado
  - Agendamento confirmado
  - Lembrete (24h antes)
  - Cancelamento
- **Critérios de Aceite:**
  - [ ] Assíncrono via Kafka
  - [ ] Templates configuráveis

---

## SPRINT 3 - Pagamentos e Prestadores
**Duração:** 2 semanas
**Objetivo:** Sistema de pagamentos funcional

### 🔴 MUST HAVE

#### TASK-021: Módulo de Pagamentos - Integração com Gateway
- **Problema:** BUG-010 (Kafka não utilizado), Sistema de pagamentos ausente
- **Prioridade:** P0
- **Story Points:** 13
- **Responsável:** Backend Dev
- **Descrição:**
  - /app/payment/
  - Integração com Stripe (ou PagSeguro)
  - /app/payment/repository/http/stripe_client.go
- **Funcionalidades:**
  - Criar customer
  - Criar subscription
  - Processar webhook
- **Critérios de Aceite:**
  - [ ] Integração testada (sandbox)
  - [ ] Webhooks validados (assinatura)
  - [ ] Retry em falhas

```go
// app/payment/repository/http/stripe_client.go
type StripeClient interface {
    CreateCustomer(email string, name string) (*stripe.Customer, error)
    CreateSubscription(customerID string, priceID string) (*stripe.Subscription, error)
    CancelSubscription(subscriptionID string) error
}
```

#### TASK-022: Módulo de Pagamentos - Planos e Assinaturas
- **Prioridade:** P0
- **Story Points:** 8
- **Responsável:** Backend Dev
- **Descrição:**
  - Criar tabela de planos
  - Criar tabela de assinaturas
  - CRUD de assinaturas
- **Endpoints:**
  - GET /api/v1/billing/plans
  - POST /api/v1/billing/subscribe
  - GET /api/v1/billing/subscription
  - PUT /api/v1/billing/subscription (upgrade/downgrade)
  - DELETE /api/v1/billing/subscription (cancel)
- **Critérios de Aceite:**
  - [ ] Planos cadastrados via seed
  - [ ] Upgrade imediato
  - [ ] Downgrade no fim do período

#### TASK-023: Webhook Handler para Pagamentos
- **Prioridade:** P0
- **Story Points:** 5
- **Responsável:** Backend Dev
- **Descrição:**
  - POST /webhooks/stripe
  - Validar assinatura do webhook
  - Processar eventos
- **Eventos:**
  - payment_intent.succeeded
  - payment_intent.payment_failed
  - customer.subscription.deleted
- **Critérios de Aceite:**
  - [ ] Signature validation
  - [ ] Idempotência (processar 1x)
  - [ ] Logs detalhados

#### TASK-024: Kafka Producers para Pagamentos
- **Problema:** BUG-010
- **Prioridade:** P1
- **Story Points:** 5
- **Responsável:** Backend Dev
- **Descrição:**
  - /pkg/kafka/producer.go
  - Enviar eventos de pagamento
- **Tópicos:**
  - payment.created
  - payment.confirmed
  - payment.failed
  - subscription.activated
  - subscription.cancelled
- **Critérios de Aceite:**
  - [ ] Producer funcionando
  - [ ] Mensagens com schema (Avro ou JSON)

#### TASK-025: Módulo de Prestadores - Models e Repository
- **Prioridade:** P0
- **Story Points:** 8
- **Responsável:** Backend Dev
- **Descrição:**
  - /app/provider/entity.go
  - Provider (extends User)
  - ProviderLink (vínculo com lab)
  - BankAccount
  - Availability
  - Withdrawal (saque)
- **Critérios de Aceite:**
  - [ ] Migrations criadas
  - [ ] Relacionamentos configurados

#### TASK-026: Módulo de Prestadores - Use Cases
- **Prioridade:** P0
- **Story Points:** 8
- **Responsável:** Backend Dev
- **Funcionalidades:**
  - Solicitar vínculo
  - Aprovar/Rejeitar vínculo
  - Definir disponibilidade
  - Calcular comissão
  - Solicitar saque
- **Critérios de Aceite:**
  - [ ] Cálculo de comissão correto
  - [ ] Validações de saque

#### TASK-027: Módulo de Prestadores - Controllers
- **Prioridade:** P0
- **Story Points:** 5
- **Responsável:** Backend Dev
- **Endpoints:**
  - POST /api/v1/providers/register
  - POST /api/v1/providers/link-request
  - GET /api/v1/providers/links
  - POST /api/v1/providers/availability
  - GET /api/v1/providers/revenue
  - POST /api/v1/providers/withdrawals
- **Critérios de Aceite:**
  - [ ] Todos funcionais
  - [ ] Testes de integração

---

### 🟡 SHOULD HAVE

#### TASK-028: Histórico de Pagamentos e Notas Fiscais
- **Prioridade:** P2
- **Story Points:** 5
- **Responsável:** Backend Dev
- **Descrição:**
  - GET /api/v1/billing/history
  - GET /api/v1/billing/invoices/:id/pdf
  - Gerar PDF de nota fiscal
- **Critérios de Aceite:**
  - [ ] PDFs gerados corretamente
  - [ ] Storage em S3 ou similar

---

## SPRINT 4 - Melhorias e Segurança
**Duração:** 2 semanas
**Objetivo:** Produção-ready

### 🔴 MUST HAVE

#### TASK-029: Implementar Validações de Segurança
- **Problema:** VULN-002 (Sem validação de input)
- **Prioridade:** P0
- **Story Points:** 5
- **Responsável:** Security Dev
- **Descrição:**
  - Instalar go-playground/validator/v10
  - Validar todos os DTOs
  - Sanitizar inputs
  - Prevenir SQL Injection (prepared statements)
  - Prevenir XSS
- **Critérios de Aceite:**
  - [ ] Todos os endpoints validados
  - [ ] Testes de segurança passando

#### TASK-030: Configurar HTTPS/TLS
- **Problema:** VULN-005 (Sem HTTPS)
- **Prioridade:** P0
- **Story Points:** 3
- **Responsável:** DevOps
- **Descrição:**
  - Configurar TLS no Gin
  - Let's Encrypt para prod
  - Redirecionar HTTP → HTTPS
- **Critérios de Aceite:**
  - [ ] TLS 1.3
  - [ ] Certificado válido
  - [ ] A+ no SSL Labs

#### TASK-031: Implementar Auditoria e Logging
- **Problema:** A09 - Logging Failures
- **Prioridade:** P1
- **Story Points:** 5
- **Responsável:** Backend Dev
- **Descrição:**
  - Criar tabela audit_logs
  - Logar ações sensíveis
  - Retenção de logs (90 dias)
- **Eventos Auditados:**
  - Login/Logout
  - Alteração de senha
  - Criação/Cancelamento de agendamento
  - Transações financeiras
- **Critérios de Aceite:**
  - [ ] Logs imutáveis
  - [ ] Incluir IP, user_agent, timestamp

#### TASK-032: Health Check e Readiness Probe
- **Problema:** BUG-012 (Health check inexistente)
- **Prioridade:** P1
- **Story Points:** 2
- **Responsável:** DevOps
- **Descrição:**
  - GET /health (liveness)
  - GET /ready (readiness)
  - Checar DB, Kafka, Redis
- **Critérios de Aceite:**
  - [ ] Retorna 200 se saudável
  - [ ] 503 se unhealthy
  - [ ] JSON com detalhes

```go
// GET /health
{
    "status": "healthy",
    "timestamp": "2025-11-12T10:00:00Z",
    "services": {
        "database": "up",
        "kafka": "up",
        "redis": "up"
    }
}
```

#### TASK-033: Documentação Swagger/OpenAPI
- **Problema:** BUG-011 (Documentação API ausente)
- **Prioridade:** P1
- **Story Points:** 5
- **Responsável:** Backend Dev
- **Descrição:**
  - Instalar swaggo/swag
  - Anotar controllers
  - Gerar swagger.json
  - Servir UI em /swagger
- **Critérios de Aceite:**
  - [ ] Todos os endpoints documentados
  - [ ] Exemplos de request/response
  - [ ] Schemas de erro

```go
// @Summary Login de usuário
// @Description Autentica usuário e retorna JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body dto.LoginRequest true "Credenciais"
// @Success 200 {object} dto.LoginResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
    // ...
}
```

#### TASK-034: Testes de Segurança (OWASP)
- **Problema:** VULN-001 a VULN-005
- **Prioridade:** P0
- **Story Points:** 8
- **Responsável:** Security Dev
- **Descrição:**
  - SQL Injection tests
  - XSS tests
  - CSRF tests
  - Brute force tests
  - IDOR tests
- **Ferramentas:**
  - OWASP ZAP
  - sqlmap (manual)
  - Custom scripts
- **Critérios de Aceite:**
  - [ ] 0 vulnerabilidades críticas
  - [ ] Relatório de segurança

---

### 🟡 SHOULD HAVE

#### TASK-035: Implementar Cache com Redis
- **Problema:** IMPROVEMENT-003 (Sem cache)
- **Prioridade:** P2
- **Story Points:** 5
- **Responsável:** Backend Dev
- **Descrição:**
  - Configurar Redis
  - Cache de listagens
  - Cache de dashboard
  - TTL configurável
- **Critérios de Aceite:**
  - [ ] Hit rate > 70%
  - [ ] Invalidação correta

#### TASK-036: Métricas com Prometheus
- **Problema:** IMPROVEMENT-001
- **Prioridade:** P2
- **Story Points:** 5
- **Responsável:** DevOps
- **Descrição:**
  - Instalar prometheus/client_golang
  - Expor /metrics
  - Métricas customizadas
- **Métricas:**
  - Latência de requests
  - Taxa de erro
  - Requests por endpoint
  - DB connection pool
- **Critérios de Aceite:**
  - [ ] Prometheus scraping
  - [ ] Dashboard Grafana

---

### 🟢 COULD HAVE

#### TASK-037: Internacionalização (i18n)
- **Prioridade:** P3
- **Story Points:** 8
- **Responsável:** Backend Dev
- **Descrição:**
  - Suporte a pt-BR, en-US, es-ES
  - Mensagens de erro traduzidas
  - Emails em múltiplos idiomas

#### TASK-038: Exportação de Relatórios
- **Prioridade:** P3
- **Story Points:** 8
- **Responsável:** Backend Dev
- **Descrição:**
  - GET /api/v1/reports/appointments/pdf
  - GET /api/v1/reports/financial/excel
  - Geração assíncrona (jobs)

---

## BACKLOG FUTURO

### 🔵 Módulo de Laboratórios/Clínicas
- **Estimativa:** 2 sprints
- **Tasks:**
  - CRUD de laboratórios
  - Gestão de equipe
  - Agenda unificada
  - Relatórios

### 🔵 Sistema de Pacientes
- **Estimativa:** 1 sprint
- **Tasks:**
  - CRUD de pacientes
  - Histórico médico
  - Documentos (anexos)
  - Portal do paciente

### 🔵 Telemedicina
- **Estimativa:** 3 sprints
- **Tasks:**
  - Videochamadas (WebRTC)
  - Chat em tempo real
  - Compartilhamento de tela
  - Gravação de consultas

### 🔵 Mobile App
- **Estimativa:** 4 sprints
- **Tasks:**
  - App iOS (Swift)
  - App Android (Kotlin)
  - Push notifications
  - Sincronização offline

---

## MATRIZ DE RISCOS

| ID | Risco | Probabilidade | Impacto | Mitigação |
|----|-------|---------------|---------|-----------|
| R-001 | Atraso na integração com gateway de pagamento | Alta | Alto | Começar integração cedo, ter fallback |
| R-002 | Performance ruim com muitos agendamentos | Média | Alto | Testes de carga desde sprint 2 |
| R-003 | Vulnerabilidades de segurança | Média | Crítico | Security review a cada sprint |
| R-004 | Conflitos de horário (race condition) | Alta | Alto | Usar locks de DB, testes de concorrência |
| R-005 | Falha no envio de emails | Média | Médio | Queue com retry, logs detalhados |
| R-006 | Dívida técnica acumulada | Alta | Médio | Code review rigoroso, refactoring contínuo |
| R-007 | Kafka downtime | Baixa | Alto | Fallback para processamento síncrono |

---

## ESTIMATIVAS TOTAIS

| Sprint | Story Points | Dias Úteis | Velocidade Esperada |
|--------|--------------|------------|---------------------|
| Sprint 0 | 12 | 5 | 12 |
| Sprint 1 | 42 | 10 | 21 |
| Sprint 2 | 51 | 10 | 25 |
| Sprint 3 | 52 | 10 | 26 |
| Sprint 4 | 33 | 10 | 16 |
| **TOTAL MVP** | **190** | **45 dias** | **~95** |

### Esforço por Desenvolvedor

**Assumindo:**
- Velocidade: 20 SP/sprint por dev
- 1 Backend Dev: 42 SP/sprint → 2 devs necessários
- 1 Security Dev (part-time)
- 1 DevOps (part-time)

**Equipe Recomendada:**
- 2 Backend Developers (full-time)
- 1 QA Engineer (full-time a partir de Sprint 1)
- 1 Security Engineer (part-time)
- 1 DevOps Engineer (part-time)

**Custo Estimado (MVP):**
- 2 Backend Devs × 2 meses × R$ 15k = R$ 60k
- 1 QA × 1.5 meses × R$ 12k = R$ 18k
- 1 Security (part) × 1 mês × R$ 10k = R$ 10k
- 1 DevOps (part) × 1 mês × R$ 12k = R$ 12k
- **Total: R$ 100k**

---

## DEPENDÊNCIAS

### Serviços Externos
- [ ] Stripe/PagSeguro (conta criada)
- [ ] SMTP provider (SendGrid, Mailgun)
- [ ] SMS provider (Twilio)
- [ ] Storage (S3/DigitalOcean Spaces)

### Infraestrutura
- [ ] Servidor de produção (AWS, GCP, DigitalOcean)
- [ ] Redis (ElastiCache ou self-hosted)
- [ ] Kafka (Confluent Cloud ou self-hosted)
- [ ] CI/CD (GitHub Actions, GitLab CI)
- [ ] Monitoramento (Datadog, New Relic)

---

## CRITÉRIOS DE DONE (Definition of Done)

Para uma task ser considerada completa:

- [ ] Código escrito seguindo CODING_GUIDELINES.md
- [ ] Testes unitários com cobertura > 80%
- [ ] Testes de integração (quando aplicável)
- [ ] Code review aprovado
- [ ] Linters passando (make lint)
- [ ] Documentação atualizada
- [ ] Swagger atualizado
- [ ] Migrations testadas (up/down)
- [ ] Logs estruturados adicionados
- [ ] Deploy em ambiente de staging
- [ ] QA aprovado

---

## PRIORIZAÇÃO (MoSCoW)

### Must Have (MVP)
- Autenticação completa
- CRUD de profissionais
- CRUD de agendamentos
- Sistema de pagamentos básico
- Segurança (OWASP)

### Should Have (Pós-MVP)
- Prestadores
- Laboratórios
- Notificações
- Relatórios básicos

### Could Have (Futuro)
- Telemedicina
- Mobile app
- i18n
- Analytics avançado

### Won't Have (Fora do escopo)
- Integração com wearables
- AI/ML para diagnósticos
- Blockchain para prontuários

---

## CONCLUSÃO

Este plano de ação fornece um roadmap claro para transformar o **access-manager** de um boilerplate em um sistema funcional de gerenciamento de acesso para profissionais de saúde.

### Marcos Principais

1. **Fim do Sprint 1 (2 semanas):** Sistema de autenticação funcionando
2. **Fim do Sprint 2 (4 semanas):** Agendamentos operacionais
3. **Fim do Sprint 3 (6 semanas):** Pagamentos e prestadores
4. **Fim do Sprint 4 (8 semanas):** Produção-ready (MVP)

### Próximos Passos Imediatos

1. Montar equipe de desenvolvimento
2. Setup de infraestrutura (Redis, Kafka, staging)
3. Criar repositório de front-end
4. Começar Sprint 0

---

**Documento criado por:** QA Senior
**Data:** 2025-11-12
**Versão:** 1.0
**Próxima Revisão:** Após Sprint 1
