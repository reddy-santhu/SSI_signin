# SSI Sign-In System

A **Self-Sovereign Identity (SSI)** authentication system that enables passwordless login using verifiable credentials stored in mobile wallets like Lissi and Sovio.

## 🎯 Project Flow

### 1. Credential Issuance Flow
```
User → Scan Connection QR → Accept Connection → Receive Credential Offer → Accept Credential → Credential Stored in Wallet
```

### 2. Authentication Flow
```
User → Visit Login Page → Scan Login QR Code → Wallet Requests Proof → User Shares Credentials → Backend Verifies → User Logged In
```

### Complete Flow Diagram
```
┌─────────────┐
│   User       │
│  (Mobile)    │
└──────┬───────┘
       │
       │ 1. Scan Connection QR
       ▼
┌─────────────┐         ┌──────────────┐         ┌─────────────┐
│   Frontend  │────────▶│   Backend    │────────▶│   Issuer    │
│  (React)    │         │   (Go API)   │         │   Agent     │
└─────────────┘         └──────────────┘         └─────────────┘
                               │                         │
                               │                         │
                               ▼                         ▼
                        ┌──────────────┐         ┌─────────────┐
                        │  PostgreSQL  │         │ Indy Ledger  │
                        │   Database   │         │ (VON Network)│
                        └──────────────┘         └─────────────┘
                               │
                               │ 2. Credential Issued
                               ▼
                        ┌─────────────┐
                        │ User Wallet │
                        │  (Lissi/    │
                        │   Sovio)    │
                        └──────┬──────┘
                               │
                               │ 3. Login Request
                               ▼
┌─────────────┐         ┌──────────────┐         ┌─────────────┐
│   Frontend  │────────▶│   Backend    │────────▶│  Verifier  │
│  (React)    │         │   (Go API)   │         │   Agent    │
└─────────────┘         └──────────────┘         └─────────────┘
                               │                         │
                               │                         │
                               ▼                         ▼
                        ┌──────────────┐         ┌─────────────┐
                        │  PostgreSQL  │         │ Indy Ledger │
                        │   Database   │         │ (Verify)    │
                        └──────────────┘         └─────────────┘
                               │
                               │ 4. Proof Verified
                               ▼
                        ┌─────────────┐
                        │   Session   │
                        │   Created   │
                        └─────────────┘
```

## 🏗️ Architecture

### System Components

```
┌─────────────────────────────────────────────────────────────┐
│                      Frontend Layer                          │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  React Application (Port 3000)                        │  │
│  │  - Login Page with QR Code                            │  │
│  │  - Dashboard                                          │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                            │
                            │ HTTP/REST
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                      Backend Layer                            │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  Go API Server (Port 8080)                            │  │
│  │  - Authentication Handlers                            │  │
│  │  - Credential Management                              │  │
│  │  - Session Management                                 │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
         │                    │                    │
         │                    │                    │
         ▼                    ▼                    ▼
┌──────────────┐    ┌──────────────┐    ┌──────────────┐
│  PostgreSQL  │    │   Issuer     │    │  Verifier    │
│  Database    │    │   Agent      │    │   Agent      │
│  (Port 5432) │    │  (Port 8001) │    │  (Port 8003) │
└──────────────┘    └──────────────┘    └──────────────┘
                            │                    │
                            │                    │
                            └──────────┬─────────┘
                                       │
                                       ▼
                            ┌──────────────┐
                            │ Indy Ledger  │
                            │ (VON Network)│
                            │ (Port 9000)  │
                            └──────────────┘
```

### Component Responsibilities

**Frontend (React)**
- Displays login QR codes
- Handles user interface interactions
- Polls for login status
- Shows dashboard after authentication

**Backend (Go)**
- Manages authentication flow
- Makes HTTP REST API calls to ACA-Py agents
- Generates proof requests via Verifier Agent API
- Receives proof verification results
- Manages user sessions
- Stores user data
- No Aries Go SDK - uses standard HTTP client

**Issuer Agent (ACA-Py - Python)**
- Issues verifiable credentials
- Manages connections with wallets
- Creates schemas and credential definitions
- Interacts with Indy ledger
- Exposes REST API for backend to call

**Verifier Agent (ACA-Py - Python)**
- Creates proof requests
- Verifies credential proofs
- Manages authentication sessions
- Validates credential attributes
- Exposes REST API for backend to call

**PostgreSQL Database**
- Stores user information
- Manages session tokens
- Tracks proof request states

**Indy Ledger (VON Network)**
- Stores DIDs (Decentralized Identifiers)
- Stores schemas
- Stores credential definitions
- Provides decentralized trust

## 🛠️ Tech Stack

- **Backend**: Go 1.21+ (Echo framework)
  - Makes HTTP REST API calls to ACA-Py agents
  - No Aries Go SDK used - pure HTTP client implementation
- **Frontend**: React (JavaScript)
- **Agents**: Hyperledger Aries Cloud Agent Python (ACA-Py) v0.8.1
  - Python-based agents running in Docker containers
  - Backend communicates via HTTP REST API
- **Ledger**: Hyperledger Indy (BCovrin Test Network / VON Network)
- **Database**: PostgreSQL 15
- **Containerization**: Docker & Docker Compose
- **Mobile Wallets**: Lissi Wallet, Sovio Wallet

## 📋 Requirements

### System Requirements

- **Docker**: Version 20.10+
- **Docker Compose**: Version 2.0+
- **Operating System**: Linux, macOS, or Windows with WSL2

### Development Requirements (Optional)

- **Go**: Version 1.21 or higher
- **Node.js**: Version 16 or higher
- **Python**: Version 3.9 or higher (for agent scripts)

### Network Requirements

- **Ports Required**:
  - `3000`: Frontend
  - `5432`: PostgreSQL
  - `8000`: Issuer Agent Admin
  - `8001`: Issuer Agent Inbound
  - `8002`: Verifier Agent Admin
  - `8003`: Verifier Agent Inbound
  - `8080`: Backend API
  - `9000`: VON Network (if running locally)

- **Network Access**:
  - Mobile devices must be able to reach the server IP on ports 8001 and 8003
  - Server must have access to BCovrin Test Network (or local VON Network)

### External Dependencies

- **BCovrin Test Network**: For Indy ledger access (or local VON Network)
- **Mobile Wallet**: Lissi Wallet or Sovio Wallet installed on user's device

## 📁 Project Structure

```
SSI_signin/
├── backend/              # Go backend API
│   ├── handlers/        # HTTP request handlers
│   ├── services/        # Business logic layer
│   ├── models/          # Data models
│   ├── repositories/    # Database access layer
│   ├── config/          # Configuration
│   ├── middleware/      # HTTP middleware
│   └── main.go         # Application entry point
├── frontend/            # React frontend
│   ├── src/
│   │   ├── pages/      # Page components
│   │   ├── services/   # API client services
│   │   └── App.js      # Main app component
│   └── public/         # Static assets
├── issuer-agent/        # Issuer agent scripts
├── docker-compose.yml   # Service orchestration
└── README.md           # This file
```

## 🔄 Data Flow

### Credential Issuance
1. Backend creates connection invitation via Issuer Agent
2. User scans QR code with mobile wallet
3. Connection established between wallet and Issuer Agent
4. Backend sends credential offer through Issuer Agent
5. User accepts credential in wallet
6. Credential stored in user's wallet

### Authentication
1. User visits login page
2. Backend creates proof request via Verifier Agent
3. Backend generates QR code with proof request
4. User scans QR code with mobile wallet
5. Wallet presents proof of credentials
6. Verifier Agent verifies proof against ledger
7. Backend creates session and returns token
8. User redirected to dashboard

## 🔐 Security Features

- **Zero-Knowledge Proofs**: Users share only requested attributes
- **Decentralized Trust**: Credentials verified against public ledger
- **No Password Storage**: Authentication based on cryptographic proofs
- **Session Management**: Secure token-based sessions
- **DID-based Identity**: Decentralized identifiers for all parties

---

**Note**: This project uses the BCovrin Test Network for development. For production use, deploy your own Indy network or use a production-grade network.
