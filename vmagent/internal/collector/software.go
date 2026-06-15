package collector

import (
	"bufio"
	"encoding/json"
	"os/exec"
	"runtime"
	"strings"

	"github.com/Tridents10/Vader-Pro-VDR/vmagent/internal/models"
)

// CollectSoftware dispatches to the correct OS implementation.
func CollectSoftware() ([]models.Application, error) {
	switch runtime.GOOS {
	case "windows":
		return collectWindowsSoftware() // implemented in software_windows.go
	case "linux":
		return collectLinuxSoftware()
	case "darwin":
		return collectMacSoftware()
	default:
		return []models.Application{}, nil
	}
}

// ---------------- Linux ----------------

func collectLinuxSoftware() ([]models.Application, error) {
	if _, err := exec.LookPath("dpkg-query"); err == nil {
		return collectDpkg()
	}
	if _, err := exec.LookPath("rpm"); err == nil {
		return collectRpm()
	}
	return []models.Application{}, nil
}

func collectDpkg() ([]models.Application, error) {
	var apps []models.Application

	// Format: Name|Version|Maintainer
	cmd := exec.Command("dpkg-query", "-W",
		"-f=${Package}|${Version}|${Maintainer}\n")
	out, err := cmd.Output()
	if err != nil {
		return apps, err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		parts := strings.SplitN(scanner.Text(), "|", 3)
		if len(parts) < 2 || parts[0] == "" {
			continue
		}
		app := models.Application{
			Name:    parts[0],
			Version: parts[1],
		}
		if len(parts) == 3 {
			app.Publisher = parts[2]
		}
		apps = append(apps, app)
	}

	return apps, nil
}

func collectRpm() ([]models.Application, error) {
	var apps []models.Application

	// Format: Name|Version|Vendor|InstallDate
	cmd := exec.Command("rpm", "-qa",
		"--queryformat", "%{NAME}|%{VERSION}-%{RELEASE}|%{VENDOR}|%{INSTALLTIME:date}\n")
	out, err := cmd.Output()
	if err != nil {
		return apps, err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		parts := strings.SplitN(scanner.Text(), "|", 4)
		if len(parts) < 2 || parts[0] == "" {
			continue
		}
		app := models.Application{
			Name:    parts[0],
			Version: parts[1],
		}
		if len(parts) >= 3 {
			app.Publisher = parts[2]
		}
		if len(parts) == 4 {
			app.InstallDate = parts[3]
		}
		apps = append(apps, app)
	}

	return apps, nil
}

// ---------------- macOS ----------------

func collectMacSoftware() ([]models.Application, error) {
	var apps []models.Application

	cmd := exec.Command("system_profiler", "SPApplicationsDataType", "-json")
	out, err := cmd.Output()
	if err != nil {
		return apps, err
	}

	var parsed struct {
		Applications []struct {
			Name         string `json:"_name"`
			Version      string `json:"version"`
			ObtainedFrom string `json:"obtained_from"`
			LastModified string `json:"lastModified"`
			Info         string `json:"info"`
		} `json:"SPApplicationsDataType"`
	}

	if jerr := json.Unmarshal(out, &parsed); jerr != nil {
		return apps, jerr
	}

	for _, a := range parsed.Applications {
		if a.Name == "" {
			continue
		}
		apps = append(apps, models.Application{
			Name:        a.Name,
			Version:     a.Version,
			Publisher:   a.ObtainedFrom,
			InstallDate: a.LastModified,
		})
	}

	return apps, nil
}
