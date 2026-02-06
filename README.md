<div align="center">
  <img src="assets/agron.png" alt="Agron Logo" width="180"/>
  
  <h1>Agron</h1>
  
  <p>
    <strong>Secure, lightweight, and flexible AES-GCM encryption for Go.</strong>
  </p>

  <p>
    <a href="https://goreportcard.com/report/github.com/elbader17/agron">
      <img src="https://goreportcard.com/badge/github.com/elbader17/agron" alt="Go Report Card">
    </a>
    <a href="https://pkg.go.dev/github.com/elbader17/agron">
      <img src="https://pkg.go.dev/badge/github.com/elbader17/agron.svg" alt="GoDoc">
    </a>
    <img src="https://img.shields.io/github/license/elbader17/agron" alt="License">
    <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go&logoColor=white" alt="Go Version">
  </p>
  
  <br />
</div>

**Agron** is a lightweight Go library designed for **AES-GCM authenticated encryption** with a focus on developer experience. It features pluggable key loaders, support for Associated Data (AAD), and zero external dependencies.

---

## ⚡ Features

- 🔒 **AES-256-GCM Encryption:** Industrial-grade authenticated encryption.
- 🛡️ **Context-Aware:** Bind encryption to specific contexts (AAD) for tamper resistance.
- 🔌 **Pluggable Architecture:** Flexible key management via Environment variables or Files.
- 📦 **Zero Dependencies:** Built entirely with the Go standard library.
- 🚀 **Docker Ready:** Native support for Docker Secrets.

## 📥 Installation

```bash
go get github.com/elbader17/agron
```

## Usage

### Basic Encryption/Decryption

```go
package main

import (
    "fmt"
    "log"
    "github.com/elbader17/agron"
)

func main() {
    // Create a 32-byte key (256 bits)
    key := make([]byte, 32)
    // In production, load this securely from env/file
    
    vault, err := agron.NewVault(key)
    if err != nil {
        log.Fatal(err)
    }
    
    plaintext := []byte("secret message")
    context := []byte("user-123") // Associated data
    
    // Encrypt
    ciphertext, err := vault.Encrypt(plaintext, context)
    if err != nil {
        log.Fatal(err)
    }
    
    // Decrypt (must use same context)
    decrypted, err := vault.Decrypt(ciphertext, context)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(string(decrypted)) // "secret message"
}
```

### Key Loaders

#### Environment Variable Loader

```go
loader := &agron.EnvHexLoader{
    VarName: "MASTER_KEY",
}

key, err := loader.Load()
if err != nil {
    log.Fatal(err)
}

vault, err := agron.NewVault(key)
```

#### File Loader (Docker Secrets compatible)

```go
loader := &agron.FileHexLoader{
    Path: "/run/secrets/master_key",
}

key, err := loader.Load()
if err != nil {
    log.Fatal(err)
}

vault, err := agron.NewVault(key)
```

## API Reference

### Vault

```go
// Create a new vault with a 32-byte key
vault, err := agron.NewVault(key []byte)

// Encrypt plaintext with context (AAD)
ciphertext, err := vault.Encrypt(plaintext, context []byte)

// Decrypt ciphertext (must use same context)
plaintext, err := vault.Decrypt(ciphertext, context []byte)
```

### Key Loaders

Both loaders expect hex-encoded keys (64 hex characters = 32 bytes).

```go
// EnvHexLoader loads from environment variable
type EnvHexLoader struct {
    VarName string
}

// FileHexLoader loads from file (Docker Secrets compatible)
type FileHexLoader struct {
    Path string
}
```

### Errors

```go
var ErrInvalidKeySize = errors.New("agron: key must be 32 bytes")
var ErrDecryption     = errors.New("agron: decryption failed")
```

## Testing

Run the test suite:

```bash
go test -v ./...
```

## Security Notes

- Keys must be exactly 32 bytes (256 bits)
- The `context` parameter provides additional authentication binding
- Decryption will fail if the context or ciphertext is tampered with
- Always use secure key management in production (env vars, Docker secrets, KMS)

## License

MIT
