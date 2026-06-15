package collector

import (
	"github.com/Tridents10/Vader-Pro-VDR/vmagent/internal/models"
	"github.com/shirou/gopsutil/v3/disk"
)

// CollectDisks gathers usage for all physical partitions.
func CollectDisks() ([]models.Disk, error) {
	var disks []models.Disk

	// physical-only partitions (set to true to include all)
	partitions, err := disk.Partitions(false)
	if err != nil {
		return disks, err
	}

	for _, p := range partitions {
		usage, uerr := disk.Usage(p.Mountpoint)
		if uerr != nil {
			// skip mount points we cannot read (e.g. CD-ROM)
			continue
		}

		disks = append(disks, models.Disk{
			MountPoint:  p.Mountpoint,
			Filesystem:  p.Fstype,
			TotalBytes:  usage.Total,
			UsedBytes:   usage.Used,
			FreeBytes:   usage.Free,
			UsedPercent: usage.UsedPercent,
		})
	}

	return disks, nil
}
