package agron

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

// KeyLoader define cómo buscar la llave maestra
type KeyLoader interface {
	Load() ([]byte, error)
}

// EnvHexLoader lee una variable de entorno que contiene un string HEX
// Ejemplo: MASTER_KEY="a1b2..."
type EnvHexLoader struct {
	VarName string
}

func (l *EnvHexLoader) Load() ([]byte, error) {
	val := os.Getenv(l.VarName)
	if val == "" {
		return nil, fmt.Errorf("env var %s is empty", l.VarName)
	}
	return hex.DecodeString(strings.TrimSpace(val))
}

// FileHexLoader lee un archivo que contiene un string HEX (Ideal Docker Secrets)
type FileHexLoader struct {
	Path string
}

func (l *FileHexLoader) Load() ([]byte, error) {
	content, err := os.ReadFile(l.Path)
	if err != nil {
		return nil, err
	}
	return hex.DecodeString(strings.TrimSpace(string(content)))
}
