package agron

import (
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
)

func TestEnvHexLoader_Load(t *testing.T) {
	testKey := make([]byte, 32)
	for i := range testKey {
		testKey[i] = byte(i)
	}
	testKeyHex := hex.EncodeToString(testKey)

	envVarName := "TEST_AGRON_KEY"
	os.Setenv(envVarName, testKeyHex)
	defer os.Unsetenv(envVarName)

	loader := &EnvHexLoader{VarName: envVarName}
	loadedKey, err := loader.Load()
	if err != nil {
		t.Fatalf("EnvHexLoader.Load failed: %v", err)
	}

	if len(loadedKey) != 32 {
		t.Errorf("Expected 32 bytes, got %d", len(loadedKey))
	}

	if string(loadedKey) != string(testKey) {
		t.Error("Loaded key doesn't match original")
	}
}

func TestEnvHexLoader_Load_Empty(t *testing.T) {
	envVarName := "TEST_AGRON_KEY_EMPTY"
	os.Unsetenv(envVarName)

	loader := &EnvHexLoader{VarName: envVarName}
	_, err := loader.Load()
	if err == nil {
		t.Error("EnvHexLoader.Load should fail with empty env var")
	}
}

func TestEnvHexLoader_Load_InvalidHex(t *testing.T) {
	envVarName := "TEST_AGRON_KEY_INVALID"
	os.Setenv(envVarName, "not-valid-hex!!!")
	defer os.Unsetenv(envVarName)

	loader := &EnvHexLoader{VarName: envVarName}
	_, err := loader.Load()
	if err == nil {
		t.Error("EnvHexLoader.Load should fail with invalid hex")
	}
}

func TestFileHexLoader_Load(t *testing.T) {
	testKey := make([]byte, 32)
	for i := range testKey {
		testKey[i] = byte(i)
	}
	testKeyHex := hex.EncodeToString(testKey)

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test-key.txt")

	err := os.WriteFile(tmpFile, []byte(testKeyHex), 0600)
	if err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	loader := &FileHexLoader{Path: tmpFile}
	loadedKey, err := loader.Load()
	if err != nil {
		t.Fatalf("FileHexLoader.Load failed: %v", err)
	}

	if len(loadedKey) != 32 {
		t.Errorf("Expected 32 bytes, got %d", len(loadedKey))
	}

	if string(loadedKey) != string(testKey) {
		t.Error("Loaded key doesn't match original")
	}
}

func TestFileHexLoader_Load_MissingFile(t *testing.T) {
	loader := &FileHexLoader{Path: "/nonexistent/path/to/key.txt"}
	_, err := loader.Load()
	if err == nil {
		t.Error("FileHexLoader.Load should fail with missing file")
	}
}

func TestFileHexLoader_Load_InvalidHex(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "invalid-key.txt")

	err := os.WriteFile(tmpFile, []byte("not-valid-hex!!!"), 0600)
	if err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	loader := &FileHexLoader{Path: tmpFile}
	_, err = loader.Load()
	if err == nil {
		t.Error("FileHexLoader.Load should fail with invalid hex")
	}
}
