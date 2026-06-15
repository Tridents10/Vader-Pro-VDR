package identity

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

const idFileName = "agent.id"

// GetOrCreateUUID returns a persisted UUID, generating one on first run
// and reusing it on subsequent runs.
func GetOrCreateUUID(outputDir string) (string, error) {
	idPath := filepath.Join(outputDir, idFileName)

	// Reuse existing UUID if valid
	if data, err := os.ReadFile(idPath); err == nil {
		existing := strings.TrimSpace(string(data))
		if _, perr := uuid.Parse(existing); perr == nil {
			return existing, nil
		}
	}

	// Generate and persist a new UUID
	newID := uuid.NewString()

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(idPath, []byte(newID), 0o644); err != nil {
		return "", err
	}

	return newID, nil
}
