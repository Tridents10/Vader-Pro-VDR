//go:build windows

package collector

import (
	"strings"

	"github.com/Tridents10/Vader-Pro-VDR/vmagent/internal/models"

	"golang.org/x/sys/windows/registry"
)

// uninstallPath represents a registry hive + subkey to scan.
type uninstallPath struct {
	hive registry.Key
	path string
}

var windowsUninstallPaths = []uninstallPath{
	{registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`},
	{registry.LOCAL_MACHINE, `SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall`},
	{registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`},
}

func collectWindowsSoftware() ([]models.Application, error) {
	seen := make(map[string]bool)
	var apps []models.Application

	for _, up := range windowsUninstallPaths {
		entries := scanUninstallKey(up.hive, up.path)
		for _, app := range entries {
			key := app.Name + "|" + app.Version
			if app.Name == "" || seen[key] {
				continue
			}
			seen[key] = true
			apps = append(apps, app)
		}
	}

	return apps, nil
}

func scanUninstallKey(hive registry.Key, path string) []models.Application {
	var apps []models.Application

	root, err := registry.OpenKey(hive, path, registry.READ)
	if err != nil {
		return apps
	}
	defer root.Close()

	subkeys, err := root.ReadSubKeyNames(-1)
	if err != nil {
		return apps
	}

	for _, sk := range subkeys {
		k, err := registry.OpenKey(hive, path+`\`+sk, registry.QUERY_VALUE)
		if err != nil {
			continue
		}

		name, _, _ := k.GetStringValue("DisplayName")
		if strings.TrimSpace(name) == "" {
			k.Close()
			continue
		}

		version, _, _ := k.GetStringValue("DisplayVersion")
		publisher, _, _ := k.GetStringValue("Publisher")
		installDate, _, _ := k.GetStringValue("InstallDate")

		// Skip system components / updates with no version where appropriate
		systemComponent, _, _ := k.GetIntegerValue("SystemComponent")
		if systemComponent == 1 {
			k.Close()
			continue
		}

		apps = append(apps, models.Application{
			Name:        strings.TrimSpace(name),
			Version:     strings.TrimSpace(version),
			Publisher:   strings.TrimSpace(publisher),
			InstallDate: formatInstallDate(installDate),
		})

		k.Close()
	}

	return apps
}

// formatInstallDate converts YYYYMMDD to YYYY-MM-DD when possible.
func formatInstallDate(raw string) string {
	raw = strings.TrimSpace(raw)
	if len(raw) == 8 {
		return raw[0:4] + "-" + raw[4:6] + "-" + raw[6:8]
	}
	return raw
}
