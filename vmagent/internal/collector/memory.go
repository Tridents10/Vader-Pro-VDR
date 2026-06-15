package collector

import (
	"github.com/Tridents10/Vader-Pro-VDR/vmagent/internal/models"
	"github.com/shirou/gopsutil/v3/mem"
)

// CollectMemory gathers RAM totals and utilization.
func CollectMemory() (models.Memory, error) {
	out := models.Memory{}

	vm, err := mem.VirtualMemory()
	if err != nil {
		return out, err
	}

	out.TotalBytes = vm.Total
	out.AvailableBytes = vm.Available
	out.UsedBytes = vm.Used
	out.UsedPercent = vm.UsedPercent

	return out, nil
}
