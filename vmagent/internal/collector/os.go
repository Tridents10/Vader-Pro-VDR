package collector

import (
	"runtime"
	"strings"
	"time"

	"github.com/Tridents10/VDR/vmagent/internal/models"
	"github.com/shirou/gopsutil/v3/host"
)

// CollectOS gathers operating system details using gopsutil.
func CollectOS() (models.OS, error) {
	out := models.OS{
		Architecture: runtime.GOARCH,
	}

	info, err := host.Info()
	if err != nil {
		return out, err
	}

	out.Name = strings.Title(info.OS) // e.g. "Windows", "Linux", "Darwin"
	out.Edition = info.Platform       // e.g. "ubuntu", "Microsoft Windows 11 Pro"
	out.Version = info.PlatformVersion
	out.Kernel = info.KernelVersion
	out.BuildNumber = info.KernelVersion

	// Normalize macOS naming
	if info.OS == "darwin" {
		out.Name = "macOS"
	}

	// Boot time
	out.BootTime = time.Unix(int64(info.BootTime), 0).UTC().Format(time.RFC3339)

	return out, nil
}
