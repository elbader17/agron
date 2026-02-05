<div align="center">
  <img src="assets/agron.png" alt="Agron Logo" width="300"/>
</div>

# agron

Agron is a lightweight Go library for AES-GCM authenticated encryption with pluggable key loaders. It provides secure encryption/decryption with associated data (AAD) support and convenient key management through environment variables or files.

## Features

- **AES-256-GCM encryption** with authenticated encryption
- **Context-aware encryption** using Associated Data (AAD)
- **Pluggable key loaders** for flexible key management
- **Zero external dependencies** - uses only Go standard library
- **Tamper detection** - ciphertext integrity verification

## Installation

```bash
go get github.com/yourusername/agron
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
