# Zero-Knowledge Password Manager

> This is the Go implementation of [Gassword-API](https://github.com/G4MB1T24/Gassword-API).

A highly secure, zero-knowledge password manager consisting of a **Go (Gin + SQLite)** backend API and a **React JS (Vite)** frontend client. 

All encryption and decryption happen strictly client-side. The server stores only salted bcrypt hashes for session authentication and fully encrypted ciphertexts for vault entries, ensuring that your raw credentials never leave your browser.

---

## 🔒 Security Architecture

### 1. Key Derivation (Argon2id)
- **Master Password**: The user's secret master password.
- **User Salt**: Generated server-side during registration via a CSPRNG (`crypto/rand` in Go) and returned to the client.
- **Process**: The client derives a 256-bit (32-byte) **Master Key** using **Argon2id** (configured with 64MB memory, 3 iterations, 32-byte binary output) using WebAssembly-based hashing.
- **Persistence**: The key is stored in browser-scoped `sessionStorage`. It survives tab reloads but is immediately wiped when the tab is closed or when the user logs out.

### 2. Vault Encryption (AES-256-GCM)
- **Encryption**: Every credential field (email/username, password) is encrypted client-side using the native Web Crypto API with **AES-256-GCM**.
- **Initialization Vector (Nonce)**: A random 12-byte nonce is generated per-field.
- **Payload**: The 12-byte nonce is prepended directly to the ciphertext before Base64 encoding.
- **Decryption**: The client parses the nonce from the beginning of the Base64 ciphertext and decrypts the field using the active Master Key.

---

## 🚀 Getting Started

### 1. Backend Setup 

The backend is built with Go, SQLite, and Gin, featuring Swagger documentation and automated migrations.

#### Prerequisite Environment
Create a `.env` file from the sample config:
```bash
cp .env.example .env
```

Ensure your `.env` contains:
```env
PORT=8080
DB_PATH=app.DB
JWT_TOK=your-random-32-byte-hexadecimal-jwt-secret-key
GIN_MODE=debug
```

#### Run the Server
```bash
# Install Go dependencies
go mod download

# Start the server (runs on port 8080 by default)
go run main.go
```

- **Swagger API Docs**: Once running, visit [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html) to interact with the API endpoints.

---

### 2. Frontend Setup (`/frontend`)

The frontend is a React application built with Vite and vanilla CSS.

#### Install Dependencies
Navigate to the sibling `frontend` folder:
```bash
cd ../frontend
npm install
```

#### Run Dev Server
```bash
npm run dev
```
The Dev Server configures a proxy to automatically route `/auth`, `/passwordcrud`, and `/protected` endpoints directly to `http://localhost:8080`, bypassing CORS.

#### Build for Production
```bash
npm run build
```