# CASOS DE TESTE FUNCIONAIS - ACCESS MANAGER

**Projeto:** Access Manager
**Data:** 2025-11-12
**QA:** QA Senior
**Tipo:** Testes Funcionais E2E

---

## ÍNDICE
1. [Módulo de Autenticação](#1-módulo-de-autenticação)
2. [Módulo de Profissionais](#2-módulo-de-profissionais)
3. [Módulo de Laboratórios](#3-módulo-de-laboratórios)
4. [Módulo de Prestadores](#4-módulo-de-prestadores)
5. [Módulo de Agendamentos](#5-módulo-de-agendamentos)
6. [Módulo de Pagamentos](#6-módulo-de-pagamentos)
7. [Testes de Integração](#7-testes-de-integração)
8. [Testes de Segurança](#8-testes-de-segurança)
9. [Testes de Performance](#9-testes-de-performance)

---

## 1. MÓDULO DE AUTENTICAÇÃO

### 1.1 Registro de Usuário

#### TC-AUTH-001: Registro de Profissional com Dados Válidos
**Prioridade:** P0 - Crítico
**Tipo:** Positivo
**Status:** ❌ FALHA (Endpoint não existe)

**Pré-condições:**
- Sistema está rodando
- Banco de dados está acessível

**Passos:**
1. Enviar POST para `/api/v1/auth/register`
2. Body:
```json
{
  "type": "professional",
  "name": "Dr. João Silva",
  "email": "joao.silva@email.com",
  "password": "Senha@123",
  "document": "123.456.789-10",
  "phone": "+55 11 98765-4321",
  "specialties": ["Cardiologia"],
  "address": {
    "street": "Rua das Flores",
    "number": "100",
    "city": "São Paulo",
    "state": "SP",
    "zipCode": "01234-567"
  }
}
```

**Resultado Esperado:**
- Status: 201 Created
- Response:
```json
{
  "id": 1,
  "name": "Dr. João Silva",
  "email": "joao.silva@email.com",
  "type": "professional",
  "message": "Verifique seu email para ativar a conta"
}
```

**Resultado Atual:**
- Status: 404 Not Found
- Response: Endpoint não existe

**Validações:**
- ✅ Email deve ser único
- ✅ Senha deve ter mínimo 8 caracteres
- ✅ CPF deve ser válido
- ✅ Email de verificação deve ser enviado

---

#### TC-AUTH-002: Registro com Email Duplicado
**Prioridade:** P0 - Crítico
**Tipo:** Negativo
**Status:** ❌ FALHA

**Pré-condições:**
- Email `joao.silva@email.com` já existe no banco

**Passos:**
1. Tentar registrar com email existente

**Resultado Esperado:**
- Status: 409 Conflict
- Response:
```json
{
  "error": "EMAIL_ALREADY_EXISTS",
  "message": "Este email já está em uso"
}
```

**Resultado Atual:**
- Endpoint não existe

---

#### TC-AUTH-003: Registro com Senha Fraca
**Prioridade:** P1 - Alto
**Tipo:** Negativo
**Status:** ❌ FALHA

**Passos:**
1. Registrar com senha "123456"

**Resultado Esperado:**
- Status: 400 Bad Request
- Response:
```json
{
  "error": "WEAK_PASSWORD",
  "message": "Senha deve conter mínimo 8 caracteres, maiúsculas, números e símbolos"
}
```

---

#### TC-AUTH-004: Registro com CPF Inválido
**Prioridade:** P1 - Alto
**Tipo:** Negativo
**Status:** ❌ FALHA

**Passos:**
1. Registrar com CPF "111.111.111-11"

**Resultado Esperado:**
- Status: 400 Bad Request
- Response:
```json
{
  "error": "INVALID_DOCUMENT",
  "message": "CPF inválido"
}
```

---

#### TC-AUTH-005: Registro sem Campos Obrigatórios
**Prioridade:** P0 - Crítico
**Tipo:** Negativo
**Status:** ❌ FALHA

**Passos:**
1. Enviar request sem campo "email"

**Resultado Esperado:**
- Status: 400 Bad Request
- Response:
```json
{
  "error": "VALIDATION_ERROR",
  "fields": {
    "email": "Campo obrigatório"
  }
}
```

---

### 1.2 Login

#### TC-AUTH-006: Login com Credenciais Válidas
**Prioridade:** P0 - Crítico
**Tipo:** Positivo
**Status:** ❌ FALHA

**Pré-condições:**
- Usuário cadastrado e email verificado

**Passos:**
1. POST `/api/v1/auth/login`
```json
{
  "email": "joao.silva@email.com",
  "password": "Senha@123"
}
```

**Resultado Esperado:**
- Status: 200 OK
- Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refreshToken": "refresh_token_here",
  "expiresIn": 3600,
  "user": {
    "id": 1,
    "name": "Dr. João Silva",
    "email": "joao.silva@email.com",
    "type": "professional"
  }
}
```

**Validações:**
- ✅ Token JWT válido
- ✅ Token expira em 1 hora
- ✅ RefreshToken expira em 7 dias

---

#### TC-AUTH-007: Login com Email Inválido
**Prioridade:** P1 - Alto
**Tipo:** Negativo
**Status:** ❌ FALHA

**Passos:**
1. Login com email inexistente

**Resultado Esperado:**
- Status: 401 Unauthorized
- Response:
```json
{
  "error": "INVALID_CREDENTIALS",
  "message": "Email ou senha incorretos"
}
```

**Nota de Segurança:** Não informar se o email existe ou não

---

#### TC-AUTH-008: Login com Senha Incorreta
**Prioridade:** P1 - Alto
**Tipo:** Negativo
**Status:** ❌ FALHA

**Passos:**
1. Login com senha errada

**Resultado Esperado:**
- Status: 401 Unauthorized
- Response igual a TC-AUTH-007

---

#### TC-AUTH-009: Login com Conta Não Verificada
**Prioridade:** P1 - Alto
**Tipo:** Negativo
**Status:** ❌ FALHA

**Pré-condições:**
- Usuário registrado mas não verificou email

**Resultado Esperado:**
- Status: 403 Forbidden
- Response:
```json
{
  "error": "EMAIL_NOT_VERIFIED",
  "message": "Verifique seu email antes de fazer login"
}
```

---

#### TC-AUTH-010: Bloqueio Após 5 Tentativas Falhas
**Prioridade:** P0 - Crítico (Segurança)
**Tipo:** Negativo
**Status:** ❌ FALHA

**Passos:**
1. Fazer 5 tentativas de login com senha errada
2. Tentar 6ª vez

**Resultado Esperado:**
- Status: 429 Too Many Requests
- Response:
```json
{
  "error": "ACCOUNT_LOCKED",
  "message": "Conta bloqueada por 15 minutos após múltiplas tentativas falhas",
  "retryAfter": 900
}
```

---

### 1.3 Refresh Token

#### TC-AUTH-011: Renovar Token com Refresh Token Válido
**Prioridade:** P0 - Crítico
**Tipo:** Positivo
**Status:** ❌ FALHA

**Passos:**
1. POST `/api/v1/auth/refresh`
```json
{
  "refreshToken": "valid_refresh_token"
}
```

**Resultado Esperado:**
- Status: 200 OK
- Response: Novo JWT

---

#### TC-AUTH-012: Refresh Token Expirado
**Prioridade:** P1 - Alto
**Tipo:** Negativo
**Status:** ❌ FALHA

**Resultado Esperado:**
- Status: 401 Unauthorized
- Usuário deve fazer login novamente

---

### 1.4 Logout

#### TC-AUTH-013: Logout com Token Válido
**Prioridade:** P1 - Alto
**Tipo:** Positivo
**Status:** ❌ FALHA

**Passos:**
1. POST `/api/v1/auth/logout`
2. Header: `Authorization: Bearer {token}`

**Resultado Esperado:**
- Status: 200 OK
- Token deve ser invalidado (blacklist)

---

### 1.5 Recuperação de Senha

#### TC-AUTH-014: Solicitar Recuperação de Senha
**Prioridade:** P1 - Alto
**Tipo:** Positivo
**Status:** ❌ FALHA

**Passos:**
1. POST `/api/v1/auth/forgot-password`
```json
{
  "email": "joao.silva@email.com"
}
```

**Resultado Esperado:**
- Status: 200 OK
- Email com link de reset enviado
- Link expira em 1 hora

---

#### TC-AUTH-015: Resetar Senha com Token Válido
**Prioridade:** P1 - Alto
**Tipo:** Positivo
**Status:** ❌ FALHA

**Passos:**
1. POST `/api/v1/auth/reset-password`
```json
{
  "token": "reset_token",
  "newPassword": "NovaSenha@123"
}
```

**Resultado Esperado:**
- Status: 200 OK
- Senha alterada
- Tokens antigos invalidados

---

### 1.6 Verificação de Email

#### TC-AUTH-016: Verificar Email com Token Válido
**Prioridade:** P0 - Crítico
**Tipo:** Positivo
**Status:** ❌ FALHA

**Passos:**
1. GET `/api/v1/auth/verify-email?token=verification_token`

**Resultado Esperado:**
- Status: 200 OK
- Conta ativada
- Usuário pode fazer login

---

---

## 2. MÓDULO DE PROFISSIONAIS

### 2.1 CRUD de Profissionais

#### TC-PROF-001: Listar Todos os Profissionais
**Prioridade:** P1 - Alto
**Tipo:** Positivo
**Status:** ❌ FALHA

**Pré-condições:**
- Autenticado como admin

**Passos:**
1. GET `/api/v1/professionals`
2. Header: `Authorization: Bearer {token}`

**Resultado Esperado:**
- Status: 200 OK
- Response:
```json
{
  "data": [
    {
      "id": 1,
      "name": "Dr. João Silva",
      "email": "joao.silva@email.com",
      "specialties": ["Cardiologia"],
      "createdAt": "2025-11-12T10:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 1,
    "totalPages": 1
  }
}
```

---

#### TC-PROF-002: Listar Profissionais com Paginação
**Prioridade:** P1 - Alto
**Tipo:** Positivo
**Status:** ❌ FALHA

**Passos:**
1. GET `/api/v1/professionals?page=2&limit=10`

**Resultado Esperado:**
- Retornar 10 registros da página 2

---

#### TC-PROF-003: Filtrar Profissionais por Especialidade
**Prioridade:** P2 - Médio
**Tipo:** Positivo
**Status:** ❌ FALHA

**Passos:**
1. GET `/api/v1/professionals?specialty=Cardiologia`

**Resultado Esperado:**
- Apenas cardiologistas

---

#### TC-PROF-004: Buscar Profissional por ID
**Prioridade:** P1 - Alto
**Tipo:** Positivo
**Status:** ❌ FALHA

**Passos:**
1. GET `/api/v1/professionals/1`

**Resultado Esperado:**
- Status: 200 OK
- Detalhes completos do profissional

---

#### TC-PROF-005: Buscar Profissional Inexistente
**Prioridade:** P1 - Alto
**Tipo:** Negativo
**Status:** ❌ FALHA

**Passos:**
1. GET `/api/v1/professionals/99999`

**Resultado Esperado:**
- Status: 404 Not Found
```json
{
  "error": "PROFESSIONAL_NOT_FOUND",
  "message": "Profissional não encontrado"
}
```

---

#### TC-PROF-006: Atualizar Dados do Profissional
**Prioridade:** P1 - Alto
**Tipo:** Positivo
**Status:** ❌ FALHA

**Pré-condições:**
- Autenticado como próprio profissional

**Passos:**
1. PUT `/api/v1/professionals/1`
```json
{
  "name": "Dr. João Silva Filho",
  "phone": "+55 11 91234-5678"
}
```

**Resultado Esperado:**
- Status: 200 OK
- Dados atualizados

---

#### TC-PROF-007: Tentar Atualizar Profissional de Outro Usuário
**Prioridade:** P0 - Crítico (Segurança)
**Tipo:** Negativo
**Status:** ❌ FALHA

**Pré-condições:**
- Autenticado como profissional ID 2
- Tentar atualizar profissional ID 1

**Resultado Esperado:**
- Status: 403 Forbidden
```json
{
  "error": "FORBIDDEN",
  "message": "Você não tem permissão para esta ação"
}
```

---

#### TC-PROF-008: Deletar Profissional
**Prioridade:** P1 - Alto
**Tipo:** Positivo
**Status:** ❌ FALHA

**Passos:**
1. DELETE `/api/v1/professionals/1`

**Resultado Esperado:**
- Status: 204 No Content
- Soft delete (manter histórico)

---

### 2.2 Dashboard do Profissional

#### TC-PROF-009: Visualizar Dashboard
**Prioridade:** P1 - Alto
**Tipo:** Positivo
**Status:** ❌ FALHA

**Passos:**
1. GET `/api/v1/professionals/dashboard`

**Resultado Esperado:**
- Status: 200 OK
```json
{
  "metrics": {
    "totalAppointmentsThisMonth": 45,
    "totalPatientsThisMonth": 32,
    "revenueThisMonth": 12500.00,
    "cancellationRate": 5.2,
    "nextAppointments": [...]
  },
  "charts": {
    "appointmentsByDay": [...],
    "revenueByMonth": [...]
  }
}
```

---

---

## 3. MÓDULO DE LABORATÓRIOS

### 3.1 CRUD de Laboratórios

#### TC-LAB-001: Criar Laboratório
**Prioridade:** P1 - Alto
**Tipo:** Positivo
**Status:** ❌ FALHA

**Passos:**
1. POST `/api/v1/laboratories`
```json
{
  "name": "Laboratório ABC",
  "cnpj": "12.345.678/0001-90",
  "email": "contato@lababc.com",
  "phone": "+55 11 3333-4444",
  "examTypes": ["Hemograma", "Raio-X"],
  "address": {...}
}
```

**Resultado Esperado:**
- Status: 201 Created
- Laboratório criado

---

#### TC-LAB-002: Criar Laboratório com CNPJ Duplicado
**Prioridade:** P1 - Alto
**Tipo:** Negativo
**Status:** ❌ FALHA

**Resultado Esperado:**
- Status: 409 Conflict
```json
{
  "error": "CNPJ_ALREADY_EXISTS",
  "message": "Este CNPJ já está cadastrado"
}
```

---

#### TC-LAB-003: Adicionar Profissional à Equipe
**Prioridade:** P1 - Alto
**Tipo:** Positivo
**Status:** ❌ FALHA

**Passos:**
1. POST `/api/v1/laboratories/1/team`
```json
{
  "professionalId": 5,
  "role": "Técnico de Raio-X",
  "permissions": ["view_patients", "create_reports"]
}
```

**Resultado Esperado:**
- Status: 201 Created
- Profissional vinculado

---

#### TC-LAB-004: Remover Profissional da Equipe
**Prioridade:** P1 - Alto
**Tipo:** Positivo
**Status:** ❌ FALHA

**Passos:**
1. DELETE `/api/v1/laboratories/1/team/5`

**Resultado Esperado:**
- Status: 204 No Content
- Vínculo removido (soft delete)

---

#### TC-LAB-005: Listar Equipe do Laboratório
**Prioridade:** P1 - Alto
**Tipo:** Positivo
**Status:** ❌ FALHA

**Passos:**
1. GET `/api/v1/laboratories/1/team`

**Resultado Esperado:**
- Status: 200 OK
- Lista de profissionais vinculados

---

---

## 4. MÓDULO DE PRESTADORES

### 4.1 Onboarding de Prestadores

#### TC-PROV-001: Registrar como Prestador
**Prioridade:** P1 - Alto
**Tipo:** Positivo
**Status:** ❌ FALHA

**Passos:**
1. POST `/api/v1/providers/register`
```json
{
  "name": "Dr. Maria Oliveira",
  "email": "maria.oliveira@email.com",
  "password": "Senha@123",
  "document": "987.654.321-00",
  "specialties": ["Radiologia"],
  "bankAccount": {
    "bank": "001",
    "agency": "1234",
    "account": "12345-6",
    "accountType": "checking"
  }
}
```

**Resultado Esperado:**
- Status: 201 Created
- Prestador cadastrado

---

#### TC-PROV-002: Solicitar Vínculo com Clínica
**Prioridade:** P1 - Alto
**Tipo:** Positivo
**Status:** ❌ FALHA

**Passos:**
1. POST `/api/v1/providers/link-request`
```json
{
  "laboratoryId": 1,
  "message": "Gostaria de trabalhar na clínica ABC"
}
```

**Resultado Esperado:**
- Status: 201 Created
- Solicitação enviada
- Notificação para administrador da clínica

---

#### TC-PROV-003: Clínica Aprovar Vínculo
**Prioridade:** P1 - Alto
**Tipo:** Positivo
**Status:** ❌ FALHA

**Pré-condições:**
- Autenticado como admin da clínica

**Passos:**
1. PUT `/api/v1/laboratories/invitations/1`
```json
{
  "status": "approved",
  "terms": {
    "commissionRate": 30,
    "paymentFrequency": "weekly"
  }
}
```

**Resultado Esperado:**
- Status: 200 OK
- Vínculo ativo
- Notificação para prestador

---

#### TC-PROV-004: Clínica Rejeitar Vínculo
**Prioridade:** P2 - Médio
**Tipo:** Positivo
**Status:** ❌ FALHA

**Passos:**
1. PUT `/api/v1/laboratories/invitations/1`
```json
{
  "status": "rejected",
  "reason": "Não temos vagas no momento"
}
```

**Resultado Esperado:**
- Status: 200 OK
- Notificação para prestador

---

### 4.2 Disponibilidade do Prestador

#### TC-PROV-005: Definir Disponibilidade
**Prioridade:** P1 - Alto
**Tipo:** Positivo
**Status:** ❌ FALHA

**Passos:**
1. POST `/api/v1/providers/availability`
```json
{
  "laboratoryId": 1,
  "weekdays": [
    {
      "day": "monday",
      "slots": [
        {"start": "08:00", "end": "12:00"},
        {"start": "14:00", "end": "18:00"}
      ]
    }
  ]
}
```

**Resultado Esperado:**
- Status: 201 Created
- Disponibilidade salva

---

#### TC-PROV-006: Bloquear Horário Específico
**Prioridade:** P2 - Médio
**Tipo:** Positivo
**Status:** ❌ FALHA

**Passos:**
1. POST `/api/v1/providers/availability/block`
```json
{
  "date": "2025-11-20",
  "start": "10:00",
  "end": "12:00",
  "reason": "Compromisso pessoal"
}
```

**Resultado Esperado:**
- Status: 201 Created
- Horário bloqueado

---

#### TC-PROV-007: Definir Período de Férias
**Prioridade:** P2 - Médio
**Tipo:** Positivo
**Status:** ❌ FALHA

**Passos:**
1. POST `/api/v1/providers/availability/vacation`
```json
{
  "startDate": "2025-12-20",
  "endDate": "2026-01-05"
}
```

**Resultado Esperado:**
- Status: 201 Created
- Período bloqueado

---

### 4.3 Financeiro do Prestador

#### TC-PROV-008: Visualizar Receita Total
**Prioridade:** P1 - Alto
**Tipo:** Positivo
**Status:** ❌ FALHA

**Passos:**
1. GET `/api/v1/providers/revenue`

**Resultado Esperado:**
- Status: 200 OK
```json
{
  "totalRevenue": 15000.00,
  "availableForWithdrawal": 12000.00,
  "pendingPayments": 3000.00,
  "byLaboratory": [
    {
      "laboratoryId": 1,
      "laboratoryName": "Lab ABC",
      "revenue": 10000.00
    }
  ]
}
```

---

#### TC-PROV-009: Solicitar Saque
**Prioridade:** P1 - Alto
**Tipo:** Positivo
**Status:** ❌ FALHA

**Passos:**
1. POST `/api/v1/providers/withdrawals`
```json
{
  "amount": 5000.00,
  "bankAccountId": 1
}
```

**Resultado Esperado:**
- Status: 201 Created
- Saque em processamento
- Evento Kafka: `withdrawal.requested`

---

#### TC-PROV-010: Solicitar Saque Maior que Disponível
**Prioridade:** P1 - Alto
**Tipo:** Negativo
**Status:** ❌ FALHA

**Passos:**
1. Tentar sacar R$ 20.000 com apenas R$ 12.000 disponível

**Resultado Esperado:**
- Status: 400 Bad Request
```json
{
  "error": "INSUFFICIENT_FUNDS",
  "message": "Saldo insuficiente para saque"
}
```

---

#### TC-PROV-011: Visualizar Histórico de Saques
**Prioridade:** P2 - Médio
**Tipo:** Positivo
**Status:** ❌ FALHA

**Passos:**
1. GET `/api/v1/providers/withdrawals`

**Resultado Esperado:**
- Status: 200 OK
- Lista de saques (pending, completed, failed)

---

---

## 5. MÓDULO DE AGENDAMENTOS

### 5.1 CRUD de Agendamentos

#### TC-APPT-001: Criar Agendamento
**Prioridade:** P0 - Crítico
**Tipo:** Positivo
**Status:** ❌ FALHA

**Passos:**
1. POST `/api/v1/appointments`
```json
{
  "professionalId": 1,
  "patientId": 10,
  "serviceType": "Consulta de Cardiologia",
  "scheduledAt": "2025-11-20T14:00:00Z",
  "duration": 30,
  "notes": "Paciente com histórico de hipertensão"
}
```

**Resultado Esperado:**
- Status: 201 Created
- Agendamento criado
- Notificações enviadas (email/SMS)

---

#### TC-APPT-002: Criar Agendamento em Horário Ocupado
**Prioridade:** P0 - Crítico
**Tipo:** Negativo
**Status:** ❌ FALHA

**Pré-condições:**
- Profissional já tem agendamento às 14:00

**Resultado Esperado:**
- Status: 409 Conflict
```json
{
  "error": "SLOT_UNAVAILABLE",
  "message": "Horário não disponível",
  "availableSlots": ["14:30", "15:00", "15:30"]
}
```

---

#### TC-APPT-003: Criar Agendamento Fora do Horário de Trabalho
**Prioridade:** P1 - Alto
**Tipo:** Negativo
**Status:** ❌ FALHA

**Passos:**
1. Tentar agendar às 22:00 (profissional trabalha 8h-18h)

**Resultado Esperado:**
- Status: 400 Bad Request
```json
{
  "error": "OUTSIDE_WORKING_HOURS",
  "message": "Profissional não trabalha neste horário"
}
```

---

#### TC-APPT-004: Listar Agendamentos do Profissional
**Prioridade:** P1 - Alto
**Tipo:** Positivo
**Status:** ❌ FALHA

**Passos:**
1. GET `/api/v1/appointments?professionalId=1&date=2025-11-20`

**Resultado Esperado:**
- Status: 200 OK
- Lista de agendamentos do dia

---

#### TC-APPT-005: Filtrar Agendamentos por Status
**Prioridade:** P2 - Médio
**Tipo:** Positivo
**Status:** ❌ FALHA

**Passos:**
1. GET `/api/v1/appointments?status=pending`

**Resultado Esperado:**
- Apenas agendamentos pendentes

---

#### TC-APPT-006: Confirmar Agendamento
**Prioridade:** P1 - Alto
**Tipo:** Positivo
**Status:** ❌ FALHA

**Passos:**
1. PATCH `/api/v1/appointments/1/confirm`

**Resultado Esperado:**
- Status: 200 OK
- Status alterado para "confirmed"
- Notificação enviada ao paciente

---

#### TC-APPT-007: Cancelar Agendamento
**Prioridade:** P1 - Alto
**Tipo:** Positivo
**Status:** ❌ FALHA

**Passos:**
1. DELETE `/api/v1/appointments/1`
```json
{
  "reason": "Paciente não pode comparecer"
}
```

**Resultado Esperado:**
- Status: 200 OK
- Status alterado para "cancelled"
- Horário liberado
- Notificação enviada

---

#### TC-APPT-008: Cancelar Agendamento com Menos de 24h
**Prioridade:** P1 - Alto
**Tipo:** Positivo
**Status:** ❌ FALHA

**Pré-condições:**
- Agendamento para daqui 2 horas

**Resultado Esperado:**
- Status: 200 OK
- Taxa de cancelamento aplicada (se configurado)
- Warning retornado

---

#### TC-APPT-009: Reagendar Consulta
**Prioridade:** P1 - Alto
**Tipo:** Positivo
**Status:** ❌ FALHA

**Passos:**
1. PUT `/api/v1/appointments/1`
```json
{
  "scheduledAt": "2025-11-21T10:00:00Z"
}
```

**Resultado Esperado:**
- Status: 200 OK
- Horário antigo liberado
- Novo horário reservado

---

#### TC-APPT-010: Marcar Presença
**Prioridade:** P2 - Médio
**Tipo:** Positivo
**Status:** ❌ FALHA

**Passos:**
1. PATCH `/api/v1/appointments/1/check-in`

**Resultado Esperado:**
- Status: 200 OK
- CheckedInAt preenchido

---

#### TC-APPT-011: Marcar Conclusão
**Prioridade:** P1 - Alto
**Tipo:** Positivo
**Status:** ❌ FALHA

**Passos:**
1. PATCH `/api/v1/appointments/1/complete`
```json
{
  "notes": "Paciente apresentou melhora",
  "prescriptions": ["..."]
}
```

**Resultado Esperado:**
- Status: 200 OK
- Status alterado para "completed"
- Pagamento processado

---

---

## 6. MÓDULO DE PAGAMENTOS

### 6.1 Planos e Assinaturas

#### TC-PAY-001: Listar Planos Disponíveis
**Prioridade:** P1 - Alto
**Tipo:** Positivo
**Status:** ❌ FALHA

**Passos:**
1. GET `/api/v1/billing/plans`

**Resultado Esperado:**
- Status: 200 OK
```json
{
  "plans": [
    {
      "id": "basic",
      "name": "Básico",
      "price": 49.90,
      "interval": "monthly",
      "features": [
        "Até 50 agendamentos/mês",
        "1 usuário",
        "Suporte por email"
      ]
    },
    {
      "id": "professional",
      "name": "Profissional",
      "price": 99.90,
      "interval": "monthly",
      "features": [
        "Agendamentos ilimitados",
        "Até 5 usuários",
        "Suporte prioritário"
      ]
    }
  ]
}
```

---

#### TC-PAY-002: Criar Assinatura
**Prioridade:** P0 - Crítico
**Tipo:** Positivo
**Status:** ❌ FALHA

**Passos:**
1. POST `/api/v1/billing/subscribe`
```json
{
  "planId": "professional",
  "paymentMethod": "credit_card",
  "cardToken": "tok_visa_4242"
}
```

**Resultado Esperado:**
- Status: 201 Created
- Assinatura ativa
- Evento Kafka: `subscription.activated`

---

#### TC-PAY-003: Criar Assinatura com Cartão Inválido
**Prioridade:** P1 - Alto
**Tipo:** Negativo
**Status:** ❌ FALHA

**Resultado Esperado:**
- Status: 400 Bad Request
```json
{
  "error": "PAYMENT_FAILED",
  "message": "Cartão recusado"
}
```

---

#### TC-PAY-004: Visualizar Assinatura Atual
**Prioridade:** P1 - Alto
**Tipo:** Positivo
**Status:** ❌ FALHA

**Passos:**
1. GET `/api/v1/billing/subscription`

**Resultado Esperado:**
- Status: 200 OK
```json
{
  "subscription": {
    "id": "sub_123",
    "planId": "professional",
    "status": "active",
    "currentPeriodStart": "2025-11-01",
    "currentPeriodEnd": "2025-12-01",
    "cancelAtPeriodEnd": false
  }
}
```

---

#### TC-PAY-005: Atualizar Plano (Upgrade)
**Prioridade:** P1 - Alto
**Tipo:** Positivo
**Status:** ❌ FALHA

**Passos:**
1. PUT `/api/v1/billing/subscription`
```json
{
  "newPlanId": "enterprise"
}
```

**Resultado Esperado:**
- Status: 200 OK
- Plano atualizado imediatamente
- Cobrança proporcional

---

#### TC-PAY-006: Atualizar Plano (Downgrade)
**Prioridade:** P1 - Alto
**Tipo:** Positivo
**Status:** ❌ FALHA

**Passos:**
1. PUT `/api/v1/billing/subscription`
```json
{
  "newPlanId": "basic"
}
```

**Resultado Esperado:**
- Status: 200 OK
- Downgrade aplicado no fim do período atual

---

#### TC-PAY-007: Cancelar Assinatura
**Prioridade:** P1 - Alto
**Tipo:** Positivo
**Status:** ❌ FALHA

**Passos:**
1. DELETE `/api/v1/billing/subscription`
```json
{
  "reason": "Não estou mais usando",
  "cancelImmediately": false
}
```

**Resultado Esperado:**
- Status: 200 OK
- Assinatura cancelada no fim do período
- Evento Kafka: `subscription.cancelled`

---

#### TC-PAY-008: Cancelar Assinatura Imediatamente
**Prioridade:** P2 - Médio
**Tipo:** Positivo
**Status:** ❌ FALHA

**Passos:**
1. DELETE com `cancelImmediately: true`

**Resultado Esperado:**
- Status: 200 OK
- Acesso bloqueado imediatamente
- Possível reembolso proporcional

---

### 6.2 Histórico de Pagamentos

#### TC-PAY-009: Visualizar Histórico
**Prioridade:** P2 - Médio
**Tipo:** Positivo
**Status:** ❌ FALHA

**Passos:**
1. GET `/api/v1/billing/history`

**Resultado Esperado:**
- Status: 200 OK
```json
{
  "payments": [
    {
      "id": "pay_123",
      "amount": 99.90,
      "status": "succeeded",
      "date": "2025-11-01",
      "invoiceUrl": "https://..."
    }
  ]
}
```

---

#### TC-PAY-010: Download de Nota Fiscal
**Prioridade:** P2 - Médio
**Tipo:** Positivo
**Status:** ❌ FALHA

**Passos:**
1. GET `/api/v1/billing/invoices/inv_123/pdf`

**Resultado Esperado:**
- Status: 200 OK
- Content-Type: application/pdf
- Download de PDF

---

### 6.3 Webhooks de Pagamento

#### TC-PAY-011: Webhook de Pagamento Confirmado
**Prioridade:** P0 - Crítico
**Tipo:** Integração
**Status:** ❌ FALHA

**Passos:**
1. Gateway envia POST `/webhooks/payment`
```json
{
  "event": "payment.succeeded",
  "paymentId": "pay_123",
  "amount": 99.90
}
```

**Resultado Esperado:**
- Status: 200 OK
- Assinatura ativada
- Evento Kafka: `payment.confirmed`

---

#### TC-PAY-012: Webhook de Pagamento Falho
**Prioridade:** P0 - Crítico
**Tipo:** Integração
**Status:** ❌ FALHA

**Passos:**
1. Gateway envia evento `payment.failed`

**Resultado Esperado:**
- Status: 200 OK
- Assinatura suspensa
- Email de notificação enviado
- Evento Kafka: `payment.failed`

---

---

## 7. TESTES DE INTEGRAÇÃO

### 7.1 Fluxo E2E Completo: Profissional

#### TC-INT-001: Fluxo Completo de Onboarding e Primeiro Agendamento
**Prioridade:** P0 - Crítico
**Tipo:** E2E
**Status:** ❌ FALHA

**Cenário:**
Um novo profissional se cadastra e cria seu primeiro agendamento

**Passos:**
1. Registrar profissional
2. Verificar email
3. Fazer login
4. Assinar plano "Professional"
5. Configurar perfil
6. Cadastrar paciente
7. Criar agendamento
8. Confirmar agendamento
9. Marcar presença
10. Concluir atendimento

**Resultado Esperado:**
- Todos os passos executados com sucesso
- Dados consistentes em todas as etapas
- Notificações enviadas corretamente

---

### 7.2 Fluxo E2E: Prestador

#### TC-INT-002: Fluxo Completo de Vínculo com Clínica
**Prioridade:** P1 - Alto
**Tipo:** E2E
**Status:** ❌ FALHA

**Passos:**
1. Prestador se registra
2. Solicita vínculo com clínica
3. Admin da clínica aprova
4. Prestador define disponibilidade
5. Clínica cria agendamento
6. Prestador visualiza agenda
7. Prestador completa atendimento
8. Sistema calcula comissão
9. Prestador solicita saque

**Resultado Esperado:**
- Fluxo completo sem erros
- Cálculos financeiros corretos

---

---

## 8. TESTES DE SEGURANÇA

### 8.1 Autenticação e Autorização

#### TC-SEC-001: Acessar Endpoint Protegido Sem Token
**Prioridade:** P0 - Crítico
**Tipo:** Segurança
**Status:** ❌ FALHA

**Passos:**
1. GET `/api/v1/professionals` sem header Authorization

**Resultado Esperado:**
- Status: 401 Unauthorized

---

#### TC-SEC-002: Acessar com Token Expirado
**Prioridade:** P0 - Crítico
**Tipo:** Segurança
**Status:** ❌ FALHA

**Resultado Esperado:**
- Status: 401 Unauthorized
```json
{
  "error": "TOKEN_EXPIRED",
  "message": "Token expirado. Use refresh token."
}
```

---

#### TC-SEC-003: Acessar com Token Inválido
**Prioridade:** P0 - Crítico
**Tipo:** Segurança
**Status:** ❌ FALHA

**Passos:**
1. Enviar token malformado

**Resultado Esperado:**
- Status: 401 Unauthorized

---

#### TC-SEC-004: SQL Injection em Filtros
**Prioridade:** P0 - Crítico
**Tipo:** Segurança
**Status:** ❌ FALHA

**Passos:**
1. GET `/api/v1/professionals?name='; DROP TABLE users; --`

**Resultado Esperado:**
- Query parametrizada (sem execução de SQL malicioso)
- Retorno normal ou erro de validação

---

#### TC-SEC-005: XSS em Campos de Texto
**Prioridade:** P0 - Crítico
**Tipo:** Segurança
**Status:** ❌ FALHA

**Passos:**
1. Criar profissional com name: `<script>alert('XSS')</script>`

**Resultado Esperado:**
- Texto sanitizado
- Sem execução de script

---

#### TC-SEC-006: CSRF em Endpoints Críticos
**Prioridade:** P1 - Alto
**Tipo:** Segurança
**Status:** ❌ FALHA

**Passos:**
1. Tentar fazer DELETE de outro domínio sem token CSRF

**Resultado Esperado:**
- Status: 403 Forbidden (se CSRF token obrigatório)

---

#### TC-SEC-007: Rate Limiting em Login
**Prioridade:** P0 - Crítico
**Tipo:** Segurança
**Status:** ❌ FALHA

**Passos:**
1. Fazer 100 requisições de login em 1 minuto

**Resultado Esperado:**
- Após limite (ex: 5/minuto):
  - Status: 429 Too Many Requests

---

#### TC-SEC-008: Brute Force de Senha
**Prioridade:** P0 - Crítico
**Tipo:** Segurança
**Status:** ❌ FALHA

**Passos:**
1. Tentar 1000 senhas diferentes

**Resultado Esperado:**
- Conta bloqueada após 5 tentativas
- CAPTCHA obrigatório

---

### 8.2 Controle de Acesso

#### TC-SEC-009: Profissional Acessar Dados de Outro
**Prioridade:** P0 - Crítico
**Tipo:** Segurança
**Status:** ❌ FALHA

**Passos:**
1. Profissional A tentar GET `/api/v1/professionals/B/patients`

**Resultado Esperado:**
- Status: 403 Forbidden

---

#### TC-SEC-010: IDOR - Modificar ID em Request
**Prioridade:** P0 - Crítico
**Tipo:** Segurança
**Status:** ❌ FALHA

**Passos:**
1. PUT `/api/v1/appointments/999` (agendamento de outro)

**Resultado Esperado:**
- Status: 403 Forbidden ou 404 Not Found

---

---

## 9. TESTES DE PERFORMANCE

### 9.1 Carga e Stress

#### TC-PERF-001: Listar 10.000 Profissionais
**Prioridade:** P2 - Médio
**Tipo:** Performance
**Status:** ❌ FALHA

**Objetivo:**
- Tempo de resposta < 500ms

**Passos:**
1. Popular banco com 10k profissionais
2. GET `/api/v1/professionals?limit=100`

**Métricas:**
- Tempo de resposta
- Uso de CPU
- Uso de memória
- Queries executadas

---

#### TC-PERF-002: 100 Requisições Concorrentes de Login
**Prioridade:** P1 - Alto
**Tipo:** Performance
**Status:** ❌ FALHA

**Objetivo:**
- Taxa de sucesso > 99%
- Tempo médio < 200ms

---

#### TC-PERF-003: Criar 1000 Agendamentos Simultâneos
**Prioridade:** P1 - Alto
**Tipo:** Stress
**Status:** ❌ FALHA

**Objetivo:**
- Sistema não deve travar
- Sem deadlocks
- Todos salvos corretamente

---

#### TC-PERF-004: Paginação com 1M de Registros
**Prioridade:** P2 - Médio
**Tipo:** Performance
**Status:** ❌ FALHA

**Objetivo:**
- Primeira página < 200ms
- Página 10.000 < 500ms

---

---

## 10. RESUMO EXECUTIVO DE TESTES

### 10.1 Estatísticas Gerais

| Categoria | Total de Casos | Passaram | Falharam | Cobertura |
|-----------|----------------|----------|----------|-----------|
| Autenticação | 16 | 0 | 16 | 0% |
| Profissionais | 9 | 0 | 9 | 0% |
| Laboratórios | 5 | 0 | 5 | 0% |
| Prestadores | 11 | 0 | 11 | 0% |
| Agendamentos | 11 | 0 | 11 | 0% |
| Pagamentos | 12 | 0 | 12 | 0% |
| Integração | 2 | 0 | 2 | 0% |
| Segurança | 10 | 0 | 10 | 0% |
| Performance | 4 | 0 | 4 | 0% |
| **TOTAL** | **80** | **0** | **80** | **0%** |

### 10.2 Casos Críticos (P0)

| ID | Descrição | Status |
|----|-----------|--------|
| TC-AUTH-001 | Registro de usuário | ❌ FALHA |
| TC-AUTH-006 | Login com credenciais | ❌ FALHA |
| TC-AUTH-010 | Bloqueio após tentativas | ❌ FALHA |
| TC-APPT-001 | Criar agendamento | ❌ FALHA |
| TC-APPT-002 | Conflito de horário | ❌ FALHA |
| TC-PAY-002 | Criar assinatura | ❌ FALHA |
| TC-SEC-001 | Proteção de endpoints | ❌ FALHA |
| TC-SEC-007 | Rate limiting | ❌ FALHA |

**Total Crítico:** 8 casos
**Taxa de Falha Crítica:** 🔴 100%

### 10.3 Recomendação

**Status do Sistema:** 🔴 NÃO PRONTO PARA PRODUÇÃO

**Motivos:**
1. Nenhuma funcionalidade implementada
2. 0% de cobertura de testes
3. Vulnerabilidades críticas de segurança
4. Impossível testar APIs inexistentes

**Próximos Passos:**
1. Implementar funcionalidades básicas
2. Executar testes assim que endpoints existirem
3. Criar suite de testes automatizados
4. CI/CD com testes obrigatórios

---

**Documento gerado por:** QA Senior
**Data:** 2025-11-12
**Versão:** 1.0
