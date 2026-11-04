# Projeto PpDM - Agregador de Ofertas e Jogos Grátis

Projeto prático desenvolvido para a disciplina de **Programação para Dispositivos Móveis (PpDM)**. Trata-se de uma aplicação voltada para monitoramento e agregação de promoções e jogos gratuitos (Steam e Epic Games), com sistema de autenticação e painel administrativo.

---

## 🌐 Acesso Online (Deploy no Render)

A aplicação está hospedada e disponível nos seguintes links:

- **Frontend (Web App):** [https://project-ppdm.onrender.com](https://project-ppdm.onrender.com)
- **Backend (API REST):** [https://project-ppdm-backend.onrender.com](https://project-ppdm-backend.onrender.com)

> ⚠️ **Nota:** Como os serviços utilizam a camada gratuita do Render, a primeira requisição pode levar cerca de 30 a 50 segundos para responder devido ao *spin-up* (*cold start*) da instância inativa.

---

## 🛠️ Tecnologias Utilizadas

### Backend
- **Linguagem:** Go (Golang)
- **Roteamento & Servidor:** Gorilla Mux
- **Banco de Dados:** MongoDB (com Driver Oficial Go)
- **Autenticação:** JWT (JSON Web Tokens) e hash de senha com bcrypt
- **Integrações:** Consumo de APIs públicas de ofertas e jogos grátis (Steam e Epic Games)

### Frontend
- **Framework:** Vue.js 3 (Composition API / `<script setup>`)
- **Linguagem:** TypeScript
- **Estilização:** Tailwind CSS v4
- **Roteamento:** Vue Router
- **Build Tool:** Vite

---

## 🚀 Funcionalidades

- **Vitrine de Ofertas e Gratuitos:**
  - Consulta de promoções atuais e jogos 100% gratuitos por plataforma (Epic Games e Steam).
  - Filtros por preço, desconto, título e plataforma.
  - Interface responsiva com foco na usabilidade móvel.
- **Autenticação de Usuários:**
  - Cadastro e login seguro com validação e token JWT.
  - Perfil do usuário logado (`/api/me`).
- **Painel Administrativo:**
  - Gestão de usuários cadastrados (listagem, exclusão e alteração de privilégios).
  - Estatísticas de uso do sistema.

---

## 📂 Estrutura do Projeto

```text
ppdm-project/
├── backend/
│   ├── db/            # Conexão MongoDB e repositórios
│   ├── handlers/      # Controladores (Auth, Deals, Admin)
│   ├── middlewares/   # Autenticação JWT, Admin e CORS
│   ├── models/        # Modelos de dados e schemas
│   ├── routes/        # Definição e inicialização de rotas
│   ├── services/      # Integração com APIs externas de deals
│   └── main.go        # Ponto de entrada da API
└── frontend/
    ├── src/
    │   ├── components/ # Componentes visuais e skeletons
    │   ├── views/      # Telas (Home/Deals, Auth, Dashboard)
    │   ├── router.ts   # Configuração de rotas da SPA
    │   └── api.ts      # Cliente HTTP para o backend
    └── package.json
```

---

## ⚙️ Como Executar

### 1. Backend

1. Certifique-se de ter o **Go** (1.20+) e o **MongoDB** instalados/configurados.
2. Acesse a pasta `backend`:
   ```bash
   cd backend
   ```
3. Crie ou ajuste o arquivo `.env`:
   ```env
   PORT=8080
   MONGO_URI=mongodb://localhost:27017
   JWT_SECRET=sua_chave_secreta_aqui
   ```
4. Baixe as dependências e inicie o servidor:
   ```bash
   go mod tidy
   go run main.go
   ```
   *O backend estará rodando em `http://localhost:8080`.*

---

### 2. Frontend

1. Certifique-se de ter o **Node.js** instalado.
2. Acesse a pasta `frontend`:
   ```bash
   cd frontend
   ```
3. Instale as dependências:
   ```bash
   npm install
   ```
4. Inicie o servidor de desenvolvimento:
   ```bash
   npm run dev
   ```
   *Acesse a aplicação no navegador pelo endereço indicado pelo Vite (normalmente `http://localhost:5173`).*

