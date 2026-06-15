package collector

import (
	"runtime"

	"github.com/Tridents10/Vader-Pro-VDR/vmagent/internal/models"
	"github.com/shirou/gopsutil/v3/cpu"
)

// CollectCPU gathers CPU model, cores, vendor, and architecture.
func CollectCPU() (models.CPU, error) {
	out := models.CPU{
		Architecture: runtime.GOARCH,
	}

	infos, err := cpu.Info()
	if err != nil {
		return out, err
	}

	if len(infos) > 0 {
		out.Model = infos[0].ModelName
		out.Vendor = infos[0].VendorID
	}

	// Physical cores
	if physical, perr := cpu.Counts(false); perr == nil {
		out.PhysicalCores = physical
	}

	// Logical cores
	if logical, lerr := cpu.Counts(true); lerr == nil {
		out.LogicalCores = logical
	}

	return out, nil
}
