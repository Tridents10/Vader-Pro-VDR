package storage

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/Tridents10/Vader-Pro-VDR/vmagent/internal/models"
)

// The AES-256 Encryption Key
// This scrambles the data. It MUST match the SECRET_KEY in your Python server.py.
const secretKey = "oJ4JIXqjUD3wG0A7oXk29a1IjxbP7kZI"

// WriteInventory writes the inventory object as an AES-GCM encrypted payload
// to output/inventory.enc.
func WriteInventory(outputDir string, inv *models.Inventory) error {
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return err
	}

	// 1. Convert struct to raw JSON bytes
	jsonData, err := json.Marshal(inv)
	if err != nil {
		return err
	}

	// 2. Setup AES-GCM Cryptography
	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return err
	}

	aesGcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	// 3. Generate a cryptographic nonce
	nonce := make([]byte, aesGcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	// 4. Encrypt the JSON data
	encryptedData := aesGcm.Seal(nonce, nonce, jsonData, nil)

	// 5. Save to a secure .enc file instead of raw JSON
	outPath := filepath.Join(outputDir, "inventory.enc")
	return os.WriteFile(outPath, encryptedData, 0o644)
}
