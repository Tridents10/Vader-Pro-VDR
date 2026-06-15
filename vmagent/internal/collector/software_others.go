//go:build !windows

package collector

import "github.com/Tridents10/Vader-Pro-VDR/vmagent/internal/models"

// collectWindowsSoftware is a no-op stub on non-Windows platforms.
func collectWindowsSoftware() ([]models.Application, error) {
	return []models.Application{}, nil
}
